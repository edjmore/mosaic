package util_test

import (
	"image/jpeg"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/edjmore/mosaic/util"
)

func TestPixelate(t *testing.T) {
	wd, err := os.Getwd()
	checkError(err)
	inpath := filepath.Join(wd, "testdata", "yellow_flowers.jpeg")

	// Load input image.
	f, err := os.Open(inpath)
	checkError(err)
	im, err := jpeg.Decode(f)
	f.Close()
	checkError(err)

	// Create tempdir for the pixelated file.
	outdir, err := ioutil.TempDir("", "pixelate_test_")
	checkError(err)
	defer os.RemoveAll(outdir)
	outpath := filepath.Join(outdir, "yellow_flowers.jpeg")

	// Pixelate and save output image.
	out := util.Pixelate(im, 64)
	w, err := os.Create(outpath)
	checkError(err)
	err = jpeg.Encode(w, out, nil)
	w.Close()
	checkError(err)

	// Verify that the output image matches pre-computed checksum.
	checksum, err := util.ComputeChecksum(outpath)
	checkError(err)
	expectedChecksum := "\xaa\xeaqm\x93E\xec@\xa3\x8d\xefJ1\xffEk\xb6(\xe4\x94\xcf`\xbf,\xea\xfc\xc4<\xd0zqz"
	if checksum != expectedChecksum {
		t.Errorf("checksums don't match: expected %q, but was %q", expectedChecksum, checksum)
	}
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
