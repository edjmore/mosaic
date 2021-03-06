package main

import (
	"context"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/edjmore/mosaic/mosaic"
	"github.com/edjmore/mosaic/preprocess"
	"github.com/edjmore/mosaic/tifig"
	"golang.org/x/sync/semaphore"
)

const (
	PREPROCESS_SIZE = 50
)

func main() {
	// TODO: there must be a better way to parse arguments (e.g. Python's argparse)
	imgdir := os.Args[1]

	// Pre-processed image files will be stored in workdir.
	workdir, err := ioutil.TempDir("", "mosaic_")
	checkError(err)
	defer os.RemoveAll(workdir)

	files, err := ioutil.ReadDir(imgdir)
	checkError(err)

	var wg sync.WaitGroup
	ctx := context.TODO()
	sem := semaphore.NewWeighted(32) // limit on num open file descriptors

	for _, file := range files {

		if err = sem.Acquire(ctx, 1); err != nil {
			log.Printf("error aquiring semaphore for %q: %v", file.Name(), err)
			continue
		}

		// Start a new goroutine for each image file so we can process concurrently.
		wg.Add(1)
		go func(file os.FileInfo) {
			defer wg.Done()
			defer sem.Release(1)

			path := filepath.Join(imgdir, file.Name())
			if _, err := preprocess.ImageFile(path, workdir, PREPROCESS_SIZE); err != nil {
				log.Printf("error processing %q: %v", path, err)
			}
		}(file)
	}
	wg.Wait()

	// infinite command loop
	for {
		fmt.Printf("> ")

		// Parse user command.
		var tgtpath string
		var size, pxSize int
		_, err := fmt.Fscanf(os.Stdin, "%s %d %d\n", &tgtpath, &size, &pxSize)
		checkError(err)

		// Load the target image.
		tgt, err := loadImage(tgtpath)
		checkError(err)

		// Generate a mosaic.
		out, err := mosaic.Generate(tgt, workdir, size, pxSize)
		checkError(err)

		// Save the mosaic as a JPEG.
		w, err := os.Create("mosaic.jpeg")
		checkError(err)
		err = jpeg.Encode(w, out, nil)
		w.Close()
		checkError(err)
	}
}

func loadImage(path string) (image.Image, error) {
	if strings.HasSuffix(path, ".heic") {

		outpath := path + ".jpeg"
		size := 10 * 1000 // tifig won't actually crop if image is smaller

		err := tifig.ConvertAndResize(path, outpath, size, size)
		if err != nil {
			return nil, err
		}

		defer os.Remove(outpath)
		path = outpath
	}

	var decode func(io.Reader) (image.Image, error)
	if strings.HasSuffix(path, ".jpeg") || strings.HasSuffix(path, ".jpg") {
		decode = jpeg.Decode
	} else if strings.HasSuffix(path, ".png") {
		decode = png.Decode
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

func checkError(err error) {
	if err != nil {
		log.Panicln(err)
	}
}
