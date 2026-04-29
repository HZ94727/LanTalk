package chat

import (
	"encoding/base64"
	"os"
	"path/filepath"
	"testing"
)

func TestSaveMessageDeduplicatesVisibleFilePath(t *testing.T) {
	store, err := NewStoreFromDirectory(t.TempDir())
	if err != nil {
		t.Fatalf("NewStoreFromDirectory failed: %v", err)
	}
	defer func() {
		_ = store.db.Close()
	}()

	raw := []byte("%PDF-1.4 test payload")
	dataURL := "data:application/pdf;base64," + base64.StdEncoding.EncodeToString(raw)

	first := ChatMessage{
		ID:         "m1",
		PeerID:     "peer-1",
		SenderID:   "self",
		SenderName: "Tester",
		Kind:       "file",
		Text:       dataURL,
		MediaName:  "1.pdf",
		MediaType:  "application/pdf",
		MediaSize:  int64(len(raw)),
		Timestamp:  1,
		Direction:  "outbound",
		Status:     "sent",
	}
	second := first
	second.ID = "m2"
	second.Timestamp = 2

	if err := store.SaveMessage(first); err != nil {
		t.Fatalf("SaveMessage(first) failed: %v", err)
	}
	if err := store.SaveMessage(second); err != nil {
		t.Fatalf("SaveMessage(second) failed: %v", err)
	}

	page, err := store.LoadConversationPage("peer-1", 0, "", 10)
	if err != nil {
		t.Fatalf("LoadConversationPage failed: %v", err)
	}
	if len(page.Messages) != 2 {
		t.Fatalf("expected 2 messages, got %d", len(page.Messages))
	}
	if page.Messages[0].Text != page.Messages[1].Text {
		t.Fatalf("expected shared file ref, got %q and %q", page.Messages[0].Text, page.Messages[1].Text)
	}

	storedPath, err := store.ResolveStoredFilePath(page.Messages[0].Text)
	if err != nil {
		t.Fatalf("ResolveStoredFilePath failed: %v", err)
	}
	if _, err := os.Stat(storedPath); err != nil {
		t.Fatalf("stored file missing: %v", err)
	}

	if got := countFilesUnder(t, filepath.Join(store.dataDir, "media", "files")); got != 1 {
		t.Fatalf("expected 1 visible stored file, got %d", got)
	}
	if got := countFilesUnder(t, filepath.Join(store.dataDir, "media", fileBlobDirName)); got != 0 {
		t.Fatalf("expected 0 legacy blob files, got %d", got)
	}

	if err := store.RemoveMessage(page.Messages[0]); err != nil {
		t.Fatalf("RemoveMessage(first) failed: %v", err)
	}
	if _, err := os.Stat(storedPath); err != nil {
		t.Fatalf("shared file should remain after first delete: %v", err)
	}

	if err := store.RemoveMessage(page.Messages[1]); err != nil {
		t.Fatalf("RemoveMessage(second) failed: %v", err)
	}
	if _, err := os.Stat(storedPath); !os.IsNotExist(err) {
		t.Fatalf("shared file should be removed after last delete, stat err=%v", err)
	}
}

