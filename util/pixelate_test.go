package util_test

import (
	"image/jpeg"
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
	defer f.Close()
	im, err := jpeg.Decode(f)
	checkError(err)

	// Pixelate and save output image.
	out := util.Pixelate(im, 64)
	w, err := os.Create(filepath.Join(wd, "testdata", "yellow_flowers_pix.jpeg"))
	checkError(err)
	defer w.Close()
	err = jpeg.Encode(w, out, nil)
	checkError(err)
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
