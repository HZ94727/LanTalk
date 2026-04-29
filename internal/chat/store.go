package chat

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	_ "modernc.org/sqlite"
)

type persistedState struct {
	Profile       Profile                  `json:"profile"`
	Settings      Settings                 `json:"settings"`
	Peers         []Peer                   `json:"peers"`
	Conversations map[string][]ChatMessage `json:"conversations"`
}

type jsonState struct {
	Profile       Profile                  `json:"profile"`
	Settings      Settings                 `json:"settings"`
	Conversations map[string][]ChatMessage `json:"conversations"`
}

type Store struct {
	dataDir         string
	mediaDir        string
	path            string
	keyPath         string
	db              *sql.DB
	crypto          *cryptoManager
	mediaUsageBytes int64
	mediaUsageKnown bool
	mu              sync.Mutex
}

func NewStore(appName string) (*Store, error) {
	baseDir, err := os.UserConfigDir()
	if err != nil {
		return nil, err
	}

	return NewStoreFromDirectory(filepath.Join(baseDir, appName))
}

func NewStoreFromDirectory(dataDir string) (*Store, error) {
	if err := os.MkdirAll(dataDir, 0o755); err != nil {
		return nil, err
	}

	dbPath := filepath.Join(dataDir, "lantalk.db")
	keyPath := filepath.Join(dataDir, "master.key")
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}

	crypto, err := loadCryptoManager(keyPath)
	if err != nil {
		_ = db.Close()
		return nil, err
	}

	store := &Store{
		dataDir:  dataDir,
		mediaDir: filepath.Join(dataDir, "media"),
		path:     dbPath,
		keyPath:  keyPath,
		db:       db,
		crypto:   crypto,
	}

	if err := os.MkdirAll(store.mediaDir, 0o755); err != nil {
		_ = db.Close()
		return nil, err
	}

	if err := store.init(); err != nil {
		_ = db.Close()
		return nil, err
	}
	if err := store.migrateJSONIfNeeded(filepath.Join(dataDir, "state.json")); err != nil {
		_ = db.Close()
		return nil, err
	}
	if err := store.migratePlaintextProfile(); err != nil {
		_ = db.Close()
		return nil, err
	}
	if err := store.migratePlaintextPeers(); err != nil {
		_ = db.Close()
		return nil, err
	}
	if err := store.migratePlaintextMessages(); err != nil {
		_ = db.Close()
		return nil, err
	}
	if err := store.migrateInlineImageMessages(); err != nil {
		_ = db.Close()
		return nil, err
	}
	if err := store.migrateStoredImageMessages(); err != nil {
		_ = db.Close()
		return nil, err
	}
	if err := store.migrateStoredFileMessages(); err != nil {
		_ = db.Close()
		return nil, err
	}

	return store, nil
}

func (s *Store) init() error {
	statements := []string{
		`CREATE TABLE IF NOT EXISTS profile (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			listen_port INTEGER NOT NULL
		);`,
		`CREATE TABLE IF NOT EXISTS settings (
			id INTEGER PRIMARY KEY CHECK (id = 1),
			language TEXT NOT NULL,
			theme TEXT NOT NULL
		);`,
		`CREATE TABLE IF NOT EXISTS peers (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			address TEXT NOT NULL,
			listen_port INTEGER NOT NULL,
			last_seen INTEGER NOT NULL,
			source TEXT NOT NULL
		);`,
		`CREATE TABLE IF NOT EXISTS messages (
			id TEXT PRIMARY KEY,
			peer_id TEXT NOT NULL,
			sender_id TEXT NOT NULL,
			sender_name TEXT NOT NULL,
			kind TEXT NOT NULL DEFAULT 'text',
			text TEXT NOT NULL,
			media_name TEXT NOT NULL DEFAULT '',
			media_type TEXT NOT NULL DEFAULT '',
			media_size INTEGER NOT NULL DEFAULT 0,
			media_hash TEXT NOT NULL DEFAULT '',
			timestamp INTEGER NOT NULL,
			direction TEXT NOT NULL,
			status TEXT NOT NULL
		);`,
		`CREATE INDEX IF NOT EXISTS idx_messages_peer_time ON messages(peer_id, timestamp);`,
	}

	for _, statement := range statements {
		if _, err := s.db.Exec(statement); err != nil {
			return err
		}
	}

	if err := s.ensureMessageColumns(); err != nil {
		return err
	}
	if _, err := s.db.Exec(`CREATE INDEX IF NOT EXISTS idx_messages_kind_hash ON messages(kind, media_hash);`); err != nil {
		return err
	}

	return nil
}