func TestMigrateStoredFileMessagesCollapsesLegacyHardlinks(t *testing.T) {
	store, err := NewStoreFromDirectory(t.TempDir())
	if err != nil {
		t.Fatalf("NewStoreFromDirectory failed: %v", err)
	}
	defer func() {
		_ = store.db.Close()
	}()

	raw := []byte("%PDF-1.4 legacy payload")
	hash := mediaContentHash(raw)
	_, blobPath := store.fileBlobPathUnlocked(hash)
	if err := os.MkdirAll(filepath.Dir(blobPath), 0o755); err != nil {
		t.Fatalf("MkdirAll(blob) failed: %v", err)
	}
	if err := os.WriteFile(blobPath, raw, 0o644); err != nil {
		t.Fatalf("WriteFile(blob) failed: %v", err)
	}

	oldRef1, oldPath1 := legacyMessageFilePath(store.dataDir, "18", "186ae93e-6a77-4f4a-bb48-3c604dfe13ad", "1.pdf")
	oldRef2, oldPath2 := legacyMessageFilePath(store.dataDir, "26", "265fc562-02ed-4788-8a5e-36c05ac91cd9", "1.pdf")
	if err := os.MkdirAll(filepath.Dir(oldPath1), 0o755); err != nil {
		t.Fatalf("MkdirAll(oldPath1) failed: %v", err)
	}
	if err := os.MkdirAll(filepath.Dir(oldPath2), 0o755); err != nil {
		t.Fatalf("MkdirAll(oldPath2) failed: %v", err)
	}
	if err := os.Link(blobPath, oldPath1); err != nil {
		t.Fatalf("Link(oldPath1) failed: %v", err)
	}
	if err := os.Link(blobPath, oldPath2); err != nil {
		t.Fatalf("Link(oldPath2) failed: %v", err)
	}

	insertLegacyMessage(t, store, ChatMessage{
		ID:         "186ae93e-6a77-4f4a-bb48-3c604dfe13ad",
		PeerID:     "peer-1",
		SenderID:   "self",
		SenderName: "Tester",
		Kind:       "file",
		Text:       oldRef1,
		MediaName:  "1.pdf",
		MediaType:  "application/pdf",
		MediaSize:  int64(len(raw)),
		MediaHash:  hash,
		Timestamp:  1,
		Direction:  "outbound",
		Status:     "sent",
	})
	insertLegacyMessage(t, store, ChatMessage{
		ID:         "265fc562-02ed-4788-8a5e-36c05ac91cd9",
		PeerID:     "peer-1",
		SenderID:   "self",
		SenderName: "Tester",
		Kind:       "file",
		Text:       oldRef2,
		MediaName:  "1.pdf",
		MediaType:  "application/pdf",
		MediaSize:  int64(len(raw)),
		MediaHash:  hash,
		Timestamp:  2,
		Direction:  "outbound",
		Status:     "sent",
	})

	if err := store.migrateStoredFileMessages(); err != nil {
		t.Fatalf("migrateStoredFileMessages failed: %v", err)
	}

	page, err := store.LoadConversationPage("peer-1", 0, "", 10)
	if err != nil {
		t.Fatalf("LoadConversationPage failed: %v", err)
	}
	if len(page.Messages) != 2 {
		t.Fatalf("expected 2 messages after migration, got %d", len(page.Messages))
	}
	if page.Messages[0].Text != page.Messages[1].Text {
		t.Fatalf("expected migrated messages to share ref, got %q and %q", page.Messages[0].Text, page.Messages[1].Text)
	}

	migratedPath, err := store.ResolveStoredFilePath(page.Messages[0].Text)
	if err != nil {
		t.Fatalf("ResolveStoredFilePath failed: %v", err)
	}
	if _, err := os.Stat(migratedPath); err != nil {
		t.Fatalf("migrated file missing: %v", err)
	}
	if _, err := os.Stat(oldPath1); !os.IsNotExist(err) {
		t.Fatalf("expected old message file 1 removed, stat err=%v", err)
	}
	if _, err := os.Stat(oldPath2); !os.IsNotExist(err) {
		t.Fatalf("expected old message file 2 removed, stat err=%v", err)
	}
	if _, err := os.Stat(blobPath); !os.IsNotExist(err) {
		t.Fatalf("expected legacy blob removed, stat err=%v", err)
	}
	if got := countFilesUnder(t, filepath.Join(store.dataDir, "media", "files")); got != 1 {
		t.Fatalf("expected 1 visible stored file after migration, got %d", got)
	}
}

