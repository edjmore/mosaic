package main

import (
	"fmt"
	"image/jpeg"
	"os"

	"github.com/edjmore/mosaic/tifig"
	"github.com/edjmore/mosaic/util"
)

func main() {
	inpath, outpath := os.Args[1], os.Args[2]

	// Convert from HEIF to JPEG.
	fmt.Printf("converting %q to JPEG\n", inpath)
	if err := tifig.ConvertAndResize(inpath, outpath, 300, 300); err != nil {
		panic(err)
	}

	// Load converted image.
	f, err := os.Open(outpath)
	checkError(err)
	im, err := jpeg.Decode(f)
	f.Close()
	checkError(err)

	// Pixelate and overwrite converted image.
	fmt.Printf("pixelating image %q\n", outpath)
	pix := util.Pixelate(im, 3)
	w, err := os.Create(outpath)
	checkError(err)
	err = jpeg.Encode(w, pix, nil)
	checkError(err)
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
