package chat

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"LanTalk/internal/mediautil"
)

const imageRefPrefix = "media:v1:"
const fileRefPrefix = "file:v1:"
const maxMediaStorageBytes int64 = 1024 * 1024 * 1024
const imageDirName = "images"
const fileBlobDirName = "blob-files"

func isImageDataURL(value string) bool {
	return strings.HasPrefix(strings.TrimSpace(value), "data:image/")
}

func isFileDataURL(value string) bool {
	return strings.HasPrefix(strings.TrimSpace(value), "data:")
}

func isStoredImageRef(value string) bool {
	return strings.HasPrefix(strings.TrimSpace(value), imageRefPrefix)
}

func isStoredFileRef(value string) bool {
	return strings.HasPrefix(strings.TrimSpace(value), fileRefPrefix)
}

func mediaReference(relativePath string) string {
	return imageRefPrefix + filepath.ToSlash(relativePath)
}

func fileReference(relativePath string) string {
	return fileRefPrefix + filepath.ToSlash(relativePath)
}

func mediaPathFromRef(ref string) string {
	return strings.TrimPrefix(strings.TrimSpace(ref), imageRefPrefix)
}

func filePathFromRef(ref string) string {
	return strings.TrimPrefix(strings.TrimSpace(ref), fileRefPrefix)
}

func mediaTypeForPath(path string) string {
	switch strings.ToLower(filepath.Ext(path)) {
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".webp":
		return "image/webp"
	case ".gif":
		return "image/gif"
	case ".bmp":
		return "image/bmp"
	case ".ico":
		return "image/x-icon"
	default:
		return "image/png"
	}
}

func fileMediaTypeForPath(path string) string {
	return mediautil.MediaTypeForFileName(path)
}

func buildImageDataURL(raw []byte, mediaType string) string {
	return "data:" + mediaType + ";base64," + base64.StdEncoding.EncodeToString(raw)
}

func mediaContentHash(raw []byte) string {
	sum := sha256.Sum256(raw)
	return hex.EncodeToString(sum[:])
}

func (s *Store) ensureMediaUsageUnlocked() error {
	if s.mediaUsageKnown {
		return nil
	}

	var total int64
	err := filepath.Walk(s.mediaDir, func(path string, info os.FileInfo, walkErr error) error {
		if walkErr != nil {
			if errors.Is(walkErr, os.ErrNotExist) {
				return nil
			}
			return walkErr
		}
		if info == nil || info.IsDir() {
			return nil
		}
		total += info.Size()
		return nil
	})
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}

	s.mediaUsageBytes = total
	s.mediaUsageKnown = true
	return nil
}

func (s *Store) cleanupEmptyParentDirsUnlocked(path, stopDir string) {
	stopDir = filepath.Clean(stopDir)
	current := filepath.Dir(filepath.Clean(path))

	for current != stopDir {
		if err := os.Remove(current); err != nil {
			break
		}
		next := filepath.Dir(current)
		if next == current {
			break
		}
		current = next
	}
}

func (s *Store) reserveMediaSpaceUnlocked(size int64) error {
	if size <= 0 {
		return nil
	}
	if err := s.ensureMediaUsageUnlocked(); err != nil {
		return err
	}
	if s.mediaUsageBytes+size > maxMediaStorageBytes {
		return fmt.Errorf("media storage limit reached")
	}
	s.mediaUsageBytes += size
	return nil
}

func (s *Store) rollbackReservedMediaSpaceUnlocked(size int64) {
	if !s.mediaUsageKnown || size <= 0 {
		return
	}
	s.mediaUsageBytes -= size
	if s.mediaUsageBytes < 0 {
		s.mediaUsageBytes = 0
	}
}

func (s *Store) removeMediaFileUnlocked(absolutePath string) error {
	var size int64
	info, err := os.Stat(absolutePath)
	if err == nil && info != nil && !info.IsDir() {
		size = info.Size()
	} else if err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}

	if err := os.Remove(absolutePath); err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}

	if errors.Is(err, os.ErrNotExist) {
		s.mediaUsageKnown = false
		s.mediaUsageBytes = 0
	} else if s.mediaUsageKnown {
		s.mediaUsageBytes -= size
		if s.mediaUsageBytes < 0 {
			s.mediaUsageBytes = 0
		}
	}

	s.cleanupEmptyParentDirsUnlocked(absolutePath, s.mediaDir)
	return nil
}

