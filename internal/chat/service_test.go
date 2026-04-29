package chat

import (
	"path/filepath"
	"testing"
	"time"
)

func TestMergeConversationOnlyPeersAddsHistoryPlaceholder(t *testing.T) {
	store, err := NewStoreFromDirectory(t.TempDir())
	if err != nil {
		t.Fatalf("NewStoreFromDirectory failed: %v", err)
	}
	defer func() {
		_ = store.db.Close()
	}()

	if err := store.SaveMessage(ChatMessage{
		ID:         "m1",
		PeerID:     "peer-history",
		SenderID:   "self",
		SenderName: "Tester",
		Kind:       "text",
		Text:       "hello",
		Timestamp:  1,
		Direction:  "outbound",
		Status:     "sent",
	}); err != nil {
		t.Fatalf("SaveMessage failed: %v", err)
	}

	peers := map[string]Peer{}
	if err := mergeConversationOnlyPeers(store, peers); err != nil {
		t.Fatalf("mergeConversationOnlyPeers failed: %v", err)
	}

	peer, ok := peers["peer-history"]
	if !ok {
		t.Fatalf("expected placeholder peer to be added")
	}
	if peer.Source != "history" {
		t.Fatalf("expected source history, got %q", peer.Source)
	}
	if peer.Name == "" {
		t.Fatalf("expected placeholder peer name")
	}
}

func TestCleanupPeersKeepsExpiredLanPeerWithHistory(t *testing.T) {
	store, err := NewStoreFromDirectory(t.TempDir())
	if err != nil {
		t.Fatalf("NewStoreFromDirectory failed: %v", err)
	}
	defer func() {
		_ = store.db.Close()
	}()

	peer := Peer{
		ID:         "peer-keep",
		Name:       "Alice",
		Address:    "192.168.1.8",
		ListenPort: 48556,
		LastSeen:   time.Now().Add(-peerExpiry - time.Second).UnixMilli(),
		Source:     "lan",
	}
	if err := store.SavePeer(peer); err != nil {
		t.Fatalf("SavePeer failed: %v", err)
	}
	if err := store.SaveMessage(ChatMessage{
		ID:         "m-history",
		PeerID:     peer.ID,
		SenderID:   "self",
		SenderName: "Tester",
		Kind:       "text",
		Text:       "history",
		Timestamp:  1,
		Direction:  "outbound",
		Status:     "sent",
	}); err != nil {
		t.Fatalf("SaveMessage failed: %v", err)
	}

	service := &Service{
		store: store,
		peers: map[string]Peer{
			peer.ID: peer,
		},
	}

	service.cleanupPeers()

	if _, ok := service.peers[peer.ID]; !ok {
		t.Fatalf("expected expired peer with history to remain")
	}
	state, err := store.Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	if len(state.Peers) != 1 {
		t.Fatalf("expected stored peer to remain, got %d", len(state.Peers))
	}
}

func TestCleanupPeersRemovesExpiredLanPeerWithoutHistory(t *testing.T) {
	store, err := NewStoreFromDirectory(t.TempDir())
	if err != nil {
		t.Fatalf("NewStoreFromDirectory failed: %v", err)
	}
	defer func() {
		_ = store.db.Close()
	}()

	peer := Peer{
		ID:         "peer-drop",
		Name:       "Bob",
		Address:    "192.168.1.9",
		ListenPort: 48557,
		LastSeen:   time.Now().Add(-peerExpiry - time.Second).UnixMilli(),
		Source:     "lan",
	}
	if err := store.SavePeer(peer); err != nil {
		t.Fatalf("SavePeer failed: %v", err)
	}

	service := &Service{
		store: store,
		peers: map[string]Peer{
			peer.ID: peer,
		},
	}

	service.cleanupPeers()

	if _, ok := service.peers[peer.ID]; ok {
		t.Fatalf("expected expired peer without history to be removed")
	}
	state, err := store.Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	if len(state.Peers) != 0 {
		t.Fatalf("expected stored peer to be deleted, got %d", len(state.Peers))
	}
}

func TestServiceDataPathReturnsDirectory(t *testing.T) {
	dataDir := t.TempDir()
	store, err := NewStoreFromDirectory(dataDir)
	if err != nil {
		t.Fatalf("NewStoreFromDirectory failed: %v", err)
	}
	defer func() {
		_ = store.db.Close()
	}()

	service := &Service{store: store}

	if got := service.DataPath(); got != dataDir {
		t.Fatalf("expected data directory %q, got %q", dataDir, got)
	}

	targetDir := filepath.Join(t.TempDir(), "lantalk-data")
	got, err := service.MoveDataFile(targetDir)
	if err != nil {
		t.Fatalf("MoveDataFile failed: %v", err)
	}
	if got != targetDir {
		t.Fatalf("expected moved data directory %q, got %q", targetDir, got)
	}
	if got := service.DataPath(); got != targetDir {
		t.Fatalf("expected service data path %q after move, got %q", targetDir, got)
	}
}