func (s *Store) Load() (persistedState, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	state := persistedState{
		Conversations: make(map[string][]ChatMessage),
	}

	if err := s.loadProfile(&state); err != nil {
		return state, err
	}
	if err := s.loadSettings(&state); err != nil {
		return state, err
	}
	if err := s.loadPeers(&state); err != nil {
		return state, err
	}

	return state, nil
}

func (s *Store) SaveProfile(profile Profile) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	encrypted, err := s.encryptProfile(profile)
	if err != nil {
		return err
	}

	_, err = s.db.Exec(`
		INSERT INTO profile(id, name, listen_port)
		VALUES (?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET
			name = excluded.name,
			listen_port = excluded.listen_port
	`, encrypted.ID, encrypted.Name, encrypted.ListenPort)
	return err
}

func (s *Store) SaveSettings(settings Settings) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, err := s.db.Exec(`
		INSERT INTO settings(id, language, theme)
		VALUES (1, ?, ?)
		ON CONFLICT(id) DO UPDATE SET
			language = excluded.language,
			theme = excluded.theme
	`, settings.Language, settings.Theme)
	return err
}

func (s *Store) SavePeer(peer Peer) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	encrypted, err := s.encryptPeer(peer)
	if err != nil {
		return err
	}

	_, err = s.db.Exec(`
		INSERT INTO peers(id, name, address, listen_port, last_seen, source)
		VALUES (?, ?, ?, ?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET
			name = excluded.name,
			address = excluded.address,
			listen_port = excluded.listen_port,
			last_seen = excluded.last_seen,
			source = excluded.source
	`, encrypted.ID, encrypted.Name, encrypted.Address, encrypted.ListenPort, encrypted.LastSeen, encrypted.Source)
	return err
}

func (s *Store) DeletePeer(peerID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, err := s.db.Exec(`DELETE FROM peers WHERE id = ?`, peerID)
	return err
}

func (s *Store) HasMessagesForPeer(peerID string) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var count int
	if err := s.db.QueryRow(`SELECT COUNT(1) FROM messages WHERE peer_id = ?`, peerID).Scan(&count); err != nil {
		return false, err
	}
	return count > 0, nil
}