func TestSaveMessageDeduplicatesVisibleImagePath(t *testing.T) {
	store, err := NewStoreFromDirectory(t.TempDir())
	if err != nil {
		t.Fatalf("NewStoreFromDirectory failed: %v", err)
	}
	defer func() {
		_ = store.db.Close()
	}()

	raw := []byte("fake png payload")
	dataURL := "data:image/png;base64," + base64.StdEncoding.EncodeToString(raw)

	first := ChatMessage{
		ID:         "img-1",
		PeerID:     "peer-1",
		SenderID:   "self",
		SenderName: "Tester",
		Kind:       "image",
		Text:       dataURL,
		MediaName:  "avatar.png",
		MediaType:  "image/png",
		MediaSize:  int64(len(raw)),
		Timestamp:  1,
		Direction:  "outbound",
		Status:     "sent",
	}
	second := first
	second.ID = "img-2"
	second.Timestamp = 2

	if err := store.SaveMessage(first); err != nil {
		t.Fatalf("SaveMessage(first) failed: %v", err)
	}
	if err := store.SaveMessage(second); err != nil {
		t.Fatalf("SaveMessage(second) failed: %v", err)
	}

	page, err := store.LoadConversationPage("peer-1", 0, "", 10)
	if err != nil {
		t.Fatalf("LoadConversationPage failed: %v", err)
	}
	if len(page.Messages) != 2 {
		t.Fatalf("expected 2 messages, got %d", len(page.Messages))
	}
	if page.Messages[0].Text != page.Messages[1].Text {
		t.Fatalf("expected shared image ref, got %q and %q", page.Messages[0].Text, page.Messages[1].Text)
	}

	storedPath, err := store.mediaAbsolutePathFromRefUnlocked(page.Messages[0].Text)
	if err != nil {
		t.Fatalf("mediaAbsolutePathFromRefUnlocked failed: %v", err)
	}
	if _, err := os.Stat(storedPath); err != nil {
		t.Fatalf("stored image missing: %v", err)
	}

	if got := countFilesUnder(t, filepath.Join(store.dataDir, "media", imageDirName)); got != 1 {
		t.Fatalf("expected 1 visible stored image, got %d", got)
	}

	if err := store.RemoveMessage(page.Messages[0]); err != nil {
		t.Fatalf("RemoveMessage(first) failed: %v", err)
	}
	if _, err := os.Stat(storedPath); err != nil {
		t.Fatalf("shared image should remain after first delete: %v", err)
	}

	if err := store.RemoveMessage(page.Messages[1]); err != nil {
		t.Fatalf("RemoveMessage(second) failed: %v", err)
	}
	if _, err := os.Stat(storedPath); !os.IsNotExist(err) {
		t.Fatalf("shared image should be removed after last delete, stat err=%v", err)
	}
}

func TestMigrateStoredImageMessagesCollapsesLegacyCopies(t *testing.T) {
	store, err := NewStoreFromDirectory(t.TempDir())
	if err != nil {
		t.Fatalf("NewStoreFromDirectory failed: %v", err)
	}
	defer func() {
		_ = store.db.Close()
	}()

	raw := []byte("legacy image payload")
	oldRef1, oldPath1 := legacyImageFilePath(store.dataDir, "18", "186ae93e-6a77-4f4a-bb48-3c604dfe13ad", ".png")
	oldRef2, oldPath2 := legacyImageFilePath(store.dataDir, "26", "265fc562-02ed-4788-8a5e-36c05ac91cd9", ".png")

	if err := os.MkdirAll(filepath.Dir(oldPath1), 0o755); err != nil {
		t.Fatalf("MkdirAll(oldPath1) failed: %v", err)
	}
	if err := os.MkdirAll(filepath.Dir(oldPath2), 0o755); err != nil {
		t.Fatalf("MkdirAll(oldPath2) failed: %v", err)
	}
	if err := os.WriteFile(oldPath1, raw, 0o644); err != nil {
		t.Fatalf("WriteFile(oldPath1) failed: %v", err)
	}
	if err := os.WriteFile(oldPath2, raw, 0o644); err != nil {
		t.Fatalf("WriteFile(oldPath2) failed: %v", err)
	}

	insertLegacyMessage(t, store, ChatMessage{
		ID:         "186ae93e-6a77-4f4a-bb48-3c604dfe13ad",
		PeerID:     "peer-1",
		SenderID:   "self",
		SenderName: "Tester",
		Kind:       "image",
		Text:       oldRef1,
		MediaName:  "avatar.png",
		MediaType:  "image/png",
		MediaSize:  int64(len(raw)),
		Timestamp:  1,
		Direction:  "outbound",
		Status:     "sent",
	})
	insertLegacyMessage(t, store, ChatMessage{
		ID:         "265fc562-02ed-4788-8a5e-36c05ac91cd9",
		PeerID:     "peer-1",
		SenderID:   "self",
		SenderName: "Tester",
		Kind:       "image",
		Text:       oldRef2,
		MediaName:  "avatar.png",
		MediaType:  "image/png",
		MediaSize:  int64(len(raw)),
		Timestamp:  2,
		Direction:  "outbound",
		Status:     "sent",
	})

	if err := store.migrateStoredImageMessages(); err != nil {
		t.Fatalf("migrateStoredImageMessages failed: %v", err)
	}

	page, err := store.LoadConversationPage("peer-1", 0, "", 10)
	if err != nil {
		t.Fatalf("LoadConversationPage failed: %v", err)
	}
	if len(page.Messages) != 2 {
		t.Fatalf("expected 2 messages after migration, got %d", len(page.Messages))
	}
	if page.Messages[0].Text != page.Messages[1].Text {
		t.Fatalf("expected migrated images to share ref, got %q and %q", page.Messages[0].Text, page.Messages[1].Text)
	}

	migratedPath, err := store.mediaAbsolutePathFromRefUnlocked(page.Messages[0].Text)
	if err != nil {
		t.Fatalf("mediaAbsolutePathFromRefUnlocked failed: %v", err)
	}
	if _, err := os.Stat(migratedPath); err != nil {
		t.Fatalf("migrated image missing: %v", err)
	}
	if _, err := os.Stat(oldPath1); !os.IsNotExist(err) {
		t.Fatalf("expected old image 1 removed, stat err=%v", err)
	}
	if _, err := os.Stat(oldPath2); !os.IsNotExist(err) {
		t.Fatalf("expected old image 2 removed, stat err=%v", err)
	}
	if got := countFilesUnder(t, filepath.Join(store.dataDir, "media", imageDirName)); got != 1 {
		t.Fatalf("expected 1 visible stored image after migration, got %d", got)
	}
}