func (s *Store) removeFileLinkUnlocked(absolutePath string) error {
	if err := os.Remove(absolutePath); err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}

	s.cleanupEmptyParentDirsUnlocked(absolutePath, s.mediaDir)
	return nil
}

func (s *Store) NormalizeMessageMedia(message ChatMessage) (ChatMessage, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.normalizeMessageMediaUnlocked(message)
}

func (s *Store) normalizeMessageMediaUnlocked(message ChatMessage) (ChatMessage, error) {
	if message.Kind != "image" && message.Kind != "file" {
		return message, nil
	}

	normalized := message
	normalized.MediaType = strings.TrimSpace(normalized.MediaType)
	normalized.MediaName = strings.TrimSpace(normalized.MediaName)
	if normalized.MediaSize < 0 {
		normalized.MediaSize = 0
	}

	if normalized.Kind == "file" {
		if isStoredFileRef(normalized.Text) || !isFileDataURL(normalized.Text) {
			if normalized.MediaType == "" {
				normalized.MediaType = fileMediaTypeForPath(normalized.MediaName)
			}
			if normalized.MediaName == "" {
				normalized.MediaName = mediautil.NormalizeFileName("", normalized.MediaType)
			}
			return normalized, nil
		}

		raw, mediaType, err := mediautil.DecodeDataURL(normalized.Text)
		if err != nil {
			return ChatMessage{}, err
		}
		if normalized.MediaType == "" {
			normalized.MediaType = mediaType
		}
		normalized.MediaName = mediautil.NormalizeFileName(normalized.MediaName, normalized.MediaType)
		if normalized.MediaSize <= 0 {
			normalized.MediaSize = int64(len(raw))
		}

		ref, hash, err := s.storeFilePayloadUnlocked(raw, normalized.MediaType, normalized.MediaName)
		if err != nil {
			return ChatMessage{}, err
		}
		normalized.Text = ref
		normalized.MediaHash = hash
		return normalized, nil
	}

	if isStoredImageRef(normalized.Text) || !isImageDataURL(normalized.Text) {
		if normalized.MediaType == "" {
			normalized.MediaType = mediaTypeForPath(normalized.Text)
		}
		if normalized.MediaName == "" {
			normalized.MediaName = mediautil.NormalizeImageFileName("", normalized.MediaType)
		}
		return normalized, nil
	}

	raw, mediaType, err := mediautil.DecodeImageDataURL(normalized.Text)
	if err != nil {
		return ChatMessage{}, err
	}
	if normalized.MediaType == "" {
		normalized.MediaType = mediaType
	}
	normalized.MediaName = mediautil.NormalizeImageFileName(normalized.MediaName, normalized.MediaType)

	if normalized.MediaSize <= 0 {
		normalized.MediaSize = int64(len(raw))
	}

	ref, hash, err := s.storeImagePayloadUnlocked(raw, normalized.MediaType, normalized.MediaName)
	if err != nil {
		return ChatMessage{}, err
	}
	normalized.Text = ref
	normalized.MediaHash = hash
	return normalized, nil
}

func (s *Store) storeImagePayloadUnlocked(raw []byte, mediaType, fileName string) (string, string, error) {
	hash := mediaContentHash(raw)
	relativePath, absolutePath, err := s.imageStoragePathUnlocked(hash, fileName, mediaType)
	if err != nil {
		return "", "", err
	}

	mediaType = strings.TrimSpace(mediaType)
	if mediaType == "" {
		mediaType = mediaTypeForPath(fileName)
	}
	reservedSize := int64(len(raw))

	if _, err := os.Stat(absolutePath); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return "", "", err
		}
		if err := s.reserveMediaSpaceUnlocked(reservedSize); err != nil {
			return "", "", err
		}

		if err := os.MkdirAll(filepath.Dir(absolutePath), 0o755); err != nil {
			s.rollbackReservedMediaSpaceUnlocked(reservedSize)
			return "", "", err
		}

		tempPath := absolutePath + ".tmp"
		if err := os.WriteFile(tempPath, raw, 0o644); err != nil {
			s.rollbackReservedMediaSpaceUnlocked(reservedSize)
			return "", "", err
		}
		if err := os.Rename(tempPath, absolutePath); err != nil {
			_ = os.Remove(tempPath)
			s.rollbackReservedMediaSpaceUnlocked(reservedSize)
			return "", "", err
		}
	}

	return mediaReference(relativePath), hash, nil
}

