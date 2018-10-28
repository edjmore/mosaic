package main

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/edjmore/mosaic/kdtree"
	"github.com/edjmore/mosaic/tifig"
	"github.com/edjmore/mosaic/util"
)

// Result image dimensions (w, h) are scaled by K/X
const (
	K = 30
	X = 20
)

func main() {
	// Create workdir to store intermediate/temp images.
	workdir, err := ioutil.TempDir("", "mosaic_")
	checkError(err)
	defer os.RemoveAll(workdir)

	dirpath, tgtpath, respath := os.Args[1], os.Args[2], os.Args[3]

	// Convert all images in dirpath to JPEGs (save to workdir).
	fmt.Printf("\nconverting files in %q to JPEG\n", dirpath)
	convertInputDir(dirpath, workdir)

	// Convert tgt image to JPEG.
	err = tifig.ConvertAndResize(tgtpath, respath, 4000, 4000)
	checkError(err)

	// Pixelate the tgt image.
	fmt.Printf("\npixelating %q\n", tgtpath)
	pix := pixelateTarget(tgtpath, respath)

	// Build kdree palette from converted input images.
	fmt.Println("\nbuilding color palette")
	pal, colormap := buildPalette(workdir)

	// Create the mosaic by choosing the best color/image for each pixel.
	fmt.Println("\ncreating mosaic")
	out := createMosaic(pix, pal, colormap)
	saveJpeg(out, respath)
}

func convertInputDir(indir, workdir string) {
	defer timeit("input dir conversion to JPEG")()

	files, err := ioutil.ReadDir(indir)
	checkError(err)

	var wg sync.WaitGroup
	for _, file := range files {

		wg.Add(1)
		go func(file os.FileInfo) {
			defer wg.Done()

			path := filepath.Join(indir, file.Name())
			outpath := filepath.Join(workdir, strings.TrimSuffix(file.Name(), "heic")+"jpeg")
			if err = tifig.ConvertAndResize(path, outpath, K, K); err != nil {
				fmt.Printf("error converting %q: %v", path, err)
			}
		}(file)
	}
	wg.Wait()
}

func pixelateTarget(tgtpath, respath string) image.Image {
	defer timeit("pixelated target image")()

	tgt := loadJpeg(respath)
	pix := util.Pixelate(tgt, X)
	saveJpeg(pix, respath)
	return pix
}

func buildPalette(workdir string) (*kdtree.Kdtree, map[color.Color]string) {
	defer timeit("built color palette (kdtree)")()

	t := kdtree.New()
	m := make(map[color.Color]string)

	files, err := ioutil.ReadDir(workdir)
	checkError(err)
	for _, file := range files {
		path := filepath.Join(workdir, file.Name())
		im := loadJpeg(path)

		c := util.ComputeAvgColor(im, im.Bounds())
		t.Add(c)
		m[c] = path
	}
	return t, m
}

func createMosaic(pix image.Image, pal *kdtree.Kdtree, colormap map[color.Color]string) image.Image {
	defer timeit("created mosaic")()

	// Filepaths for each pixel.
	grid := make([][]string, 0)
	usedPaths := make(map[string]bool)

	b := pix.Bounds()
	out := image.NewRGBA(image.Rect(0, 0, (b.Max.X-b.Min.X)*K, (b.Max.Y-b.Min.Y)*K))
	for x := 0; x < b.Max.X; x++ {
		grid = append(grid, make([]string, 0))

		for y := 0; y < b.Max.Y; y++ {
			// Find the input image that best matches this color.
			bestMatch := pal.Nearest(pix.At(x, y))
			path := colormap[bestMatch]
			grid[x] = append(grid[x], path)
			usedPaths[path] = true
		}
	}

	// Only load each matching image once.
	var wg sync.WaitGroup
	for path := range usedPaths {

		wg.Add(1)
		go func(path string) {
			defer wg.Done()
			im := loadJpeg(path)

			for x := 0; x < b.Max.X; x++ {
				for y := 0; y < b.Max.Y; y++ {

					if grid[x][y] != path {
						continue
					}

					// Copy matching image to output image.
					for xx := 0; xx < K; xx++ {
						for yy := 0; yy < K; yy++ {
							out.Set(x*K+xx, y*K+yy, im.At(xx, yy))
						}
					}
				}
			}
		}(path)
	}
	wg.Wait()
	return out
}

func loadJpeg(path string) image.Image {
	f, err := os.Open(path)
	checkError(err)
	im, err := jpeg.Decode(f)
	f.Close()
	checkError(err)
	return im
}

func saveJpeg(im image.Image, path string) {
	w, err := os.Create(path)
	checkError(err)
	jpeg.Encode(w, im, nil)
	w.Close()
	checkError(err)
}

func timeit(msg string) func() {
	start := time.Now()
	return func() {
		fmt.Printf("%s: %s\n", msg, time.Since(start))
	}
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
