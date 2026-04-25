package chat

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sync"
)

type persistedState struct {
	Profile       Profile                  `json:"profile"`
	Settings      Settings                 `json:"settings"`
	Conversations map[string][]ChatMessage `json:"conversations"`
}

type Store struct {
	path string
	mu   sync.Mutex
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

	return &Store{path: filepath.Join(dataDir, "state.json")}, nil
}

func (s *Store) Load() (persistedState, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var state persistedState

	data, err := os.ReadFile(s.path)
	if errors.Is(err, os.ErrNotExist) {
		state.Conversations = make(map[string][]ChatMessage)
		return state, nil
	}
	if err != nil {
		return state, err
	}

	if err := json.Unmarshal(data, &state); err != nil {
		return state, err
	}
	if state.Conversations == nil {
		state.Conversations = make(map[string][]ChatMessage)
	}

	return state, nil
}

func (s *Store) Save(state persistedState) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(s.path, data, 0o644)
}

func (s *Store) Path() string {
	return s.path
}

func (s *Store) Move(state persistedState, targetDir string) (string, error) {
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

	targetPath := filepath.Join(targetDir, "state.json")
	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return "", err
	}
	if err := os.WriteFile(targetPath, data, 0o644); err != nil {
		return "", err
	}

	oldPath := s.path
	s.path = targetPath
	if oldPath != targetPath {
		if _, err := os.Stat(oldPath); err == nil {
			_ = os.Remove(oldPath)
		}
	}

	return targetPath, nil
}