func (s *Store) imageStoragePathUnlocked(hash, fileName, mediaType string) (string, string, error) {
	hash = strings.TrimSpace(strings.ToLower(hash))
	if len(hash) < 2 {
		return "", "", errors.New("invalid media hash")
	}

	mediaType = strings.TrimSpace(mediaType)
	if mediaType == "" {
		mediaType = mediaTypeForPath(fileName)
	}
	relativeDir := filepath.Join("media", imageDirName, hash[:2], hash)
	absoluteDir := filepath.Join(s.dataDir, relativeDir)

	entries, err := os.ReadDir(absoluteDir)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return "", "", err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		relativePath := filepath.Join(relativeDir, entry.Name())
		return relativePath, filepath.Join(s.dataDir, relativePath), nil
	}

	fileName = mediautil.NormalizeImageFileName(fileName, mediaType)
	relativePath := filepath.Join(relativeDir, fileName)
	return relativePath, filepath.Join(s.dataDir, relativePath), nil
}

func (s *Store) storeFilePayloadUnlocked(raw []byte, mediaType, fileName string) (string, string, error) {
	hash := mediaContentHash(raw)
	relativePath, absolutePath, err := s.fileStoragePathUnlocked(hash, fileName, mediaType)
	if err != nil {
		return "", "", err
	}

	if _, err := os.Stat(absolutePath); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return "", "", err
		}
		reservedSize := int64(len(raw))
		if err := s.reserveMediaSpaceUnlocked(reservedSize); err != nil {
			return "", "", err
		}
		if err := os.MkdirAll(filepath.Dir(absolutePath), 0o755); err != nil {
			s.rollbackReservedMediaSpaceUnlocked(reservedSize)
			return "", "", err
		}
		tempPath := absolutePath + ".tmp"
		if err := os.WriteFile(tempPath, raw, 0o644); err != nil {
			s.rollbackReservedMediaSpaceUnlocked(reservedSize)
			return "", "", err
		}
		if err := os.Rename(tempPath, absolutePath); err != nil {
			_ = os.Remove(tempPath)
			s.rollbackReservedMediaSpaceUnlocked(reservedSize)
			return "", "", err
		}
	}

	return fileReference(relativePath), hash, nil
}

func (s *Store) fileStoragePathUnlocked(hash, fileName, mediaType string) (string, string, error) {
	hash = strings.TrimSpace(strings.ToLower(hash))
	if len(hash) < 2 {
		return "", "", errors.New("invalid media hash")
	}

	mediaType = strings.TrimSpace(mediaType)
	if mediaType == "" {
		mediaType = fileMediaTypeForPath(fileName)
	}
	relativeDir := filepath.Join("media", "files", hash[:2], hash)
	absoluteDir := filepath.Join(s.dataDir, relativeDir)

	entries, err := os.ReadDir(absoluteDir)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return "", "", err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		relativePath := filepath.Join(relativeDir, entry.Name())
		return relativePath, filepath.Join(s.dataDir, relativePath), nil
	}

	fileName = mediautil.NormalizeFileName(fileName, mediaType)
	relativePath := filepath.Join(relativeDir, fileName)
	return relativePath, filepath.Join(s.dataDir, relativePath), nil
}

func (s *Store) fileBlobPathUnlocked(hash string) (string, string) {
	hash = strings.TrimSpace(strings.ToLower(hash))
	if len(hash) < 2 {
		hash = "00" + hash
	}
	relativePath := filepath.Join("media", fileBlobDirName, hash[:2], hash)
	return relativePath, filepath.Join(s.dataDir, relativePath)
}

