package chat

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)

const (
	discoveryPort    = 48555
	broadcastEvery   = 2 * time.Second
	peerExpiry       = 8 * time.Second
	maxMessageLength = 2000
)

type Callbacks struct {
	OnProfileChanged  func(Profile)
	OnSettingsChanged func(Settings)
	OnDataPathChanged func(string)
	OnPeersChanged    func([]Peer)
	OnMessage         func(string, []ChatMessage)
}

type Service struct {
	store *Store

	mu            sync.RWMutex
	self          Profile
	settings      Settings
	peers         map[string]Peer
	conversations map[string][]ChatMessage

	callbacks Callbacks
	listener  net.Listener
	udpConn   *net.UDPConn
	cancel    context.CancelFunc
}

func NewService(appName string, callbacks Callbacks) (*Service, error) {
	store, err := NewStore(appName)
	if err != nil {
		return nil, err
	}

	state, err := store.Load()
	if err != nil {
		return nil, err
	}

	self := state.Profile
	if self.ID == "" {
		self.ID = uuid.NewString()
	}
	if strings.TrimSpace(self.Name) == "" {
		self.Name = defaultName(self.ID)
	}

	return &Service{
		store:         store,
		self:          self,
		settings:      normalizeSettings(state.Settings),
		peers:         make(map[string]Peer),
		conversations: state.Conversations,
		callbacks:     callbacks,
	}, nil
}

func (s *Service) Start(ctx context.Context) error {
	runCtx, cancel := context.WithCancel(ctx)
	s.cancel = cancel

	listener, err := net.Listen("tcp4", ":0")
	if err != nil {
		return err
	}
	s.listener = listener

	addr, ok := listener.Addr().(*net.TCPAddr)
	if !ok {
		return errors.New("failed to resolve tcp address")
	}

	s.mu.Lock()
	s.self.ListenPort = addr.Port
	if err := s.persistLocked(); err != nil {
		s.mu.Unlock()
		return err
	}
	profile := s.self
	s.mu.Unlock()

	udpConn, err := net.ListenUDP("udp4", &net.UDPAddr{Port: discoveryPort})
	if err != nil {
		return err
	}
	s.udpConn = udpConn

	go s.acceptLoop(runCtx)
	go s.discoveryListenLoop(runCtx)
	go s.discoveryBroadcastLoop(runCtx)
	go s.peerCleanupLoop(runCtx)

	s.emitProfile(profile)
	s.emitSettings()
	s.emitPeers()

	return nil
}

func (s *Service) Stop() {
	if s.cancel != nil {
		s.cancel()
	}
	if s.udpConn != nil {
		_ = s.udpConn.Close()
	}
	if s.listener != nil {
		_ = s.listener.Close()
	}
}

func (s *Service) Snapshot() Snapshot {
	s.mu.RLock()
	defer s.mu.RUnlock()

	conversations := make(map[string][]ChatMessage, len(s.conversations))
	for peerID, messages := range s.conversations {
		conversations[peerID] = append([]ChatMessage(nil), messages...)
	}

	return Snapshot{
		Self:          s.self,
		Settings:      s.settings,
		Peers:         s.sortedPeersLocked(),
		Conversations: conversations,
	}
}

func (s *Service) UpdateDisplayName(name string) (Profile, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return Profile{}, errors.New("nickname cannot be empty")
	}

	s.mu.Lock()
	s.self.Name = name
	if err := s.persistLocked(); err != nil {
		s.mu.Unlock()
		return Profile{}, err
	}
	profile := s.self
	s.mu.Unlock()

	s.emitProfile(profile)
	go s.broadcastAnnouncement()

	return profile, nil
}

