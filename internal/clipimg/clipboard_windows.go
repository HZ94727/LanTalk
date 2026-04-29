//go:build windows

package clipimg

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"runtime"
	"time"
	"unsafe"

	"LanTalk/internal/mediautil"
	"golang.org/x/sys/windows"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	_ "golang.org/x/image/webp"
)

const (
	cfDIBV5      = 17
	gmemMoveable = 0x0002
	biRGB        = 0
)

type bitmapV5Header struct {
	Size          uint32
	Width         int32
	Height        int32
	Planes        uint16
	BitCount      uint16
	Compression   uint32
	SizeImage     uint32
	XPelsPerMeter int32
	YPelsPerMeter int32
	ClrUsed       uint32
	ClrImportant  uint32
	RedMask       uint32
	GreenMask     uint32
	BlueMask      uint32
	AlphaMask     uint32
	CSType        uint32
	Endpoints     struct {
		CiexyzRed, CiexyzGreen, CiexyzBlue struct {
			CiexyzX, CiexyzY, CiexyzZ int32
		}
	}
	GammaRed    uint32
	GammaGreen  uint32
	GammaBlue   uint32
	Intent      uint32
	ProfileData uint32
	ProfileSize uint32
	Reserved    uint32
}

var (
	user32            = windows.NewLazySystemDLL("user32.dll")
	kernel32          = windows.NewLazySystemDLL("kernel32.dll")
	procOpenClipboard = user32.NewProc("OpenClipboard")
	procCloseClipboard = user32.NewProc("CloseClipboard")
	procEmptyClipboard = user32.NewProc("EmptyClipboard")
	procSetClipboardData = user32.NewProc("SetClipboardData")
	procGlobalAlloc   = kernel32.NewProc("GlobalAlloc")
	procGlobalLock    = kernel32.NewProc("GlobalLock")
	procGlobalUnlock  = kernel32.NewProc("GlobalUnlock")
	procGlobalFree    = kernel32.NewProc("GlobalFree")
	procMemMove       = kernel32.NewProc("RtlMoveMemory")
)

func WriteImageDataURL(dataURL string) error {
	raw, _, err := mediautil.DecodeImageDataURL(dataURL)
	if err != nil {
		return err
	}

	img, _, err := image.Decode(bytes.NewReader(raw))
	if err != nil {
		return fmt.Errorf("decode image: %w", err)
	}

	dib, err := encodeClipboardImage(img)
	if err != nil {
		return err
	}

	if err := writeDIBV5ToClipboard(dib); err != nil {
		return err
	}

	return nil
}

func encodeClipboardImage(img image.Image) ([]byte, error) {
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()
	if width <= 0 || height <= 0 {
		return nil, errors.New("image has invalid bounds")
	}

	offset := int(unsafe.Sizeof(bitmapV5Header{}))
	imageSize := 4 * width * height
	header := bitmapV5Header{
		Size:        uint32(offset),
		Width:       int32(width),
		Height:      int32(height),
		Planes:      1,
		BitCount:    32,
		Compression: biRGB,
		SizeImage:   uint32(imageSize),
		RedMask:     0x00ff0000,
		GreenMask:   0x0000ff00,
		BlueMask:    0x000000ff,
		AlphaMask:   0xff000000,
		CSType:      0x73524742,
		Intent:      4,
	}

	dib := make([]byte, offset+imageSize)
	headerBytes := unsafe.Slice((*byte)(unsafe.Pointer(&header)), offset)
	copy(dib[:offset], headerBytes)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			srcX := bounds.Min.X + x
			srcY := bounds.Min.Y + (height - 1 - y)
			idx := offset + 4*(y*width+x)
			r, g, b, a := img.At(srcX, srcY).RGBA()
			dib[idx+2] = uint8(r >> 8)
			dib[idx+1] = uint8(g >> 8)
			dib[idx] = uint8(b >> 8)
			dib[idx+3] = uint8(a >> 8)
		}
	}

	return dib, nil
}

func writeDIBV5ToClipboard(dib []byte) error {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	if err := openClipboardWithRetry(); err != nil {
		return err
	}
	defer procCloseClipboard.Call()

	if r, _, err := procEmptyClipboard.Call(); r == 0 {
		return fmt.Errorf("empty clipboard: %w", err)
	}

	hMem, _, err := procGlobalAlloc.Call(gmemMoveable, uintptr(len(dib)))
	if hMem == 0 {
		return fmt.Errorf("allocate clipboard memory: %w", err)
	}

	ptr, _, err := procGlobalLock.Call(hMem)
	if ptr == 0 {
		procGlobalFree.Call(hMem)
		return fmt.Errorf("lock clipboard memory: %w", err)
	}

	procMemMove.Call(ptr, uintptr(unsafe.Pointer(&dib[0])), uintptr(len(dib)))
	procGlobalUnlock.Call(hMem)

	if r, _, err := procSetClipboardData.Call(cfDIBV5, hMem); r == 0 {
		procGlobalFree.Call(hMem)
		return fmt.Errorf("set clipboard data: %w", err)
	}

	return nil
}

func openClipboardWithRetry() error {
	var lastErr error
	for range 8 {
		if r, _, err := procOpenClipboard.Call(0); r != 0 {
			return nil
		} else {
			lastErr = err
		}
		time.Sleep(20 * time.Millisecond)
	}
	if lastErr == nil {
		lastErr = errors.New("clipboard is busy")
	}
	return fmt.Errorf("open clipboard: %w", lastErr)
}
