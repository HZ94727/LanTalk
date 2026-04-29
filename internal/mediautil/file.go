package mediautil

import (
	"encoding/base64"
	"errors"
	"fmt"
	"mime"
	"net/http"
	"path/filepath"
	"strings"
)

func DecodeDataURL(dataURL string) ([]byte, string, error) {
	dataURL = strings.TrimSpace(dataURL)
	if !strings.HasPrefix(dataURL, "data:") {
		return nil, "", errors.New("only data URLs are supported")
	}

	parts := strings.SplitN(dataURL, ",", 2)
	if len(parts) != 2 {
		return nil, "", errors.New("invalid payload")
	}

	header := parts[0]
	if !strings.HasSuffix(header, ";base64") {
		return nil, "", errors.New("payload must be base64 encoded")
	}

	mediaType := normalizeMediaType(strings.TrimSuffix(strings.TrimPrefix(header, "data:"), ";base64"))
	raw, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, "", fmt.Errorf("decode payload: %w", err)
	}

	return raw, mediaType, nil
}

func normalizeMediaType(mediaType string) string {
	mediaType = strings.TrimSpace(strings.Split(mediaType, ";")[0])
	if mediaType == "" {
		return "application/octet-stream"
	}
	return strings.ToLower(mediaType)
}

func MediaTypeForFileName(name string) string {
	ext := strings.ToLower(filepath.Ext(strings.TrimSpace(name)))
	if ext == "" {
		return "application/octet-stream"
	}

	mediaType := normalizeMediaType(mime.TypeByExtension(ext))
	if mediaType == "" {
		return "application/octet-stream"
	}
	return mediaType
}

func DetectMediaType(name string, raw []byte) string {
	mediaType := MediaTypeForFileName(name)
	if mediaType != "application/octet-stream" {
		return mediaType
	}

	if len(raw) == 0 {
		return mediaType
	}

	sniffSize := len(raw)
	if sniffSize > 512 {
		sniffSize = 512
	}
	detected := normalizeMediaType(http.DetectContentType(raw[:sniffSize]))
	if detected == "" {
		return mediaType
	}
	return detected
}

func EncodeDataURL(raw []byte, mediaType string) string {
	return "data:" + normalizeMediaType(mediaType) + ";base64," + base64.StdEncoding.EncodeToString(raw)
}

func ExtensionForFileMediaType(mediaType string) string {
	switch normalizeMediaType(mediaType) {
	case "application/json":
		return ".json"
	case "application/pdf":
		return ".pdf"
	case "application/zip":
		return ".zip"
	case "application/x-7z-compressed":
		return ".7z"
	case "application/x-rar-compressed":
		return ".rar"
	case "application/vnd.openxmlformats-officedocument.wordprocessingml.document":
		return ".docx"
	case "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet":
		return ".xlsx"
	case "application/vnd.openxmlformats-officedocument.presentationml.presentation":
		return ".pptx"
	case "application/msword":
		return ".doc"
	case "application/vnd.ms-excel":
		return ".xls"
	case "application/vnd.ms-powerpoint":
		return ".ppt"
	case "text/plain":
		return ".txt"
	case "text/csv":
		return ".csv"
	}

	extensions, err := mime.ExtensionsByType(normalizeMediaType(mediaType))
	if err == nil {
		for _, ext := range extensions {
			if ext == "" {
				continue
			}
			return ext
		}
	}

	return ".bin"
}

func NormalizeFileName(name, mediaType string) string {
	name = strings.TrimSpace(filepath.Base(name))
	ext := ExtensionForFileMediaType(mediaType)

	if name == "" || name == "." || name == string(filepath.Separator) {
		return "file" + ext
	}

	if filepath.Ext(name) == "" && ext != "" {
		return name + ext
	}

	return name
}