func (s *Service) SendMessage(peerID, text string) error {
	text = strings.TrimSpace(text)
	if text == "" {
		return errors.New("message cannot be empty")
	}
	if len([]rune(text)) > maxMessageLength {
		return fmt.Errorf("message too long, max %d characters", maxMessageLength)
	}

	s.mu.RLock()
	peer, ok := s.peers[peerID]
	self := s.self
	s.mu.RUnlock()
	if !ok {
		return errors.New("peer is offline")
	}

	message := ChatMessage{
		ID:         uuid.NewString(),
		PeerID:     peerID,
		SenderID:   self.ID,
		SenderName: self.Name,
		Text:       text,
		Timestamp:  time.Now().UnixMilli(),
		Direction:  "outbound",
		Status:     "sent",
	}

	payload := directMessage{
		ID:         message.ID,
		SenderID:   self.ID,
		SenderName: self.Name,
		Text:       text,
		Timestamp:  message.Timestamp,
	}

	if err := s.sendDirect(peer, payload); err != nil {
		message.Status = "failed"
		s.appendMessage(peerID, message)
		return err
	}

	s.appendMessage(peerID, message)
	return nil
}

func (s *Service) UpdateLanguage(language string) (Settings, error) {
	settings := normalizeSettings(Settings{Language: language})

	s.mu.Lock()
	s.settings = settings
	if err := s.persistLocked(); err != nil {
		s.mu.Unlock()
		return Settings{}, err
	}
	s.mu.Unlock()

	s.emitSettings()
	return settings, nil
}

func (s *Service) DataPath() string {
	return s.store.Path()
}

func (s *Service) MoveDataFile(targetDir string) (string, error) {
	targetDir = strings.TrimSpace(targetDir)
	if targetDir == "" {
		return "", errors.New("target directory is empty")
	}

	s.mu.Lock()
	state := persistedState{
		Profile:       s.self,
		Settings:      s.settings,
		Conversations: s.conversations,
	}
	newPath, err := s.store.Move(state, targetDir)
	s.mu.Unlock()
	if err != nil {
		return "", err
	}

	s.emitDataPath(newPath)
	return newPath, nil
}

func (s *Service) acceptLoop(ctx context.Context) {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			if ctx.Err() != nil || errors.Is(err, net.ErrClosed) {
				return
			}
			continue
		}

		go s.handleIncoming(conn)
	}
}

func (s *Service) handleIncoming(conn net.Conn) {
	defer conn.Close()

	var payload directMessage
	if err := json.NewDecoder(conn).Decode(&payload); err != nil {
		return
	}
	if strings.TrimSpace(payload.Text) == "" {
		return
	}

	peerID := payload.SenderID
	if peerID == "" {
		return
	}

	message := ChatMessage{
		ID:         payload.ID,
		PeerID:     peerID,
		SenderID:   payload.SenderID,
		SenderName: payload.SenderName,
		Text:       payload.Text,
		Timestamp:  payload.Timestamp,
		Direction:  "inbound",
		Status:     "received",
	}

	s.appendMessage(peerID, message)
}

func (s *Service) discoveryListenLoop(ctx context.Context) {
	buffer := make([]byte, 2048)

	for {
		_ = s.udpConn.SetReadDeadline(time.Now().Add(2 * time.Second))
		n, addr, err := s.udpConn.ReadFromUDP(buffer)
		if err != nil {
			if ctx.Err() != nil {
				return
			}
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				continue
			}
			continue
		}

		var incoming announcement
		if err := json.Unmarshal(buffer[:n], &incoming); err != nil {
			continue
		}

		s.handleAnnouncement(incoming, addr)
	}
}

func (s *Service) handleAnnouncement(incoming announcement, addr *net.UDPAddr) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if incoming.ID == "" || incoming.ID == s.self.ID || incoming.ListenPort == 0 {
		return
	}

	name := strings.TrimSpace(incoming.Name)
	if name == "" {
		name = defaultName(incoming.ID)
	}

	s.peers[incoming.ID] = Peer{
		ID:         incoming.ID,
		Name:       name,
		Address:    addr.IP.String(),
		ListenPort: incoming.ListenPort,
		LastSeen:   time.Now().UnixMilli(),
	}

	s.emitPeersLocked()
}

func (s *Service) discoveryBroadcastLoop(ctx context.Context) {
	ticker := time.NewTicker(broadcastEvery)
	defer ticker.Stop()

	s.broadcastAnnouncement()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			s.broadcastAnnouncement()
		}
	}
}

