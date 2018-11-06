package preprocess

import (
	"image"
	"image/jpeg"
	"os"
	"path/filepath"
	"strings"

	"github.com/edjmore/mosaic/tifig"
	"github.com/edjmore/mosaic/util"
	"github.com/nfnt/resize"
	"github.com/oliamb/cutter"
)

// Pre-process an image so that it can be used in a mosaic.
// The image can be a JPEG, PNG, or HEIF encoded.
// The processed image will be saved in JPEG format to workdir.
// Image will be resized - and cropped if necessary - to be a square of (size, size).
func ImageFile(filename, workdir string, size int) (string, error) {
	_, name := filepath.Split(filename)
	outpath := filepath.Join(workdir, name+".jpeg")

	// For HEIF encoded images, the tifig executable handles all the pre-processing.
	if strings.HasSuffix(filename, ".heic") {
		err := tifig.ConvertAndResize(filename, outpath, size, size)
		return outpath, err
	}

	img, err := util.LoadImage(filename)
	if err != nil {
		return "", err
	}
	return outpath, preprocessImageFile(img, outpath, size)
}

func preprocessImageFile(img image.Image, outpath string, size int) error {
	// If img is not square, need to crop larger dimension before resizing.
	b := img.Bounds()
	w, h := b.Max.X-b.Min.X, b.Max.Y-b.Min.Y
	if w != h {
		d := w
		if d > h {
			d = h
		}

		var err error
		img, err = cutter.Crop(img, cutter.Config{
			Width:  d,
			Height: d,
			Mode:   cutter.Centered,
		})
		if err != nil {
			return err
		}
	}

	// Using nearest-neighbor interpolation b/c it's fast.
	img = resize.Resize(uint(size), uint(size), img, resize.NearestNeighbor)

	out, err := os.Create(outpath)
	if err != nil {
		return err
	}
	err = jpeg.Encode(out, img, nil)
	out.Close()
	return err
}
