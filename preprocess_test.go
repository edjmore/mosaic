package mosaic_test

import (
	"crypto/sha256"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/edjmore/mosaic"
)

func TestPreprocessJpeg(t *testing.T) {
	wd, err := os.Getwd()
	checkError(err)
	filename := filepath.Join(wd, "testdata", "gopher.jpeg")

	// Create tempdir for the output file.
	outdir, err := ioutil.TempDir("", "tifig_test_")
	checkError(err)
	defer os.RemoveAll(outdir)

	outpath, err := mosaic.PreprocessImageFile(filename, outdir, 200, 200)
	if err != nil {
		t.Error(err)
	}

	// Verify that the processed image matches pre-computed checksum.
	checksum, err := computeChecksum(outpath)
	checkError(err)
	expectedChecksum := "\x9a\xd8.\x19\xbf\xcb\x01!\x88J-\xf2\x1f\xb6ө\xaa\xb3\x01غvsݥ\xd1A;\xa5P\xc8\x05"
	if checksum != expectedChecksum {
		t.Errorf("checksums don't match: expected %q, but was %q", expectedChecksum, checksum)
	}
}

// Returns the file's SHA256 checksum.
func computeChecksum(filepath string) (string, error) {
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

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