func TestClearConversationPreservesSharedMediaAcrossPeers(t *testing.T) {
	store, err := NewStoreFromDirectory(t.TempDir())
	if err != nil {
		t.Fatalf("NewStoreFromDirectory failed: %v", err)
	}
	defer func() {
		_ = store.db.Close()
	}()

	fileRaw := []byte("shared pdf payload")
	fileURL := "data:application/pdf;base64," + base64.StdEncoding.EncodeToString(fileRaw)
	imageRaw := []byte("unique image payload")
	imageURL := "data:image/png;base64," + base64.StdEncoding.EncodeToString(imageRaw)

	sharedPeerOne := ChatMessage{
		ID:         "peer-one-file",
		PeerID:     "peer-1",
		SenderID:   "self",
		SenderName: "Tester",
		Kind:       "file",
		Text:       fileURL,
		MediaName:  "shared.pdf",
		MediaType:  "application/pdf",
		MediaSize:  int64(len(fileRaw)),
		Timestamp:  1,
		Direction:  "outbound",
		Status:     "sent",
	}
	sharedPeerTwo := sharedPeerOne
	sharedPeerTwo.ID = "peer-two-file"
	sharedPeerTwo.PeerID = "peer-2"
	sharedPeerTwo.Timestamp = 2
	uniqueImage := ChatMessage{
		ID:         "peer-one-image",
		PeerID:     "peer-1",
		SenderID:   "self",
		SenderName: "Tester",
		Kind:       "image",
		Text:       imageURL,
		MediaName:  "unique.png",
		MediaType:  "image/png",
		MediaSize:  int64(len(imageRaw)),
		Timestamp:  3,
		Direction:  "outbound",
		Status:     "sent",
	}

	if err := store.SaveMessage(sharedPeerOne); err != nil {
		t.Fatalf("SaveMessage(sharedPeerOne) failed: %v", err)
	}
	if err := store.SaveMessage(sharedPeerTwo); err != nil {
		t.Fatalf("SaveMessage(sharedPeerTwo) failed: %v", err)
	}
	if err := store.SaveMessage(uniqueImage); err != nil {
		t.Fatalf("SaveMessage(uniqueImage) failed: %v", err)
	}

	peerOnePage, err := store.LoadConversationPage("peer-1", 0, "", 10)
	if err != nil {
		t.Fatalf("LoadConversationPage(peer-1) failed: %v", err)
	}
	peerTwoPage, err := store.LoadConversationPage("peer-2", 0, "", 10)
	if err != nil {
		t.Fatalf("LoadConversationPage(peer-2) failed: %v", err)
	}
	sharedPath, err := store.ResolveStoredFilePath(peerOnePage.Messages[0].Text)
	if err != nil {
		t.Fatalf("ResolveStoredFilePath failed: %v", err)
	}
	uniqueImagePath, err := store.mediaAbsolutePathFromRefUnlocked(peerOnePage.Messages[1].Text)
	if err != nil {
		t.Fatalf("mediaAbsolutePathFromRefUnlocked failed: %v", err)
	}

	removedCount, err := store.ClearConversation("peer-1")
	if err != nil {
		t.Fatalf("ClearConversation failed: %v", err)
	}
	if removedCount != 2 {
		t.Fatalf("expected 2 removed messages, got %d", removedCount)
	}

	peerOneAfter, err := store.LoadConversationPage("peer-1", 0, "", 10)
	if err != nil {
		t.Fatalf("LoadConversationPage(peer-1 after clear) failed: %v", err)
	}
	if len(peerOneAfter.Messages) != 0 {
		t.Fatalf("expected peer-1 conversation empty, got %d messages", len(peerOneAfter.Messages))
	}
	peerTwoAfter, err := store.LoadConversationPage("peer-2", 0, "", 10)
	if err != nil {
		t.Fatalf("LoadConversationPage(peer-2 after clear) failed: %v", err)
	}
	if len(peerTwoAfter.Messages) != 1 {
		t.Fatalf("expected peer-2 conversation to retain 1 message, got %d", len(peerTwoAfter.Messages))
	}
	if peerTwoAfter.Messages[0].Text != peerTwoPage.Messages[0].Text {
		t.Fatalf("expected shared ref unchanged for peer-2")
	}

	if _, err := os.Stat(sharedPath); err != nil {
		t.Fatalf("shared file should remain referenced by peer-2: %v", err)
	}
	if _, err := os.Stat(uniqueImagePath); !os.IsNotExist(err) {
		t.Fatalf("unique image should be removed after clear, stat err=%v", err)
	}
}

