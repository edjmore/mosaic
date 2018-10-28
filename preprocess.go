package mosaic

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/edjmore/mosaic/tifig"
	"github.com/nfnt/resize"
	"github.com/oliamb/cutter"
)

// Pre-process an image so that it can be used in a mosaic.
// The image can be a JPEG, PNG, or HEIF encoded.
// The processed image will be saved in JPEG format to workdir.
// Image will be resized - and cropped if necessary - to match width, height.
func PreprocessImageFile(filename, workdir string, width, height int) (string, error) {
	_, name := filepath.Split(filename)
	outpath := filepath.Join(workdir, name+".jpeg")

	// For HEIF encoded images, the tifig executable handles all the pre-processing.
	if strings.HasSuffix(filename, ".heic") {
		err := tifig.ConvertAndResize(filename, outpath, width, height)
		return outpath, err
	}

	if strings.HasSuffix(filename, ".png") {
		return outpath, preprocessImageFile(filename, outpath, width, height, png.Decode)
	} else if strings.HasSuffix(filename, ".jpeg") {
		return outpath, preprocessImageFile(filename, outpath, width, height, jpeg.Decode)
	} else {
		return "", fmt.Errorf("unknown image format: %q", filename)
	}
}

func preprocessImageFile(filename, outpath string, width, height int, decode func(io.Reader) (image.Image, error)) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}

	img, err := decode(f)
	f.Close()
	if err != nil {
		return err
	}

	// If img is not square, need to crop larger dimension before resizing.
	b := img.Bounds()
	w, h := b.Max.X-b.Min.X, b.Max.Y-b.Min.Y
	if w != h {
		d := w
		if d > h {
			d = h
		}
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
	img = resize.Resize(uint(width), uint(height), img, resize.NearestNeighbor)
	if err != nil {
		return err
	}

	out, err := os.Create(outpath)
	if err != nil {
		return err
	}
	err = jpeg.Encode(out, img, nil)
	out.Close()
	return err
}
