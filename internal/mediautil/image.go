package mediautil

import (
	"encoding/base64"
	"errors"
	"fmt"
	"path/filepath"
	"strings"
)

func DecodeImageDataURL(dataURL string) ([]byte, string, error) {
	dataURL = strings.TrimSpace(dataURL)
	if !strings.HasPrefix(dataURL, "data:image/") {
		return nil, "", errors.New("only image data URLs are supported")
	}

	parts := strings.SplitN(dataURL, ",", 2)
	if len(parts) != 2 {
		return nil, "", errors.New("invalid image payload")
	}

	header := parts[0]
	if !strings.HasSuffix(header, ";base64") {
		return nil, "", errors.New("image payload must be base64 encoded")
	}

	mediaType := strings.TrimSuffix(strings.TrimPrefix(header, "data:"), ";base64")
	if mediaType == "" {
		return nil, "", errors.New("missing image media type")
	}

	raw, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, "", fmt.Errorf("decode image payload: %w", err)
	}

	return raw, mediaType, nil
}

func ExtensionForMediaType(mediaType string) string {
	switch strings.ToLower(strings.TrimSpace(mediaType)) {
	case "image/jpeg", "image/jpg":
		return ".jpg"
	case "image/png":
		return ".png"
	case "image/webp":
		return ".webp"
	case "image/gif":
		return ".gif"
	case "image/bmp":
		return ".bmp"
	case "image/x-icon", "image/vnd.microsoft.icon":
		return ".ico"
	default:
		return ".png"
	}
}

func NormalizeImageFileName(name, mediaType string) string {
	name = strings.TrimSpace(filepath.Base(name))
	ext := ExtensionForMediaType(mediaType)

	if name == "" || name == "." || name == string(filepath.Separator) {
		return "image" + ext
	}

	if filepath.Ext(name) == "" {
		return name + ext
	}

	return name
}