func (s *Store) LoadConversationPeerIDs() ([]string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	rows, err := s.db.Query(`
		SELECT DISTINCT peer_id
		FROM messages
		ORDER BY peer_id
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	peerIDs := make([]string, 0, 16)
	for rows.Next() {
		var peerID string
		if err := rows.Scan(&peerID); err != nil {
			return nil, err
		}
		peerID = strings.TrimSpace(peerID)
		if peerID == "" {
			continue
		}
		peerIDs = append(peerIDs, peerID)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return peerIDs, nil
}

func (s *Store) SaveMessage(message ChatMessage) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	normalized, err := s.normalizeMessageMediaUnlocked(message)
	if err != nil {
		return err
	}

	encrypted, err := s.encryptMessage(normalized)
	if err != nil {
		return err
	}

	_, err = s.db.Exec(`
		INSERT INTO messages(id, peer_id, sender_id, sender_name, kind, text, media_name, media_type, media_size, media_hash, timestamp, direction, status)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET
			peer_id = excluded.peer_id,
			sender_id = excluded.sender_id,
			sender_name = excluded.sender_name,
			kind = excluded.kind,
			text = excluded.text,
			media_name = excluded.media_name,
			media_type = excluded.media_type,
			media_size = excluded.media_size,
			media_hash = excluded.media_hash,
			timestamp = excluded.timestamp,
			direction = excluded.direction,
			status = excluded.status
	`, encrypted.ID, encrypted.PeerID, encrypted.SenderID, encrypted.SenderName, encrypted.Kind, encrypted.Text, encrypted.MediaName, encrypted.MediaType, encrypted.MediaSize, encrypted.MediaHash, encrypted.Timestamp, encrypted.Direction, encrypted.Status)
	return err
}

func (s *Store) MarkPendingOutgoingFailed() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, err := s.db.Exec(`
		UPDATE messages
		SET status = 'failed'
		WHERE direction = 'outbound' AND status = 'sending'
	`)
	return err
}

func (s *Store) LoadMessage(peerID, messageID string) (ChatMessage, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var message ChatMessage
	err := s.db.QueryRow(`
		SELECT id, peer_id, sender_id, sender_name, kind, text, media_name, media_type, media_size, media_hash, timestamp, direction, status
		FROM messages
		WHERE peer_id = ? AND id = ?
		LIMIT 1
	`, peerID, messageID).Scan(
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
	)
	if err != nil {
		return ChatMessage{}, err
	}

	return s.decryptMessage(message)
}

func (s *Store) DeleteMessage(messageID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, err := s.db.Exec(`DELETE FROM messages WHERE id = ?`, messageID)
	return err
}

func (s *Store) Path() string {
	return s.path
}

func (s *Store) DataDir() string {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.dataDir
}

func (s *Store) Move(targetDir string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if targetDir == "" {
		return "", errors.New("target directory is empty")
	}

	targetDir, err := filepath.Abs(targetDir)
	if err != nil {
		return "", err
	}
	if err := os.MkdirAll(targetDir, 0o755); err != nil {
		return "", err
	}

	targetPath := filepath.Join(targetDir, "lantalk.db")
	targetKeyPath := filepath.Join(targetDir, "master.key")
	targetMediaDir := filepath.Join(targetDir, "media")
	if targetPath == s.path {
		return targetPath, nil
	}

	if _, err := s.db.Exec(`VACUUM INTO ?`, targetPath); err != nil {
		return "", err
	}

	oldPath := s.path
	oldKeyPath := s.keyPath
	_ = s.db.Close()

	if _, err := copyFile(oldKeyPath, targetKeyPath); err != nil {
		return "", err
	}
	if err := copyDir(filepath.Join(s.dataDir, "media"), targetMediaDir); err != nil {
		return "", err
	}

	newDB, err := sql.Open("sqlite", targetPath)
	if err != nil {
		return "", err
	}
	s.db = newDB
	s.dataDir = targetDir
	s.mediaDir = targetMediaDir
	s.path = targetPath
	s.keyPath = targetKeyPath
	s.mediaUsageKnown = false
	s.mediaUsageBytes = 0

	if err := s.init(); err != nil {
		return "", err
	}

	if oldPath != targetPath {
		_ = os.Remove(oldPath)
		_ = os.Remove(oldKeyPath)
		_ = os.Remove(filepath.Join(filepath.Dir(oldPath), "state.json"))
		_ = os.RemoveAll(filepath.Join(filepath.Dir(oldPath), "media"))
	}

	return targetPath, nil
}

func (s *Store) loadProfile(state *persistedState) error {
	row := s.db.QueryRow(`SELECT id, name, listen_port FROM profile LIMIT 1`)
	if err := row.Scan(&state.Profile.ID, &state.Profile.Name, &state.Profile.ListenPort); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}
		return err
	}
	decrypted, err := s.decryptProfile(state.Profile)
	if err != nil {
		return err
	}
	state.Profile = decrypted
	return nil
}

func (s *Store) loadSettings(state *persistedState) error {
	row := s.db.QueryRow(`SELECT language, theme FROM settings WHERE id = 1`)
	if err := row.Scan(&state.Settings.Language, &state.Settings.Theme); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}
		return err
	}
	return nil
}

func (s *Store) loadPeers(state *persistedState) error {
	rows, err := s.db.Query(`SELECT id, name, address, listen_port, last_seen, source FROM peers ORDER BY id`)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var peer Peer
		if err := rows.Scan(&peer.ID, &peer.Name, &peer.Address, &peer.ListenPort, &peer.LastSeen, &peer.Source); err != nil {
			return err
		}
		peer, err = s.decryptPeer(peer)
		if err != nil {
			return err
		}
		state.Peers = append(state.Peers, peer)
	}

	return rows.Err()
}

func (s *Store) loadMessages(state *persistedState) error {
	rows, err := s.db.Query(`
		SELECT id, peer_id, sender_id, sender_name, kind, text, media_name, media_type, media_size, media_hash, timestamp, direction, status
		FROM messages
		ORDER BY timestamp, id
	`)
	if err != nil {
		return err
	}
	defer rows.Close()

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
		state.Conversations[message.PeerID] = append(state.Conversations[message.PeerID], message)
	}

	return rows.Err()
}

func (s *Store) LoadConversationPage(peerID string, beforeTimestamp int64, beforeID string, limit int) (ConversationPage, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	page := ConversationPage{
		PeerID:   peerID,
		Messages: []ChatMessage{},
		HasMore:  false,
	}

	if strings.TrimSpace(peerID) == "" {
		return page, errors.New("peer id is required")
	}
	if limit <= 0 {
		limit = 40
	}

	var (
		rows *sql.Rows
		err  error
	)

	if beforeTimestamp > 0 {
		rows, err = s.db.Query(`
			SELECT id, peer_id, sender_id, sender_name, kind, text, media_name, media_type, media_size, media_hash, timestamp, direction, status
			FROM messages
			WHERE peer_id = ? AND (timestamp < ? OR (timestamp = ? AND id < ?))
			ORDER BY timestamp DESC, id DESC
			LIMIT ?
		`, peerID, beforeTimestamp, beforeTimestamp, beforeID, limit+1)
	} else {
		rows, err = s.db.Query(`
			SELECT id, peer_id, sender_id, sender_name, kind, text, media_name, media_type, media_size, media_hash, timestamp, direction, status
			FROM messages
			WHERE peer_id = ?
			ORDER BY timestamp DESC, id DESC
			LIMIT ?
		`, peerID, limit+1)
	}
	if err != nil {
		return page, err
	}
	defer rows.Close()

	descMessages := make([]ChatMessage, 0, limit+1)
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
			return page, err
		}
		message, err = s.decryptMessage(message)
		if err != nil {
			return page, err
		}
		descMessages = append(descMessages, message)
	}
	if err := rows.Err(); err != nil {
		return page, err
	}

	if len(descMessages) > limit {
		page.HasMore = true
		descMessages = descMessages[:limit]
	}

	page.Messages = make([]ChatMessage, 0, len(descMessages))
	for i := len(descMessages) - 1; i >= 0; i-- {
		page.Messages = append(page.Messages, descMessages[i])
	}

	return page, nil
}

func (s *Store) migrateJSONIfNeeded(jsonPath string) error {
	if _, err := os.Stat(jsonPath); errors.Is(err, os.ErrNotExist) {
		return nil
	} else if err != nil {
		return err
	}

	var count int
	if err := s.db.QueryRow(`SELECT COUNT(*) FROM messages`).Scan(&count); err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	data, err := os.ReadFile(jsonPath)
	if err != nil {
		return err
	}

	var legacy jsonState
	if err := json.Unmarshal(data, &legacy); err != nil {
		return err
	}

	if legacy.Conversations == nil {
		legacy.Conversations = make(map[string][]ChatMessage)
	}

	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	if legacy.Profile.ID != "" {
		if _, err = tx.Exec(`
			INSERT INTO profile(id, name, listen_port)
			VALUES (?, ?, ?)
			ON CONFLICT(id) DO UPDATE SET
				name = excluded.name,
				listen_port = excluded.listen_port
		`, legacy.Profile.ID, legacy.Profile.Name, legacy.Profile.ListenPort); err != nil {
			return err
		}
	}

	settings := normalizeSettings(legacy.Settings)
	if _, err = tx.Exec(`
		INSERT INTO settings(id, language, theme)
		VALUES (1, ?, ?)
		ON CONFLICT(id) DO UPDATE SET
			language = excluded.language,
			theme = excluded.theme
	`, settings.Language, settings.Theme); err != nil {
		return err
	}

	peerMap := make(map[string]Peer)
	for peerID, messages := range legacy.Conversations {
		for _, message := range messages {
			encrypted, encryptErr := s.encryptMessage(message)
			if encryptErr != nil {
				err = encryptErr
				return err
			}
			if _, err = tx.Exec(`
				INSERT INTO messages(id, peer_id, sender_id, sender_name, kind, text, media_name, media_type, media_size, media_hash, timestamp, direction, status)
				VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
				ON CONFLICT(id) DO UPDATE SET
					peer_id = excluded.peer_id,
					sender_id = excluded.sender_id,
					sender_name = excluded.sender_name,
					kind = excluded.kind,
					text = excluded.text,
					media_name = excluded.media_name,
					media_type = excluded.media_type,
					media_size = excluded.media_size,
					media_hash = excluded.media_hash,
					timestamp = excluded.timestamp,
					direction = excluded.direction,
					status = excluded.status
			`, encrypted.ID, encrypted.PeerID, encrypted.SenderID, encrypted.SenderName, encrypted.Kind, encrypted.Text, encrypted.MediaName, encrypted.MediaType, encrypted.MediaSize, encrypted.MediaHash, encrypted.Timestamp, encrypted.Direction, encrypted.Status); err != nil {
				return err
			}

			peer := peerMap[peerID]
			peer.ID = peerID
			peer.Name = fallbackPeerName(peer.Name, message)
			peer.Source = inferPeerSource(peerID)
			peerMap[peerID] = peer
		}
	}

	peers := make([]Peer, 0, len(peerMap))
	for _, peer := range peerMap {
		peers = append(peers, peer)
	}
	sort.Slice(peers, func(i, j int) bool {
		if peers[i].Name == peers[j].Name {
			return peers[i].ID < peers[j].ID
		}
		return peers[i].Name < peers[j].Name
	})

	for _, peer := range peers {
		encryptedPeer, encryptErr := s.encryptPeer(peer)
		if encryptErr != nil {
			err = encryptErr
			return err
		}
		if _, err = tx.Exec(`
			INSERT INTO peers(id, name, address, listen_port, last_seen, source)
			VALUES (?, ?, ?, ?, ?, ?)
			ON CONFLICT(id) DO UPDATE SET
				name = excluded.name,
				address = excluded.address,
				listen_port = excluded.listen_port,
				last_seen = excluded.last_seen,
				source = excluded.source
		`, encryptedPeer.ID, encryptedPeer.Name, encryptedPeer.Address, encryptedPeer.ListenPort, encryptedPeer.LastSeen, encryptedPeer.Source); err != nil {
			return err
		}
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func fallbackPeerName(current string, message ChatMessage) string {
	if current != "" {
		return current
	}
	if message.Direction == "inbound" && message.SenderName != "" {
		return message.SenderName
	}
	if message.Direction == "outbound" && strings.TrimSpace(message.PeerID) == debugEchoPeerID {
		return "Echo Bot"
	}
	return fmt.Sprintf("Peer %s", shortID(message.PeerID))
}

func inferPeerSource(peerID string) string {
	if peerID == debugEchoPeerID {
		return "debug"
	}
	return "manual"
}

func shortID(value string) string {
	if len(value) <= 8 {
		return value
	}
	return value[:8]
}

func (s *Store) migratePlaintextMessages() error {
	rows, err := s.db.Query(`
		SELECT id, peer_id, sender_id, sender_name, kind, text, media_name, media_type, media_size, media_hash, timestamp, direction, status
		FROM messages
	`)
	if err != nil {
		return err
	}
	defer rows.Close()

	var plaintext []ChatMessage
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
		if !isEncrypted(message.Text) || !isEncrypted(message.SenderName) {
			plaintext = append(plaintext, message)
		}
	}
	if err := rows.Err(); err != nil {
		return err
	}
	if len(plaintext) == 0 {
		return nil
	}

	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	for _, message := range plaintext {
		encrypted, err := s.encryptMessage(message)
		if err != nil {
			_ = tx.Rollback()
			return err
		}
		if _, err := tx.Exec(`
			UPDATE messages
			SET sender_name = ?, text = ?, media_name = ?, media_type = ?, media_size = ?, media_hash = ?
			WHERE id = ?
		`, encrypted.SenderName, encrypted.Text, encrypted.MediaName, encrypted.MediaType, encrypted.MediaSize, encrypted.MediaHash, encrypted.ID); err != nil {
			_ = tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}

func (s *Store) migratePlaintextProfile() error {
	row := s.db.QueryRow(`SELECT id, name, listen_port FROM profile LIMIT 1`)

	var profile Profile
	if err := row.Scan(&profile.ID, &profile.Name, &profile.ListenPort); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}
		return err
	}
	if isEncrypted(profile.Name) {
		return nil
	}

	encrypted, err := s.encryptProfile(profile)
	if err != nil {
		return err
	}

	_, err = s.db.Exec(`
		UPDATE profile
		SET name = ?
		WHERE id = ?
	`, encrypted.Name, encrypted.ID)
	return err
}

func (s *Store) migratePlaintextPeers() error {
	rows, err := s.db.Query(`
		SELECT id, name, address, listen_port, last_seen, source
		FROM peers
	`)
	if err != nil {
		return err
	}
	defer rows.Close()

	var plaintext []Peer
	for rows.Next() {
		var peer Peer
		if err := rows.Scan(&peer.ID, &peer.Name, &peer.Address, &peer.ListenPort, &peer.LastSeen, &peer.Source); err != nil {
			return err
		}
		if !isEncrypted(peer.Name) {
			plaintext = append(plaintext, peer)
		}
	}
	if err := rows.Err(); err != nil {
		return err
	}
	if len(plaintext) == 0 {
		return nil
	}

	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	for _, peer := range plaintext {
		encrypted, err := s.encryptPeer(peer)
		if err != nil {
			_ = tx.Rollback()
			return err
		}
		if _, err := tx.Exec(`
			UPDATE peers
			SET name = ?
			WHERE id = ?
		`, encrypted.Name, encrypted.ID); err != nil {
			_ = tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}

func (s *Store) encryptProfile(profile Profile) (Profile, error) {
	encrypted := profile
	var err error
	encrypted.Name, err = s.crypto.EncryptString(profile.Name)
	if err != nil {
		return Profile{}, err
	}
	return encrypted, nil
}

func (s *Store) decryptProfile(profile Profile) (Profile, error) {
	decrypted := profile
	var err error
	decrypted.Name, err = s.crypto.DecryptString(profile.Name)
	if err != nil {
		return Profile{}, err
	}
	return decrypted, nil
}

func (s *Store) encryptPeer(peer Peer) (Peer, error) {
	encrypted := peer
	var err error
	encrypted.Name, err = s.crypto.EncryptString(peer.Name)
	if err != nil {
		return Peer{}, err
	}
	return encrypted, nil
}

func (s *Store) decryptPeer(peer Peer) (Peer, error) {
	decrypted := peer
	var err error
	decrypted.Name, err = s.crypto.DecryptString(peer.Name)
	if err != nil {
		return Peer{}, err
	}
	return decrypted, nil
}

func (s *Store) encryptMessage(message ChatMessage) (ChatMessage, error) {
	encrypted := message
	var err error
	if encrypted.Kind == "" {
		encrypted.Kind = "text"
	}
	encrypted.SenderName, err = s.crypto.EncryptString(message.SenderName)
	if err != nil {
		return ChatMessage{}, err
	}
	encrypted.Text, err = s.crypto.EncryptString(message.Text)
	if err != nil {
		return ChatMessage{}, err
	}
	encrypted.MediaName, err = s.crypto.EncryptString(message.MediaName)
	if err != nil {
		return ChatMessage{}, err
	}
	encrypted.MediaType, err = s.crypto.EncryptString(message.MediaType)
	if err != nil {
		return ChatMessage{}, err
	}
	encrypted.MediaSize = message.MediaSize
	encrypted.MediaHash = message.MediaHash
	return encrypted, nil
}

func (s *Store) decryptMessage(message ChatMessage) (ChatMessage, error) {
	decrypted := message
	var err error
	if decrypted.Kind == "" {
		decrypted.Kind = "text"
	}
	decrypted.SenderName, err = s.crypto.DecryptString(message.SenderName)
	if err != nil {
		return ChatMessage{}, err
	}
	decrypted.Text, err = s.crypto.DecryptString(message.Text)
	if err != nil {
		return ChatMessage{}, err
	}
	decrypted.MediaName, err = s.crypto.DecryptString(message.MediaName)
	if err != nil {
		return ChatMessage{}, err
	}
	decrypted.MediaType, err = s.crypto.DecryptString(message.MediaType)
	if err != nil {
		return ChatMessage{}, err
	}
	decrypted.MediaSize = message.MediaSize
	decrypted.MediaHash = message.MediaHash
	return decrypted, nil
}

func (s *Store) ensureMessageColumns() error {
	statements := []string{
		`ALTER TABLE messages ADD COLUMN kind TEXT NOT NULL DEFAULT 'text';`,
		`ALTER TABLE messages ADD COLUMN media_name TEXT NOT NULL DEFAULT '';`,
		`ALTER TABLE messages ADD COLUMN media_type TEXT NOT NULL DEFAULT '';`,
		`ALTER TABLE messages ADD COLUMN media_size INTEGER NOT NULL DEFAULT 0;`,
		`ALTER TABLE messages ADD COLUMN media_hash TEXT NOT NULL DEFAULT '';`,
	}

	for _, statement := range statements {
		if _, err := s.db.Exec(statement); err != nil && !strings.Contains(strings.ToLower(err.Error()), "duplicate column name") {
			return err
		}
	}

	return nil
}

func copyFile(src, dst string) (int64, error) {
	input, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer input.Close()

	output, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer output.Close()

	n, err := output.ReadFrom(input)
	if err != nil {
		return n, err
	}

	return n, output.Sync()
}
