//go:build !windows

package chat

func protectKey(raw []byte) ([]byte, error) {
	return append([]byte(nil), raw...), nil
}

func unprotectKey(protected []byte) ([]byte, error) {
	return append([]byte(nil), protected...), nil
}
