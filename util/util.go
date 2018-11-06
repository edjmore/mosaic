package util

import (
	"crypto/sha256"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"os"
	"strings"

	"github.com/edjmore/mosaic/tifig"
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

// Load an image file based on file extension.
// This function should work for "jpeg", "jpg", "png", or "heic" files.
func LoadImage(path string) (image.Image, error) {
	pathLower := strings.ToLower(path)

	// Choose a decode func based on file extension.
	var decode func(io.Reader) (image.Image, error)
	if strings.HasSuffix(pathLower, ".jpeg") || strings.HasSuffix(pathLower, ".jpg") {
		decode = jpeg.Decode
	} else if strings.HasSuffix(pathLower, ".png") {
		decode = png.Decode
	} else if strings.HasSuffix(pathLower, ".heic") {
		// For HEIF encoded files, we try to convert to JPEG with tifig.
		outpath := path + ".jpeg"
		if err := tifig.Convert(path, outpath); err != nil {
			return nil, err
		}

		// Remove the conversion result after loading image.
		defer os.Remove(outpath)
		path = outpath
		decode = jpeg.Decode
	} else {
		return nil, fmt.Errorf("unknown image format: %q", path)
	}

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return decode(f)
}
