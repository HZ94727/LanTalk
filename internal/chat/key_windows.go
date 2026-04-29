//go:build windows

package chat

import (
	"fmt"
	"unsafe"

	"golang.org/x/sys/windows"
)

func protectKey(raw []byte) ([]byte, error) {
	in := windows.DataBlob{
		Size: uint32(len(raw)),
	}
	if len(raw) > 0 {
		in.Data = &raw[0]
	}

	description, err := windows.UTF16PtrFromString("LanTalk Key")
	if err != nil {
		return nil, fmt.Errorf("protect key description: %w", err)
	}

	var out windows.DataBlob
	if err := windows.CryptProtectData(&in, description, nil, 0, nil, 0, &out); err != nil {
		return nil, fmt.Errorf("protect key: %w", err)
	}
	defer windows.LocalFree(windows.Handle(unsafe.Pointer(out.Data)))

	protected := unsafe.Slice(out.Data, out.Size)
	return append([]byte(nil), protected...), nil
}

func unprotectKey(protected []byte) ([]byte, error) {
	in := windows.DataBlob{
		Size: uint32(len(protected)),
	}
	if len(protected) > 0 {
		in.Data = &protected[0]
	}

	var out windows.DataBlob
	if err := windows.CryptUnprotectData(&in, nil, nil, 0, nil, 0, &out); err != nil {
		return nil, fmt.Errorf("unprotect key: %w", err)
	}
	defer windows.LocalFree(windows.Handle(unsafe.Pointer(out.Data)))

	raw := unsafe.Slice(out.Data, out.Size)
	return append([]byte(nil), raw...), nil
}
