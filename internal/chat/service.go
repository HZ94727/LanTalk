package chat

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"path/filepath"
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
	maxImageBytes    = 4 * 1024 * 1024
	debugEchoPeerID  = "debug-echo-bot"
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
		peers:         peersByID(state.Peers),
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
	if err := s.store.SaveProfile(s.self); err != nil {
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
	if err := s.store.SaveProfile(s.self); err != nil {
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
		Kind:       "text",
		Text:       text,
		Timestamp:  time.Now().UnixMilli(),
		Direction:  "outbound",
		Status:     "sent",
	}

	payload := directMessage{
		ID:         message.ID,
		SenderID:   self.ID,
		SenderName: self.Name,
		Kind:       message.Kind,
		Text:       text,
		MediaName:  "",
		MediaType:  "",
		Timestamp:  message.Timestamp,
	}

	if peer.Source == "debug" {
		s.appendMessage(peerID, message)
		go s.sendDebugReply(peer, message)
		return nil
	}

	if err := s.sendDirect(peer, payload); err != nil {
		message.Status = "failed"
		s.appendMessage(peerID, message)
		return err
	}

	s.appendMessage(peerID, message)
	return nil
}

func (s *Service) SendImageMessage(peerID, dataURL, fileName string) error {
	dataURL = strings.TrimSpace(dataURL)
	if dataURL == "" {
		return errors.New("image payload is empty")
	}

	mediaType, payloadSize, err := validateImageDataURL(dataURL)
	if err != nil {
		return err
	}
	if payloadSize > maxImageBytes {
		return fmt.Errorf("image too large, max %d MB", maxImageBytes/(1024*1024))
	}

	fileName = strings.TrimSpace(fileName)
	if fileName == "" {
		fileName = "image"
	}
	fileName = filepath.Base(fileName)

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
		Kind:       "image",
		Text:       dataURL,
		MediaName:  fileName,
		MediaType:  mediaType,
		Timestamp:  time.Now().UnixMilli(),
		Direction:  "outbound",
		Status:     "sent",
	}

	payload := directMessage{
		ID:         message.ID,
		SenderID:   self.ID,
		SenderName: self.Name,
		Kind:       message.Kind,
		Text:       message.Text,
		MediaName:  message.MediaName,
		MediaType:  message.MediaType,
		Timestamp:  message.Timestamp,
	}

	if peer.Source == "debug" {
		s.appendMessage(peerID, message)
		go s.sendDebugReply(peer, message)
		return nil
	}

	if err := s.sendDirect(peer, payload); err != nil {
		message.Status = "failed"
		s.appendMessage(peerID, message)
		return err
	}

	s.appendMessage(peerID, message)
	return nil
}

func (s *Service) EnsureDebugPeer() Peer {
	s.mu.Lock()
	defer s.mu.Unlock()

	peer, ok := s.peers[debugEchoPeerID]
	if ok {
		return peer
	}

	peer = Peer{
		ID:         debugEchoPeerID,
		Name:       "Echo Bot",
		Address:    "Local simulator",
		ListenPort: 0,
		LastSeen:   time.Now().UnixMilli(),
		Source:     "debug",
	}
	s.peers[peer.ID] = peer
	_ = s.store.SavePeer(peer)
	s.emitPeersLocked()
	return peer
}

func (s *Service) AddManualPeer(name, address string, port int) (Peer, error) {
	name = strings.TrimSpace(name)
	address = strings.TrimSpace(address)
	if address == "" {
		return Peer{}, errors.New("address cannot be empty")
	}
	if port <= 0 || port > 65535 {
		return Peer{}, errors.New("port must be between 1 and 65535")
	}
	if name == "" {
		name = fmt.Sprintf("Manual %s:%d", address, port)
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	for _, peer := range s.peers {
		if peer.Address == address && peer.ListenPort == port {
			if peer.Name != name {
				peer.Name = name
				peer.LastSeen = time.Now().UnixMilli()
				if peer.Source == "" {
					peer.Source = "manual"
				}
				s.peers[peer.ID] = peer
				_ = s.store.SavePeer(peer)
				s.emitPeersLocked()
			}
			return peer, nil
		}
	}

	peer := Peer{
		ID:         uuid.NewString(),
		Name:       name,
		Address:    address,
		ListenPort: port,
		LastSeen:   time.Now().UnixMilli(),
		Source:     "manual",
	}
	s.peers[peer.ID] = peer
	_ = s.store.SavePeer(peer)
	s.emitPeersLocked()
	return peer, nil
}

func (s *Service) UpdateLanguage(language string) (Settings, error) {
	s.mu.RLock()
	currentTheme := s.settings.Theme
	s.mu.RUnlock()
	settings := normalizeSettings(Settings{Language: language, Theme: currentTheme})

	s.mu.Lock()
	s.settings = settings
	if err := s.store.SaveSettings(s.settings); err != nil {
		s.mu.Unlock()
		return Settings{}, err
	}
	s.mu.Unlock()

	s.emitSettings()
	return settings, nil
}

func (s *Service) UpdateTheme(theme string) (Settings, error) {
	s.mu.RLock()
	currentLanguage := s.settings.Language
	s.mu.RUnlock()
	settings := normalizeSettings(Settings{Language: currentLanguage, Theme: theme})

	s.mu.Lock()
	s.settings = settings
	if err := s.store.SaveSettings(s.settings); err != nil {
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
	newPath, err := s.store.Move(targetDir)
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
	message, err := directPayloadToMessage(payload)
	if err != nil {
		return
	}

	peerID := payload.SenderID
	if peerID == "" {
		return
	}

	message.PeerID = peerID
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
		Source:     "lan",
	}
	_ = s.store.SavePeer(s.peers[incoming.ID])

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
		if peer.Source != "" && peer.Source != "lan" {
			continue
		}
		lastSeen := time.UnixMilli(peer.LastSeen)
		if now.Sub(lastSeen) > peerExpiry {
			delete(s.peers, id)
			_ = s.store.DeletePeer(id)
			changed = true
		}
	}

	if changed {
		s.emitPeersLocked()
	}
}

