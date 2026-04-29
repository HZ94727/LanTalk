//go:build !windows

package clipimg

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"sync"

	"LanTalk/internal/mediautil"
	"golang.design/x/clipboard"
	_ "golang.org/x/image/webp"
)

var (
	initOnce sync.Once
	initErr  error
)

func initClipboard() error {
	initOnce.Do(func() {
		initErr = clipboard.Init()
	})
	return initErr
}

func WriteImageDataURL(dataURL string) error {
	if err := initClipboard(); err != nil {
		return err
	}

	raw, _, err := mediautil.DecodeImageDataURL(dataURL)
	if err != nil {
		return err
	}

	img, _, err := image.Decode(bytes.NewReader(raw))
	if err != nil {
		return fmt.Errorf("decode image: %w", err)
	}

	var pngBytes bytes.Buffer
	if err := png.Encode(&pngBytes, img); err != nil {
		return fmt.Errorf("encode clipboard image: %w", err)
	}

	clipboard.Write(clipboard.FmtImage, pngBytes.Bytes())
	return nil
}