func TestCleanupUnusedMediaRemovesOnlyOrphanedFiles(t *testing.T) {
	store, err := NewStoreFromDirectory(t.TempDir())
	if err != nil {
		t.Fatalf("NewStoreFromDirectory failed: %v", err)
	}
	defer func() {
		_ = store.db.Close()
	}()

	fileRaw := []byte("kept file payload")
	fileURL := "data:application/pdf;base64," + base64.StdEncoding.EncodeToString(fileRaw)
	imageRaw := []byte("kept image payload")
	imageURL := "data:image/png;base64," + base64.StdEncoding.EncodeToString(imageRaw)

	if err := store.SaveMessage(ChatMessage{
		ID:         "kept-file",
		PeerID:     "peer-1",
		SenderID:   "self",
		SenderName: "Tester",
		Kind:       "file",
		Text:       fileURL,
		MediaName:  "kept.pdf",
		MediaType:  "application/pdf",
		MediaSize:  int64(len(fileRaw)),
		Timestamp:  1,
		Direction:  "outbound",
		Status:     "sent",
	}); err != nil {
		t.Fatalf("SaveMessage(kept file) failed: %v", err)
	}
	if err := store.SaveMessage(ChatMessage{
		ID:         "kept-image",
		PeerID:     "peer-1",
		SenderID:   "self",
		SenderName: "Tester",
		Kind:       "image",
		Text:       imageURL,
		MediaName:  "kept.png",
		MediaType:  "image/png",
		MediaSize:  int64(len(imageRaw)),
		Timestamp:  2,
		Direction:  "outbound",
		Status:     "sent",
	}); err != nil {
		t.Fatalf("SaveMessage(kept image) failed: %v", err)
	}

	page, err := store.LoadConversationPage("peer-1", 0, "", 10)
	if err != nil {
		t.Fatalf("LoadConversationPage failed: %v", err)
	}

	keptFilePath, err := store.ResolveStoredFilePath(page.Messages[0].Text)
	if err != nil {
		t.Fatalf("ResolveStoredFilePath failed: %v", err)
	}
	keptImagePath, err := store.mediaAbsolutePathFromRefUnlocked(page.Messages[1].Text)
	if err != nil {
		t.Fatalf("mediaAbsolutePathFromRefUnlocked failed: %v", err)
	}

	orphanFilePath := filepath.Join(store.dataDir, "media", "files", "zz", "orphan", "lost.pdf")
	orphanImagePath := filepath.Join(store.dataDir, "media", imageDirName, "yy", "orphan", "lost.png")
	if err := os.MkdirAll(filepath.Dir(orphanFilePath), 0o755); err != nil {
		t.Fatalf("MkdirAll(orphan file) failed: %v", err)
	}
	if err := os.MkdirAll(filepath.Dir(orphanImagePath), 0o755); err != nil {
		t.Fatalf("MkdirAll(orphan image) failed: %v", err)
	}
	if err := os.WriteFile(orphanFilePath, []byte("orphan file"), 0o644); err != nil {
		t.Fatalf("WriteFile(orphan file) failed: %v", err)
	}
	if err := os.WriteFile(orphanImagePath, []byte("orphan image"), 0o644); err != nil {
		t.Fatalf("WriteFile(orphan image) failed: %v", err)
	}

	result, err := store.CleanupUnusedMedia()
	if err != nil {
		t.Fatalf("CleanupUnusedMedia failed: %v", err)
	}
	if result.RemovedFiles != 2 {
		t.Fatalf("expected 2 removed orphaned files, got %d", result.RemovedFiles)
	}
	if result.ReclaimedBytes <= 0 {
		t.Fatalf("expected reclaimed bytes to be positive, got %d", result.ReclaimedBytes)
	}

	if _, err := os.Stat(keptFilePath); err != nil {
		t.Fatalf("kept file should remain after cleanup: %v", err)
	}
	if _, err := os.Stat(keptImagePath); err != nil {
		t.Fatalf("kept image should remain after cleanup: %v", err)
	}
	if _, err := os.Stat(orphanFilePath); !os.IsNotExist(err) {
		t.Fatalf("orphan file should be removed, stat err=%v", err)
	}
	if _, err := os.Stat(orphanImagePath); !os.IsNotExist(err) {
		t.Fatalf("orphan image should be removed, stat err=%v", err)
	}
}

