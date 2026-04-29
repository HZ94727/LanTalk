package chat

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"os"
)

const encryptedPrefix = "enc:v1:"

type cryptoManager struct {
	key []byte
}

func loadCryptoManager(keyPath string) (*cryptoManager, error) {
	protectedKey, err := os.ReadFile(keyPath)
	if errors.Is(err, os.ErrNotExist) {
		rawKey := make([]byte, 32)
		if _, err := io.ReadFull(rand.Reader, rawKey); err != nil {
			return nil, err
		}

		protectedKey, err = protectKey(rawKey)
		if err != nil {
			return nil, err
		}
		if err := os.WriteFile(keyPath, protectedKey, 0o600); err != nil {
			return nil, err
		}
		return &cryptoManager{key: rawKey}, nil
	}
	if err != nil {
		return nil, err
	}

	rawKey, err := unprotectKey(protectedKey)
	if err != nil {
		return nil, err
	}
	if len(rawKey) != 32 {
		return nil, fmt.Errorf("unexpected key length: %d", len(rawKey))
	}

	return &cryptoManager{key: rawKey}, nil
}

func (m *cryptoManager) EncryptString(value string) (string, error) {
	if value == "" {
		return value, nil
	}
	if isEncrypted(value) {
		return value, nil
	}

	block, err := aes.NewCipher(m.key)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(value), nil)
	return encryptedPrefix + base64.StdEncoding.EncodeToString(ciphertext), nil
}

func (m *cryptoManager) DecryptString(value string) (string, error) {
	if value == "" || !isEncrypted(value) {
		return value, nil
	}

	encoded := value[len(encryptedPrefix):]
	payload, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(m.key)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	if len(payload) < gcm.NonceSize() {
		return "", errors.New("encrypted payload too short")
	}

	nonce := payload[:gcm.NonceSize()]
	ciphertext := payload[gcm.NonceSize():]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}
	return string(plaintext), nil
}

func isEncrypted(value string) bool {
	return len(value) > len(encryptedPrefix) && value[:len(encryptedPrefix)] == encryptedPrefix
}
