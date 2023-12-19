package random

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

func Bytes(n int) ([]byte, error) {
	bytes := make([]byte, n)
	nRead, err := rand.Read(bytes)
	if err != nil {
		return nil, fmt.Errorf("error reading bytes: %w", err)
	}

	if nRead < n {
		return nil, fmt.Errorf("didn't read enough bytes")
	}

	return bytes, nil
}

func String(n int) (string, error) {
	bytes, err := Bytes(n)
	if err != nil {
		return "", fmt.Errorf("error stringifying bytes: %w", err)
	}

	return base64.URLEncoding.EncodeToString(bytes), nil
}

func SessionToken(bytesPerToken int) (string, error) {
	return String(bytesPerToken)
}
