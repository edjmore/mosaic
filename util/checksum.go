package util

import (
	"crypto/sha256"
	"io"
	"os"
)

// Returns the file's SHA256 checksum.
func ComputeChecksum(filepath string) (string, error) {
	f, err := os.Open(filepath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := sha256.New()
	if _, err = io.Copy(h, f); err != nil {
		return "", err
	}
	return string(h.Sum(nil)), nil
}