func (s *Store) mediaAbsolutePathFromRefUnlocked(ref string) (string, error) {
	if !isStoredImageRef(ref) {
		return "", errors.New("invalid image reference")
	}

	relativePath := filepath.Clean(filepath.FromSlash(mediaPathFromRef(ref)))
	if relativePath == "." || strings.HasPrefix(relativePath, "..") {
		return "", errors.New("invalid image reference path")
	}

	absolutePath := filepath.Join(s.dataDir, relativePath)
	mediaRoot := filepath.Clean(s.mediaDir)
	candidate := filepath.Clean(absolutePath)
	if candidate != mediaRoot && !strings.HasPrefix(candidate, mediaRoot+string(os.PathSeparator)) {
		return "", errors.New("image reference escapes media directory")
	}

	return candidate, nil
}

func (s *Store) fileAbsolutePathFromRefUnlocked(ref string) (string, error) {
	if !isStoredFileRef(ref) {
		return "", errors.New("invalid file reference")
	}

	relativePath := filepath.Clean(filepath.FromSlash(filePathFromRef(ref)))
	if relativePath == "." || strings.HasPrefix(relativePath, "..") {
		return "", errors.New("invalid file reference path")
	}

	absolutePath := filepath.Join(s.dataDir, relativePath)
	mediaRoot := filepath.Clean(s.mediaDir)
	candidate := filepath.Clean(absolutePath)
	if candidate != mediaRoot && !strings.HasPrefix(candidate, mediaRoot+string(os.PathSeparator)) {
		return "", errors.New("file reference escapes media directory")
	}

	return candidate, nil
}

func (s *Store) ResolveImagePayload(value, mediaType string) ([]byte, string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.resolveImagePayloadUnlocked(value, mediaType)
}

func (s *Store) resolveImagePayloadUnlocked(value, mediaType string) ([]byte, string, error) {
	value = strings.TrimSpace(value)
	mediaType = strings.TrimSpace(mediaType)

	if isImageDataURL(value) {
		return mediautil.DecodeImageDataURL(value)
	}
	if !isStoredImageRef(value) {
		return nil, "", errors.New("unsupported image source")
	}

	absolutePath, err := s.mediaAbsolutePathFromRefUnlocked(value)
	if err != nil {
		return nil, "", err
	}
	raw, err := os.ReadFile(absolutePath)
	if err != nil {
		return nil, "", err
	}

	if mediaType == "" {
		mediaType = mediaTypeForPath(absolutePath)
	}
	return raw, mediaType, nil
}

func (s *Store) ResolveImageDataURL(value, mediaType string) (string, error) {
	raw, resolvedMediaType, err := s.ResolveImagePayload(value, mediaType)
	if err != nil {
		return "", err
	}
	return buildImageDataURL(raw, resolvedMediaType), nil
}

func (s *Store) ResolveFilePayload(value, mediaType string) ([]byte, string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.resolveFilePayloadUnlocked(value, mediaType)
}

func (s *Store) resolveFilePayloadUnlocked(value, mediaType string) ([]byte, string, error) {
	value = strings.TrimSpace(value)
	mediaType = strings.TrimSpace(mediaType)

	if isFileDataURL(value) {
		return mediautil.DecodeDataURL(value)
	}
	if !isStoredFileRef(value) {
		return nil, "", errors.New("unsupported file source")
	}

	absolutePath, err := s.fileAbsolutePathFromRefUnlocked(value)
	if err != nil {
		return nil, "", err
	}
	raw, err := os.ReadFile(absolutePath)
	if err != nil {
		return nil, "", err
	}

	if mediaType == "" {
		mediaType = fileMediaTypeForPath(absolutePath)
	}
	return raw, mediaType, nil
}

func (s *Store) ResolveStoredFilePath(value string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.fileAbsolutePathFromRefUnlocked(value)
}

func (s *Store) ResolveStoredImagePath(value string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.mediaAbsolutePathFromRefUnlocked(value)
}

func (s *Store) IsStoredFileAvailable(value string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	value = strings.TrimSpace(value)
	if value == "" {
		return false
	}
	if isFileDataURL(value) {
		return true
	}
	absolutePath, err := s.fileAbsolutePathFromRefUnlocked(value)
	if err != nil {
		return false
	}
	if _, err := os.Stat(absolutePath); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			s.mediaUsageKnown = false
			s.mediaUsageBytes = 0
		}
		return false
	}
	return true
}