func insertLegacyMessage(t *testing.T, store *Store, message ChatMessage) {
	t.Helper()

	encrypted, err := store.encryptMessage(message)
	if err != nil {
		t.Fatalf("encryptMessage failed: %v", err)
	}
	if _, err := store.db.Exec(`
		INSERT INTO messages(id, peer_id, sender_id, sender_name, kind, text, media_name, media_type, media_size, media_hash, timestamp, direction, status)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, encrypted.ID, encrypted.PeerID, encrypted.SenderID, encrypted.SenderName, encrypted.Kind, encrypted.Text, encrypted.MediaName, encrypted.MediaType, encrypted.MediaSize, encrypted.MediaHash, encrypted.Timestamp, encrypted.Direction, encrypted.Status); err != nil {
		t.Fatalf("insert legacy message failed: %v", err)
	}
}

func legacyMessageFilePath(dataDir, shard, messageID, fileName string) (string, string) {
	relative := filepath.Join("media", "files", shard, messageID, fileName)
	return fileReference(relative), filepath.Join(dataDir, relative)
}

func legacyImageFilePath(dataDir, shard, messageID, ext string) (string, string) {
	relative := filepath.Join("media", shard, messageID+ext)
	return mediaReference(relative), filepath.Join(dataDir, relative)
}

func countFilesUnder(t *testing.T, root string) int {
	t.Helper()

	count := 0
	err := filepath.Walk(root, func(_ string, info os.FileInfo, walkErr error) error {
		if walkErr != nil {
			if os.IsNotExist(walkErr) {
				return nil
			}
			return walkErr
		}
		if info != nil && !info.IsDir() {
			count++
		}
		return nil
	})
	if err != nil {
		t.Fatalf("Walk(%s) failed: %v", root, err)
	}
	return count
}