func (s *Service) broadcastAnnouncement() {
	s.mu.RLock()
	payload := announcement{
		ID:         s.self.ID,
		Name:       s.self.Name,
		ListenPort: s.self.ListenPort,
	}
	s.mu.RUnlock()

	data, err := json.Marshal(payload)
	if err != nil {
		return
	}

	conn, err := net.DialUDP("udp4", nil, &net.UDPAddr{IP: net.IPv4bcast, Port: discoveryPort})
	if err != nil {
		return
	}
	defer conn.Close()

	_, _ = conn.Write(data)
}

func (s *Service) peerCleanupLoop(ctx context.Context) {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			s.cleanupPeers()
		}
	}
}

func (s *Service) cleanupPeers() {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	changed := false
	for id, peer := range s.peers {
		lastSeen := time.UnixMilli(peer.LastSeen)
		if now.Sub(lastSeen) > peerExpiry {
			delete(s.peers, id)
			changed = true
		}
	}

	if changed {
		s.emitPeersLocked()
	}
}

func (s *Service) appendMessage(peerID string, message ChatMessage) {
	s.mu.Lock()
	s.conversations[peerID] = append(s.conversations[peerID], message)
	messages := append([]ChatMessage(nil), s.conversations[peerID]...)
	_ = s.persistLocked()
	s.mu.Unlock()

	if s.callbacks.OnMessage != nil {
		s.callbacks.OnMessage(peerID, messages)
	}
}

func (s *Service) sendDirect(peer Peer, payload directMessage) error {
	conn, err := net.DialTimeout("tcp4", fmt.Sprintf("%s:%d", peer.Address, peer.ListenPort), 3*time.Second)
	if err != nil {
		return err
	}
	defer conn.Close()

	_ = conn.SetWriteDeadline(time.Now().Add(3 * time.Second))
	return json.NewEncoder(conn).Encode(payload)
}

func (s *Service) persistLocked() error {
	return s.store.Save(persistedState{
		Profile:       s.self,
		Settings:      s.settings,
		Conversations: s.conversations,
	})
}

func (s *Service) emitProfile(profile Profile) {
	if s.callbacks.OnProfileChanged != nil {
		s.callbacks.OnProfileChanged(profile)
	}
}

func (s *Service) emitSettings() {
	s.mu.RLock()
	settings := s.settings
	s.mu.RUnlock()

	if s.callbacks.OnSettingsChanged != nil {
		s.callbacks.OnSettingsChanged(settings)
	}
}

func (s *Service) emitDataPath(path string) {
	if s.callbacks.OnDataPathChanged != nil {
		s.callbacks.OnDataPathChanged(path)
	}
}

func (s *Service) emitPeers() {
	s.mu.RLock()
	peers := s.sortedPeersLocked()
	s.mu.RUnlock()

	if s.callbacks.OnPeersChanged != nil {
		s.callbacks.OnPeersChanged(peers)
	}
}

func (s *Service) emitPeersLocked() {
	if s.callbacks.OnPeersChanged != nil {
		s.callbacks.OnPeersChanged(s.sortedPeersLocked())
	}
}

func (s *Service) sortedPeersLocked() []Peer {
	peers := make([]Peer, 0, len(s.peers))
	for _, peer := range s.peers {
		peers = append(peers, peer)
	}

	sort.Slice(peers, func(i, j int) bool {
		if peers[i].Name == peers[j].Name {
			return peers[i].ID < peers[j].ID
		}
		return strings.ToLower(peers[i].Name) < strings.ToLower(peers[j].Name)
	})

	return peers
}

func defaultName(id string) string {
	shortID := id
	if len(shortID) > 8 {
		shortID = shortID[:8]
	}
	return "LanTalk-" + shortID
}

func normalizeSettings(settings Settings) Settings {
	switch settings.Language {
	case "en-US", "zh-CN":
		return settings
	default:
		return Settings{Language: "zh-CN"}
	}
}