func (s *Store) ResolveFileDataURL(value, mediaType string) (string, error) {
	raw, resolvedMediaType, err := s.ResolveFilePayload(value, mediaType)
	if err != nil {
		return "", err
	}
	return buildImageDataURL(raw, resolvedMediaType), nil
}

func (s *Store) mediaUsageBytesUnlocked() (int64, error) {
	if err := s.ensureMediaUsageUnlocked(); err != nil {
		return 0, err
	}
	return s.mediaUsageBytes, nil
}

func (s *Store) StorageStats() (StorageStats, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var stats StorageStats
	if err := s.db.QueryRow(`SELECT COUNT(1) FROM messages`).Scan(&stats.MessageCount); err != nil {
		return StorageStats{}, err
	}

	mediaBytes, err := s.mediaUsageBytesUnlocked()
	if err != nil {
		return StorageStats{}, err
	}
	stats.MediaBytes = mediaBytes
	return stats, nil
}

func (s *Store) loadConversationMessagesUnlocked(peerID string) ([]ChatMessage, error) {
	rows, err := s.db.Query(`
		SELECT id, peer_id, sender_id, sender_name, kind, text, media_name, media_type, media_size, media_hash, timestamp, direction, status
		FROM messages
		WHERE peer_id = ?
		ORDER BY timestamp, id
	`, peerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	messages := make([]ChatMessage, 0, 32)
	for rows.Next() {
		var message ChatMessage
		if err := rows.Scan(
			&message.ID,
			&message.PeerID,
			&message.SenderID,
			&message.SenderName,
			&message.Kind,
			&message.Text,
			&message.MediaName,
			&message.MediaType,
			&message.MediaSize,
			&message.MediaHash,
			&message.Timestamp,
			&message.Direction,
			&message.Status,
		); err != nil {
			return nil, err
		}
		message, err = s.decryptMessage(message)
		if err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return messages, nil
}

func (s *Store) ClearConversation(peerID string) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	peerID = strings.TrimSpace(peerID)
	if peerID == "" {
		return 0, errors.New("peer id is required")
	}

	messages, err := s.loadConversationMessagesUnlocked(peerID)
	if err != nil {
		return 0, err
	}
	if len(messages) == 0 {
		return 0, nil
	}

	if _, err := s.db.Exec(`DELETE FROM messages WHERE peer_id = ?`, peerID); err != nil {
		return 0, err
	}

	for _, message := range messages {
		if message.Kind == "image" && isStoredImageRef(message.Text) {
			_ = s.deleteImageRefIfUnreferencedUnlocked(message.Text, message.MediaHash, message.ID)
		}
		if message.Kind == "file" && isStoredFileRef(message.Text) {
			_ = s.deleteFileRefIfUnreferencedUnlocked(message.Text, message.MediaHash, message.ID)
		}
	}

	return len(messages), nil
}

func (s *Store) CleanupUnusedMedia() (MediaCleanupResult, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var result MediaCleanupResult
	referenced, err := s.referencedMediaPathsUnlocked()
	if err != nil {
		return result, err
	}

	orphanedFiles := make([]string, 0, 16)
	err = filepath.Walk(s.mediaDir, func(path string, info os.FileInfo, walkErr error) error {
		if walkErr != nil {
			if errors.Is(walkErr, os.ErrNotExist) {
				return nil
			}
			return walkErr
		}
		if info == nil || info.IsDir() {
			return nil
		}
		candidate := filepath.Clean(path)
		if _, ok := referenced[candidate]; ok {
			return nil
		}
		orphanedFiles = append(orphanedFiles, candidate)
		return nil
	})
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return result, err
	}

	for _, orphanedPath := range orphanedFiles {
		info, statErr := os.Stat(orphanedPath)
		if statErr != nil {
			if errors.Is(statErr, os.ErrNotExist) {
				continue
			}
			return result, statErr
		}
		if err := os.Remove(orphanedPath); err != nil && !errors.Is(err, os.ErrNotExist) {
			return result, err
		}
		result.RemovedFiles++
		result.ReclaimedBytes += info.Size()
		s.cleanupEmptyParentDirsUnlocked(orphanedPath, s.mediaDir)
	}

	s.mediaUsageKnown = false
	s.mediaUsageBytes = 0
	remainingBytes, err := s.mediaUsageBytesUnlocked()
	if err != nil {
		return result, err
	}
	result.MediaBytes = remainingBytes
	return result, nil
}

