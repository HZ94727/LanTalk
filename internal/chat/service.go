package chat

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"LanTalk/internal/mediautil"
	"github.com/google/uuid"
)

const (
	discoveryPort    = 48555
	broadcastEvery   = 2 * time.Second
	peerExpiry       = 8 * time.Second
	maxMessageLength = 2000
	maxImageBytes    = 4 * 1024 * 1024
	maxFileBytes     = 20 * 1024 * 1024
	debugEchoPeerID  = "debug-echo-bot"
)

type Callbacks struct {
	OnProfileChanged   func(Profile)
	OnSettingsChanged  func(Settings)
	OnDataPathChanged  func(string)
	OnPeersChanged     func([]Peer)
	OnMessage          func(string, []ChatMessage)
	OnTransferProgress func(TransferProgress)
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

	if err := store.MarkPendingOutgoingFailed(); err != nil {
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

	peers := peersByID(state.Peers)
	if err := mergeConversationOnlyPeers(store, peers); err != nil {
		return nil, err
	}

	return &Service{
		store:         store,
		self:          self,
		settings:      normalizeSettings(state.Settings),
		peers:         peers,
		conversations: make(map[string][]ChatMessage),
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

	return Snapshot{
		Self:          s.self,
		Settings:      s.settings,
		Peers:         s.sortedPeersLocked(),
		Conversations: map[string][]ChatMessage{},
	}
}

func (s *Service) LoadConversation(peerID string, beforeTimestamp int64, beforeID string, limit int) (ConversationPage, error) {
	page, err := s.store.LoadConversationPage(peerID, beforeTimestamp, beforeID, limit)
	if err != nil {
		return ConversationPage{}, err
	}

	s.mu.Lock()
	existing := s.conversations[peerID]
	if beforeTimestamp <= 0 {
		s.conversations[peerID] = append([]ChatMessage(nil), page.Messages...)
	} else {
		merged := append([]ChatMessage(nil), page.Messages...)
		seen := make(map[string]struct{}, len(page.Messages)+len(existing))
		for _, message := range merged {
			seen[message.ID] = struct{}{}
		}
		for _, message := range existing {
			if _, ok := seen[message.ID]; ok {
				continue
			}
			merged = append(merged, message)
		}
		s.conversations[peerID] = merged
	}
	snapshot := append([]ChatMessage(nil), s.conversations[peerID]...)
	s.mu.Unlock()

	page.Messages = snapshot
	return page, nil
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
	if !peerCanSend(peer) {
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
		Status:     "sending",
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

	return s.queueOutgoingMessage(peer, message, payload)
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
	if !peerCanSend(peer) {
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
		MediaSize:  int64(payloadSize),
		Timestamp:  time.Now().UnixMilli(),
		Direction:  "outbound",
		Status:     "sending",
	}

	payload := directMessage{
		ID:         message.ID,
		SenderID:   self.ID,
		SenderName: self.Name,
		Kind:       message.Kind,
		Text:       message.Text,
		MediaName:  message.MediaName,
		MediaType:  message.MediaType,
		MediaSize:  message.MediaSize,
		Timestamp:  message.Timestamp,
	}

	return s.queueOutgoingMessage(peer, message, payload)
}

func (s *Service) SendFileMessage(peerID, dataURL, fileName string) error {
	dataURL = strings.TrimSpace(dataURL)
	if dataURL == "" {
		return errors.New("file payload is empty")
	}

	mediaType, payloadSize, err := validateFileDataURL(dataURL)
	if err != nil {
		return err
	}

	fileName = strings.TrimSpace(fileName)
	fileName = filepath.Base(fileName)

	kind := "file"
	if strings.HasPrefix(strings.ToLower(strings.TrimSpace(mediaType)), "image/") {
		kind = "image"
		if payloadSize > maxImageBytes {
			return fmt.Errorf("image too large, max %d MB", maxImageBytes/(1024*1024))
		}
		if fileName == "" {
			fileName = "image"
		}
	} else {
		if payloadSize > maxFileBytes {
			return fmt.Errorf("file too large, max %d MB", maxFileBytes/(1024*1024))
		}
		if fileName == "" {
			fileName = "file"
		}
	}

	s.mu.RLock()
	peer, ok := s.peers[peerID]
	self := s.self
	s.mu.RUnlock()
	if !ok {
		return errors.New("peer is offline")
	}
	if !peerCanSend(peer) {
		return errors.New("peer is offline")
	}

	message := ChatMessage{
		ID:         uuid.NewString(),
		PeerID:     peerID,
		SenderID:   self.ID,
		SenderName: self.Name,
		Kind:       kind,
		Text:       dataURL,
		MediaName:  fileName,
		MediaType:  mediaType,
		MediaSize:  int64(payloadSize),
		Timestamp:  time.Now().UnixMilli(),
		Direction:  "outbound",
		Status:     "sending",
	}

	payload := directMessage{
		ID:         message.ID,
		SenderID:   self.ID,
		SenderName: self.Name,
		Kind:       message.Kind,
		Text:       message.Text,
		MediaName:  message.MediaName,
		MediaType:  message.MediaType,
		MediaSize:  message.MediaSize,
		Timestamp:  message.Timestamp,
	}

	return s.queueOutgoingMessage(peer, message, payload)
}

func (s *Service) SendLocalFileMessage(peerID, filePath string) error {
	filePath = strings.TrimSpace(filePath)
	if filePath == "" {
		return errors.New("file path is empty")
	}

	info, err := os.Stat(filePath)
	if err != nil {
		return err
	}
	if info.IsDir() {
		return errors.New("directories are not supported")
	}

	fileName := filepath.Base(filePath)
	mediaType := mediautil.MediaTypeForFileName(fileName)
	if strings.HasPrefix(strings.ToLower(strings.TrimSpace(mediaType)), "image/") {
		if info.Size() > maxImageBytes {
			return fmt.Errorf("image too large, max %d MB", maxImageBytes/(1024*1024))
		}
	} else if info.Size() > maxFileBytes {
		return fmt.Errorf("file too large, max %d MB", maxFileBytes/(1024*1024))
	}

	raw, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	mediaType = mediautil.DetectMediaType(fileName, raw)
	return s.SendFileMessage(peerID, mediautil.EncodeDataURL(raw, mediaType), fileName)
}

func (s *Service) RetryMessage(peerID, messageID string) error {
	peerID = strings.TrimSpace(peerID)
	messageID = strings.TrimSpace(messageID)
	if peerID == "" || messageID == "" {
		return errors.New("peer id and message id are required")
	}

	s.mu.RLock()
	peer, ok := s.peers[peerID]
	self := s.self
	s.mu.RUnlock()
	if !ok {
		return errors.New("peer is offline")
	}
	if !peerCanSend(peer) {
		return errors.New("peer is offline")
	}

	message, err := s.findMessage(peerID, messageID)
	if err != nil {
		return err
	}
	if message.Direction != "outbound" {
		return errors.New("only outbound messages can be retried")
	}
	if message.Status != "failed" {
		return errors.New("message is not retryable")
	}

	payload := directMessage{
		ID:         message.ID,
		SenderID:   self.ID,
		SenderName: self.Name,
		Kind:       message.Kind,
		MediaName:  message.MediaName,
		MediaType:  message.MediaType,
		MediaSize:  message.MediaSize,
		Timestamp:  message.Timestamp,
	}

	switch message.Kind {
	case "image":
		payload.Text, err = s.store.ResolveImageDataURL(message.Text, message.MediaType)
	case "file":
		payload.Text, err = s.store.ResolveFileDataURL(message.Text, message.MediaType)
	default:
		payload.Text = message.Text
	}
	if err != nil {
		return err
	}

	message.SenderID = self.ID
	message.SenderName = self.Name
	message.Status = "sending"
	return s.queueOutgoingMessage(peer, message, payload)
}

func (s *Service) DeleteMessage(peerID, messageID string) error {
	peerID = strings.TrimSpace(peerID)
	messageID = strings.TrimSpace(messageID)
	if peerID == "" || messageID == "" {
		return errors.New("peer id and message id are required")
	}

	s.mu.Lock()
	messages, ok := s.conversations[peerID]
	if !ok {
		s.mu.Unlock()
		return errors.New("conversation not found")
	}

	nextMessages := make([]ChatMessage, 0, len(messages))
	removed := false
	var removedMessage ChatMessage
	for _, message := range messages {
		if message.ID == messageID {
			if message.Status == "sending" {
				s.mu.Unlock()
				return errors.New("message is still sending")
			}
			removed = true
			removedMessage = message
			continue
		}
		nextMessages = append(nextMessages, message)
	}
	if !removed {
		s.mu.Unlock()
		return errors.New("message not found")
	}

	s.conversations[peerID] = nextMessages
	if err := s.store.RemoveMessage(removedMessage); err != nil {
		s.mu.Unlock()
		return err
	}
	snapshot := append([]ChatMessage(nil), nextMessages...)
	s.mu.Unlock()

	if s.callbacks.OnMessage != nil {
		s.callbacks.OnMessage(peerID, snapshot)
	}
	return nil
}

func (s *Service) ClearConversation(peerID string) (int, error) {
	peerID = strings.TrimSpace(peerID)
	if peerID == "" {
		return 0, errors.New("peer id is required")
	}

	s.mu.RLock()
	messages := append([]ChatMessage(nil), s.conversations[peerID]...)
	s.mu.RUnlock()

	for _, message := range messages {
		if message.Status == "sending" {
			return 0, errors.New("conversation has sending messages")
		}
	}

	removedCount, err := s.store.ClearConversation(peerID)
	if err != nil {
		return 0, err
	}

	s.mu.Lock()
	s.conversations[peerID] = []ChatMessage{}
	s.mu.Unlock()

	if s.callbacks.OnMessage != nil {
		s.callbacks.OnMessage(peerID, []ChatMessage{})
	}
	return removedCount, nil
}

func (s *Service) findMessage(peerID, messageID string) (ChatMessage, error) {
	s.mu.RLock()
	if messages, ok := s.conversations[peerID]; ok {
		for _, message := range messages {
			if message.ID == messageID {
				s.mu.RUnlock()
				return message, nil
			}
		}
	}
	s.mu.RUnlock()

	return s.store.LoadMessage(peerID, messageID)
}

func (s *Service) ResolveImagePayload(value, mediaType string) ([]byte, string, error) {
	return s.store.ResolveImagePayload(value, mediaType)
}

func (s *Service) ResolveImageDataURL(value, mediaType string) (string, error) {
	return s.store.ResolveImageDataURL(value, mediaType)
}

func (s *Service) ResolveFilePayload(value, mediaType string) ([]byte, string, error) {
	return s.store.ResolveFilePayload(value, mediaType)
}

func (s *Service) ResolveFileDataURL(value, mediaType string) (string, error) {
	return s.store.ResolveFileDataURL(value, mediaType)
}

func (s *Service) ResolveStoredFilePath(value string) (string, error) {
	return s.store.ResolveStoredFilePath(value)
}

func (s *Service) ResolveStoredImagePath(value string) (string, error) {
	return s.store.ResolveStoredImagePath(value)
}

func (s *Service) IsStoredFileAvailable(value string) bool {
	return s.store.IsStoredFileAvailable(value)
}

func (s *Service) LoadStorageStats() (StorageStats, error) {
	return s.store.StorageStats()
}

func (s *Service) CleanupUnusedMedia() (MediaCleanupResult, error) {
	return s.store.CleanupUnusedMedia()
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
	return s.store.DataDir()
}

func (s *Service) MoveDataFile(targetDir string) (string, error) {
	targetDir = strings.TrimSpace(targetDir)
	if targetDir == "" {
		return "", errors.New("target directory is empty")
	}

	s.mu.Lock()
	_, err := s.store.Move(targetDir)
	s.mu.Unlock()
	if err != nil {
		return "", err
	}

	newPath := s.store.DataDir()
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
	_, _ = s.appendMessage(peerID, message)
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
			hasHistory, err := s.store.HasMessagesForPeer(id)
			if err != nil {
				continue
			}
			if hasHistory {
				continue
			}
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
		MediaSize:  source.MediaSize,
		Timestamp:  time.Now().UnixMilli(),
		Direction:  "inbound",
		Status:     "received",
	}
	_, _ = s.appendMessage(peer.ID, reply)
}

func debugReplyText(message ChatMessage) string {
	if message.Kind == "image" || message.Kind == "file" {
		return message.Text
	}

	trimmed := strings.TrimSpace(message.Text)
	if trimmed == "" {
		return "Echo Bot: ..."
	}
	return "Echo Bot: " + trimmed
}

func upsertConversationMessage(messages []ChatMessage, message ChatMessage) []ChatMessage {
	next := append([]ChatMessage(nil), messages...)
	for index := range next {
		if next[index].ID == message.ID {
			next[index] = message
			return next
		}
	}
	return append(next, message)
}

func clampProgress(progress int) int {
	if progress < 0 {
		return 0
	}
	if progress > 100 {
		return 100
	}
	return progress
}

func (s *Service) queueOutgoingMessage(peer Peer, message ChatMessage, payload directMessage) error {
	normalized, err := s.appendMessage(message.PeerID, message)
	if err != nil {
		return err
	}

	if normalized.Kind != "text" {
		s.emitTransferProgress(TransferProgress{
			PeerID:    normalized.PeerID,
			MessageID: normalized.ID,
			Progress:  0,
		})
	}

	if peer.Source == "debug" {
		go s.finishDebugSend(peer, normalized)
		return nil
	}

	go s.dispatchOutgoing(peer, normalized, payload)
	return nil
}

func (s *Service) finishDebugSend(peer Peer, message ChatMessage) {
	time.Sleep(120 * time.Millisecond)

	updated := message
	updated.Status = "sent"
	if _, err := s.appendMessage(message.PeerID, updated); err != nil {
		return
	}

	if updated.Kind != "text" {
		s.emitTransferProgress(TransferProgress{
			PeerID:    updated.PeerID,
			MessageID: updated.ID,
			Progress:  100,
		})
	}

	go s.sendDebugReply(peer, updated)
}

func (s *Service) dispatchOutgoing(peer Peer, message ChatMessage, payload directMessage) {
	progressReporter := func(int) {}
	if message.Kind != "text" {
		lastReported := -5
		progressReporter = func(progress int) {
			progress = clampProgress(progress)
			if progress != 100 && progress < lastReported+5 {
				return
			}
			lastReported = progress
			s.emitTransferProgress(TransferProgress{
				PeerID:    message.PeerID,
				MessageID: message.ID,
				Progress:  progress,
			})
		}
	}

	err := s.sendDirect(peer, payload, progressReporter)
	updated := message
	if err != nil {
		updated.Status = "failed"
	} else {
		if message.Kind != "text" {
			progressReporter(100)
		}
		updated.Status = "sent"
	}
	_, _ = s.appendMessage(message.PeerID, updated)
}

func (s *Service) appendMessage(peerID string, message ChatMessage) (ChatMessage, error) {
	normalized, err := s.store.NormalizeMessageMedia(message)
	if err != nil {
		return ChatMessage{}, err
	}

	s.mu.Lock()
	previous := s.conversations[peerID]
	next := upsertConversationMessage(previous, normalized)
	s.conversations[peerID] = next
	if err := s.store.SaveMessage(normalized); err != nil {
		s.conversations[peerID] = previous
		s.mu.Unlock()
		return ChatMessage{}, err
	}
	messages := append([]ChatMessage(nil), next...)
	s.mu.Unlock()

	if s.callbacks.OnMessage != nil {
		s.callbacks.OnMessage(peerID, messages)
	}
	return normalized, nil
}

func (s *Service) sendDirect(peer Peer, payload directMessage, onProgress func(int)) error {
	conn, err := net.DialTimeout("tcp4", fmt.Sprintf("%s:%d", peer.Address, peer.ListenPort), 3*time.Second)
	if err != nil {
		return err
	}
	defer conn.Close()

	_ = conn.SetWriteDeadline(time.Now().Add(12 * time.Second))
	raw, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	raw = append(raw, '\n')

	if onProgress != nil {
		onProgress(0)
	}

	const chunkSize = 64 * 1024
	reader := bytes.NewReader(raw)
	buffer := make([]byte, chunkSize)
	written := 0
	total := len(raw)

	for {
		n, readErr := reader.Read(buffer)
		if n > 0 {
			wn, writeErr := conn.Write(buffer[:n])
			if wn > 0 {
				written += wn
				if onProgress != nil && total > 0 {
					progress := int((float64(written) / float64(total)) * 99)
					onProgress(progress)
				}
			}
			if writeErr != nil {
				return writeErr
			}
			if wn != n {
				return io.ErrShortWrite
			}
		}
		if errors.Is(readErr, io.EOF) {
			break
		}
		if readErr != nil {
			return readErr
		}
	}

	return nil
}

func mergeConversationOnlyPeers(store *Store, peers map[string]Peer) error {
	peerIDs, err := store.LoadConversationPeerIDs()
	if err != nil {
		return err
	}

	for _, peerID := range peerIDs {
		if _, ok := peers[peerID]; ok {
			continue
		}

		peers[peerID] = Peer{
			ID:         peerID,
			Name:       defaultName(peerID),
			Address:    "",
			ListenPort: 0,
			LastSeen:   0,
			Source:     "history",
		}
	}

	return nil
}

func peerCanSend(peer Peer) bool {
	if peer.ID == "" {
		return false
	}

	address := strings.TrimSpace(peer.Address)
	port := peer.ListenPort
	switch peer.Source {
	case "debug":
		return true
	case "lan":
		if address == "" || address == ":0" || port <= 0 {
			return false
		}
		return time.Since(time.UnixMilli(peer.LastSeen)) <= peerExpiry
	default:
		return address != "" && address != ":0" && port > 0
	}
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
		MediaSize:  payload.MediaSize,
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
		if message.MediaSize <= 0 {
			message.MediaSize = int64(payloadSize)
		}
	case "file":
		mediaType, payloadSize, err := validateFileDataURL(payload.Text)
		if err != nil {
			return ChatMessage{}, err
		}
		if payloadSize > maxFileBytes {
			return ChatMessage{}, fmt.Errorf("file too large")
		}
		if strings.TrimSpace(message.MediaType) == "" {
			message.MediaType = mediaType
		}
		if strings.TrimSpace(message.MediaName) == "" {
			message.MediaName = "file"
		}
		if message.MediaSize <= 0 {
			message.MediaSize = int64(payloadSize)
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
	mediaType, decodedLen, err := validateDataURL(value)
	if err != nil {
		return "", 0, err
	}
	if !strings.HasPrefix(mediaType, "image/") {
		return "", 0, errors.New("only image data URLs are supported")
	}
	return mediaType, decodedLen, nil
}

func validateFileDataURL(value string) (string, int, error) {
	return validateDataURL(value)
}

func validateDataURL(value string) (string, int, error) {
	if !strings.HasPrefix(value, "data:image/") {
		if !strings.HasPrefix(value, "data:") {
			return "", 0, errors.New("only data URLs are supported")
		}
	}
	parts := strings.SplitN(value, ",", 2)
	if len(parts) != 2 {
		return "", 0, errors.New("invalid payload")
	}
	header := parts[0]
	if !strings.HasSuffix(header, ";base64") {
		return "", 0, errors.New("payload must be base64 encoded")
	}
	mediaType := strings.TrimPrefix(header, "data:")
	mediaType = strings.TrimSuffix(mediaType, ";base64")
	if mediaType == "" {
		mediaType = "application/octet-stream"
	}

	decodedLen := base64.StdEncoding.DecodedLen(len(parts[1]))
	if decodedLen <= 0 {
		return "", 0, errors.New("invalid payload")
	}

	if _, err := base64.StdEncoding.DecodeString(parts[1]); err != nil {
		return "", 0, errors.New("invalid base64 payload")
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

func (s *Service) emitTransferProgress(progress TransferProgress) {
	if s.callbacks.OnTransferProgress != nil {
		s.callbacks.OnTransferProgress(progress)
	}
}

func (s *Service) sortedPeersLocked() []Peer {
	peers := make([]Peer, 0, len(s.peers))
	for _, peer := range s.peers {
		peers = append(peers, peer)
	}

	sort.Slice(peers, func(i, j int) bool {
		iOnline := peerCanSend(peers[i])
		jOnline := peerCanSend(peers[j])
		if iOnline != jOnline {
			return iOnline
		}
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