func (s *Service) sendDebugReply(peer Peer, source ChatMessage) {
	delay := 550 * time.Millisecond
	time.Sleep(delay)

	reply := ChatMessage{
		ID:         uuid.NewString(),
		PeerID:     peer.ID,
		SenderID:   peer.ID,
		SenderName: peer.Name,
		Kind:       source.Kind,
		Text:       debugReplyText(source),
		MediaName:  source.MediaName,
		MediaType:  source.MediaType,
		Timestamp:  time.Now().UnixMilli(),
		Direction:  "inbound",
		Status:     "received",
	}
	s.appendMessage(peer.ID, reply)
}

func debugReplyText(message ChatMessage) string {
	if message.Kind == "image" {
		return message.Text
	}

	trimmed := strings.TrimSpace(message.Text)
	if trimmed == "" {
		return "Echo Bot: ..."
	}
	return "Echo Bot: " + trimmed
}

func (s *Service) appendMessage(peerID string, message ChatMessage) {
	s.mu.Lock()
	s.conversations[peerID] = append(s.conversations[peerID], message)
	messages := append([]ChatMessage(nil), s.conversations[peerID]...)
	_ = s.store.SaveMessage(message)
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

	_ = conn.SetWriteDeadline(time.Now().Add(12 * time.Second))
	return json.NewEncoder(conn).Encode(payload)
}

func directPayloadToMessage(payload directMessage) (ChatMessage, error) {
	kind := payload.Kind
	if kind == "" {
		kind = "text"
	}

	message := ChatMessage{
		ID:         payload.ID,
		SenderID:   payload.SenderID,
		SenderName: payload.SenderName,
		Kind:       kind,
		Text:       payload.Text,
		MediaName:  payload.MediaName,
		MediaType:  payload.MediaType,
		Timestamp:  payload.Timestamp,
		Direction:  "inbound",
		Status:     "received",
	}

	switch kind {
	case "image":
		mediaType, payloadSize, err := validateImageDataURL(payload.Text)
		if err != nil {
			return ChatMessage{}, err
		}
		if payloadSize > maxImageBytes {
			return ChatMessage{}, fmt.Errorf("image too large")
		}
		if strings.TrimSpace(message.MediaType) == "" {
			message.MediaType = mediaType
		}
		if strings.TrimSpace(message.MediaName) == "" {
			message.MediaName = "image"
		}
	case "text":
		if strings.TrimSpace(payload.Text) == "" {
			return ChatMessage{}, errors.New("message cannot be empty")
		}
	default:
		return ChatMessage{}, fmt.Errorf("unsupported message kind: %s", kind)
	}

	return message, nil
}

func validateImageDataURL(value string) (string, int, error) {
	if !strings.HasPrefix(value, "data:image/") {
		return "", 0, errors.New("only image data URLs are supported")
	}
	parts := strings.SplitN(value, ",", 2)
	if len(parts) != 2 {
		return "", 0, errors.New("invalid image payload")
	}
	header := parts[0]
	if !strings.HasSuffix(header, ";base64") {
		return "", 0, errors.New("image payload must be base64 encoded")
	}
	mediaType := strings.TrimPrefix(header, "data:")
	mediaType = strings.TrimSuffix(mediaType, ";base64")
	if mediaType == "" {
		return "", 0, errors.New("missing image media type")
	}

	decodedLen := base64.StdEncoding.DecodedLen(len(parts[1]))
	if decodedLen <= 0 {
		return "", 0, errors.New("invalid image payload")
	}

	if _, err := base64.StdEncoding.DecodeString(parts[1]); err != nil {
		return "", 0, errors.New("invalid base64 image payload")
	}

	return mediaType, decodedLen, nil
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
	normalized := settings
	switch normalized.Language {
	case "en-US", "zh-CN":
	default:
		normalized.Language = "zh-CN"
	}

	switch normalized.Theme {
	case "midnight", "paper", "forest":
	default:
		normalized.Theme = "midnight"
	}

	return normalized
}

func peersByID(peers []Peer) map[string]Peer {
	peerMap := make(map[string]Peer, len(peers))
	for _, peer := range peers {
		peerMap[peer.ID] = peer
	}
	return peerMap
}