func (s *Store) referencedMediaPathsUnlocked() (map[string]struct{}, error) {
	rows, err := s.db.Query(`
		SELECT kind, text
		FROM messages
		WHERE kind IN ('image', 'file')
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	referenced := make(map[string]struct{})
	for rows.Next() {
		var (
			kind          string
			encryptedText string
		)
		if err := rows.Scan(&kind, &encryptedText); err != nil {
			return nil, err
		}

		text, err := s.crypto.DecryptString(encryptedText)
		if err != nil {
			return nil, err
		}

		var absolutePath string
		switch kind {
		case "image":
			if !isStoredImageRef(text) {
				continue
			}
			absolutePath, err = s.mediaAbsolutePathFromRefUnlocked(text)
		case "file":
			if !isStoredFileRef(text) {
				continue
			}
			absolutePath, err = s.fileAbsolutePathFromRefUnlocked(text)
		default:
			continue
		}
		if err != nil {
			continue
		}

		referenced[filepath.Clean(absolutePath)] = struct{}{}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return referenced, nil
}

func (s *Store) RemoveMessage(message ChatMessage) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, err := s.db.Exec(`DELETE FROM messages WHERE id = ?`, message.ID); err != nil {
		return err
	}

	if message.Kind == "image" && isStoredImageRef(message.Text) {
		if err := s.deleteImageRefIfUnreferencedUnlocked(message.Text, message.MediaHash, message.ID); err != nil {
			return err
		}
	}
	if message.Kind == "file" && isStoredFileRef(message.Text) {
		if err := s.deleteFileRefIfUnreferencedUnlocked(message.Text, message.MediaHash, message.ID); err != nil {
			return err
		}
	}

	return nil
}

func (s *Store) deleteImageRefIfUnreferencedUnlocked(ref, mediaHash, excludingMessageID string) error {
	mediaHash = strings.TrimSpace(mediaHash)
	if mediaHash != "" {
		var count int
		if err := s.db.QueryRow(`
			SELECT COUNT(1)
			FROM messages
			WHERE kind = 'image' AND media_hash = ?
		`, mediaHash).Scan(&count); err != nil {
			return err
		}
		if count > 0 {
			return nil
		}

		absolutePath, err := s.mediaAbsolutePathFromRefUnlocked(ref)
		if err != nil {
			return nil
		}
		return s.removeMediaFileUnlocked(absolutePath)
	}

	rows, err := s.db.Query(`
		SELECT id, text
		FROM messages
		WHERE kind = 'image' AND id <> ?
	`, excludingMessageID)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var messageID string
		var encryptedText string
		if err := rows.Scan(&messageID, &encryptedText); err != nil {
			return err
		}
		decryptedText, err := s.crypto.DecryptString(encryptedText)
		if err != nil {
			return err
		}
		if decryptedText == ref {
			return nil
		}
	}
	if err := rows.Err(); err != nil {
		return err
	}

	absolutePath, err := s.mediaAbsolutePathFromRefUnlocked(ref)
	if err != nil {
		return nil
	}
	return s.removeMediaFileUnlocked(absolutePath)
}

func (s *Store) deleteFileRefIfUnreferencedUnlocked(ref, mediaHash, excludingMessageID string) error {
	mediaHash = strings.TrimSpace(mediaHash)
	if mediaHash == "" {
		absolutePath, err := s.fileAbsolutePathFromRefUnlocked(ref)
		if err != nil {
			return nil
		}
		return s.removeMediaFileUnlocked(absolutePath)
	}

	var count int
	if err := s.db.QueryRow(`
		SELECT COUNT(1)
		FROM messages
		WHERE kind = 'file' AND media_hash = ?
	`, mediaHash).Scan(&count); err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	absolutePath, err := s.fileAbsolutePathFromRefUnlocked(ref)
	if err != nil {
		return nil
	}
	return s.removeMediaFileUnlocked(absolutePath)
}

func (s *Store) migrateInlineImageMessages() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	rows, err := s.db.Query(`
		SELECT id, peer_id, sender_id, sender_name, kind, text, media_name, media_type, media_size, media_hash, timestamp, direction, status
		FROM messages
		WHERE kind = 'image'
	`)
	if err != nil {
		return err
	}
	defer rows.Close()

	var pending []ChatMessage
	for rows.Next() {
		var message ChatMessage
		if err := rows.Scan(
			&message.ID,
			&message.PeerID,
			&message.SenderID,
			&message.SenderName,
			&message.Kind,
			&message.Text,
			&message.MediaName,
			&message.MediaType,
			&message.MediaSize,
			&message.MediaHash,
			&message.Timestamp,
			&message.Direction,
			&message.Status,
		); err != nil {
			return err
		}

		message, err = s.decryptMessage(message)
		if err != nil {
			return err
		}
		if !isImageDataURL(message.Text) {
			continue
		}
		message, err = s.normalizeMessageMediaUnlocked(message)
		if err != nil {
			return err
		}
		pending = append(pending, message)
	}
	if err := rows.Err(); err != nil {
		return err
	}

	for _, message := range pending {
		encrypted, err := s.encryptMessage(message)
		if err != nil {
			return err
		}
		if _, err := s.db.Exec(`
			UPDATE messages
			SET sender_name = ?, text = ?, media_name = ?, media_type = ?, media_size = ?, media_hash = ?
			WHERE id = ?
		`, encrypted.SenderName, encrypted.Text, encrypted.MediaName, encrypted.MediaType, encrypted.MediaSize, encrypted.MediaHash, encrypted.ID); err != nil {
			return err
		}
	}

	return nil
}

func (s *Store) migrateStoredImageMessages() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	rows, err := s.db.Query(`
		SELECT id, peer_id, sender_id, sender_name, kind, text, media_name, media_type, media_size, media_hash, timestamp, direction, status
		FROM messages
		WHERE kind = 'image'
	`)
	if err != nil {
		return err
	}
	defer rows.Close()

	var pending []ChatMessage
	for rows.Next() {
		var message ChatMessage
		if err := rows.Scan(
			&message.ID,
			&message.PeerID,
			&message.SenderID,
			&message.SenderName,
			&message.Kind,
			&message.Text,
			&message.MediaName,
			&message.MediaType,
			&message.MediaSize,
			&message.MediaHash,
			&message.Timestamp,
			&message.Direction,
			&message.Status,
		); err != nil {
			return err
		}

		message, err = s.decryptMessage(message)
		if err != nil {
			return err
		}
		if !isStoredImageRef(message.Text) {
			continue
		}

		currentAbsolutePath, err := s.mediaAbsolutePathFromRefUnlocked(message.Text)
		if err != nil {
			continue
		}
		raw, resolvedMediaType, err := s.resolveImagePayloadUnlocked(message.Text, message.MediaType)
		if err != nil {
			continue
		}
		if message.MediaType == "" {
			message.MediaType = resolvedMediaType
		}
		if message.MediaName == "" {
			message.MediaName = mediautil.NormalizeImageFileName("", message.MediaType)
		}
		if message.MediaSize <= 0 {
			message.MediaSize = int64(len(raw))
		}

		nextRef, nextHash, err := s.storeImagePayloadUnlocked(raw, message.MediaType, message.MediaName)
		if err != nil {
			return err
		}
		nextAbsolutePath := filepath.Join(s.dataDir, filepath.Clean(filepath.FromSlash(mediaPathFromRef(nextRef))))
		if currentAbsolutePath != nextAbsolutePath {
			_ = s.removeMediaFileUnlocked(currentAbsolutePath)
		}

		message.Text = nextRef
		message.MediaHash = nextHash
		pending = append(pending, message)
	}
	if err := rows.Err(); err != nil {
		return err
	}

	for _, message := range pending {
		encrypted, err := s.encryptMessage(message)
		if err != nil {
			return err
		}
		if _, err := s.db.Exec(`
			UPDATE messages
			SET sender_name = ?, text = ?, media_name = ?, media_type = ?, media_size = ?, media_hash = ?
			WHERE id = ?
		`, encrypted.SenderName, encrypted.Text, encrypted.MediaName, encrypted.MediaType, encrypted.MediaSize, encrypted.MediaHash, encrypted.ID); err != nil {
			return err
		}
	}

	s.mediaUsageKnown = false
	s.mediaUsageBytes = 0
	return nil
}

func (s *Store) migrateStoredFileMessages() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	rows, err := s.db.Query(`
		SELECT id, peer_id, sender_id, sender_name, kind, text, media_name, media_type, media_size, media_hash, timestamp, direction, status
		FROM messages
		WHERE kind = 'file'
	`)
	if err != nil {
		return err
	}
	defer rows.Close()

	var pending []ChatMessage
	for rows.Next() {
		var message ChatMessage
		if err := rows.Scan(
			&message.ID,
			&message.PeerID,
			&message.SenderID,
			&message.SenderName,
			&message.Kind,
			&message.Text,
			&message.MediaName,
			&message.MediaType,
			&message.MediaSize,
			&message.MediaHash,
			&message.Timestamp,
			&message.Direction,
			&message.Status,
		); err != nil {
			return err
		}

		message, err = s.decryptMessage(message)
		if err != nil {
			return err
		}
		if !isStoredFileRef(message.Text) {
			continue
		}

		currentAbsolutePath, err := s.fileAbsolutePathFromRefUnlocked(message.Text)
		if err != nil {
			continue
		}
		raw, resolvedMediaType, err := s.resolveFilePayloadUnlocked(message.Text, message.MediaType)
		if err != nil {
			continue
		}
		if message.MediaType == "" {
			message.MediaType = resolvedMediaType
		}
		if message.MediaName == "" {
			message.MediaName = mediautil.NormalizeFileName("", message.MediaType)
		}
		if message.MediaSize <= 0 {
			message.MediaSize = int64(len(raw))
		}

		nextRef, nextHash, err := s.storeFilePayloadUnlocked(raw, message.MediaType, message.MediaName)
		if err != nil {
			return err
		}
		nextAbsolutePath := filepath.Join(s.dataDir, filepath.Clean(filepath.FromSlash(filePathFromRef(nextRef))))
		if currentAbsolutePath != nextAbsolutePath {
			_ = s.removeFileLinkUnlocked(currentAbsolutePath)
		}
		_, legacyBlobAbsolutePath := s.fileBlobPathUnlocked(nextHash)
		if legacyBlobAbsolutePath != nextAbsolutePath {
			_ = s.removeFileLinkUnlocked(legacyBlobAbsolutePath)
		}

		message.Text = nextRef
		message.MediaHash = nextHash
		pending = append(pending, message)
	}
	if err := rows.Err(); err != nil {
		return err
	}

	for _, message := range pending {
		encrypted, err := s.encryptMessage(message)
		if err != nil {
			return err
		}
		if _, err := s.db.Exec(`
			UPDATE messages
			SET sender_name = ?, text = ?, media_name = ?, media_type = ?, media_size = ?, media_hash = ?
			WHERE id = ?
		`, encrypted.SenderName, encrypted.Text, encrypted.MediaName, encrypted.MediaType, encrypted.MediaSize, encrypted.MediaHash, encrypted.ID); err != nil {
			return err
		}
	}

	s.mediaUsageKnown = false
	s.mediaUsageBytes = 0
	return nil
}

func copyDir(src, dst string) error {
	if _, err := os.Stat(src); errors.Is(err, os.ErrNotExist) {
		return nil
	} else if err != nil {
		return err
	}

	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relativePath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		targetPath := filepath.Join(dst, relativePath)

		if info.IsDir() {
			return os.MkdirAll(targetPath, info.Mode())
		}

		if err := os.MkdirAll(filepath.Dir(targetPath), 0o755); err != nil {
			return err
		}
		if _, err := copyFile(path, targetPath); err != nil {
			return err
		}
		return os.Chmod(targetPath, info.Mode())
	})
}
