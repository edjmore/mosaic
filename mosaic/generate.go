package mosaic

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/edjmore/mosaic/kdtree"
	"github.com/edjmore/mosaic/util"
	"github.com/nfnt/resize"
)

// Generate a mosaic matching tgt using the images in workdir as "pixels".
// pxSize is the dimension of the square images in workdir.
// size x size squares of pixels from tgt will be replaced by pxSize x pxSize images.
// The dimensions of the resulting image will be w*pxSize/size, h*pxSize/size.
func Generate(tgt image.Image, workdir string, size, pxSize int) (image.Image, error) {
	files, err := ioutil.ReadDir(workdir)
	if err != nil {
		return nil, err
	}
	start := time.Now()

	// Build a "palette" from the images in workdir.
	// pathByColor gives us a mapping from a color to the corresponding image.
	palette := kdtree.New()
	pathByColor := make(map[color.Color]string)
	for _, file := range files {

		path := filepath.Join(workdir, file.Name())
		img, err := loadJpeg(path)
		if err != nil {
			return nil, fmt.Errorf("error loading %q: %v", path, err)
		}

		c := util.ComputeAvgColor(img, img.Bounds())
		palette.Add(c)
		pathByColor[c] = path
	}
	log.Printf("built palette: %v\n", time.Since(start))

	// Compute an average pixel color for each size x size square.
	// Then find the workdir image with nearest color, for each pixel in pix.
	pix := util.Pixelate(tgt, size)
	coordsByPath := make(map[string][]int)
	w, h := pix.Bounds().Max.X, pix.Bounds().Max.Y
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {

			c := pix.At(x, y)
			path := pathByColor[palette.Nearest(c)]

			coords, ok := coordsByPath[path]
			if !ok {
				coords = []int{}
			}
			coordsByPath[path] = append(coords, x*w+y)
		}
	}
	log.Printf("pixelated: %v\n", time.Since(start))

	// Create the final image.
	// For each workdir image in coordsByPath, copy the image to each coord in out.
	out := image.NewRGBA(image.Rect(0, 0, w*pxSize, h*pxSize))
	var wg sync.WaitGroup
	for path, coords := range coordsByPath {

		wg.Add(1)
		go func(path string, coords []int) {
			defer wg.Done()

			img, err := loadJpeg(path)
			if err != nil {
				log.Printf("error loading %q: %v", path, err)
				return
			}

			// Pre-processed images will be square, but may not be the right size.
			if img.Bounds().Max.X != pxSize {
				img = resize.Resize(uint(pxSize), uint(pxSize), img, resize.NearestNeighbor)
				if err != nil {
					log.Printf("error resizing %q: %v", path, err)
				}
			}

			for _, coord := range coords {
				x, y := coord/w, coord%w

				// Copy img to out. Upper left of img will be at (x*pxSize, y*pxSize) in out.
				for xx := 0; xx < pxSize; xx++ {
					for yy := 0; yy < pxSize; yy++ {
						c := img.At(xx, yy)
						out.Set(x*pxSize+xx, y*pxSize+yy, c)
					}
				}
			}
		}(path, coords)
	}
	wg.Wait()
	log.Printf("generated mosaic: %v\n", time.Since(start))
	return out, nil
}

func loadJpeg(path string) (image.Image, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return jpeg.Decode(f)
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
