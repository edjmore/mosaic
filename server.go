package main

import (
	"fmt"
	"github.com/edjmore/mosaic/preprocess"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

const (
	PREPROCESS_SIZE = 50
)

// Pre-processed image files will be stored in workdir.
var workdir string

func main() {
	var err error
	workdir, err = ioutil.TempDir("", "images_")
	checkError(err)
	defer os.RemoveAll(workdir)

	// Setup a file server so we can serve images from the workdir.
	fs := http.FileServer(http.Dir(workdir))
	http.Handle(workdir + "/", http.StripPrefix(workdir, fs))

    http.HandleFunc("/favicon.ico", http.NotFound) // no favicon 
    
    // Listen until terminated.
	http.HandleFunc("/", handler)
	log.Print("listening...")
	log.Panic(http.ListenAndServe(":8080", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "<html>")
	switch r.Method {
	case http.MethodPost:
		handlePost(w, r)
	default:
		handleGet(w, r)
	}
	fmt.Fprint(w, "</html>")
}

func handlePost(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	srcdir := r.Form.Get("srcdir")
	if srcdir == "" {
		log.Fatalf("invalid srcdir: %q", srcdir)
	}

	// Load the images from srcdir to workdir.
	files, err := ioutil.ReadDir(srcdir)
	checkError(err)
	var wg sync.WaitGroup
	for _, file := range files {

		// Start a new goroutine for each image file so we can process concurrently.
		wg.Add(1)
		go func(file os.FileInfo) {
			defer wg.Done()

			path := filepath.Join(srcdir, file.Name())
			if _, err := preprocess.ImageFile(path, workdir, PREPROCESS_SIZE); err != nil {
				log.Printf("error processing %q: %v", path, err)
			}
		}(file)
	}
	wg.Wait()

	// Return an img tag for each file in workdir.
	files, err = ioutil.ReadDir(workdir)
	checkError(err)
	for _, file := range files {
		path := filepath.Join(workdir, file.Name())
		fmt.Fprintf(w, "<div><img src='%s'/></div>", path)
	}
}

func handleGet(w http.ResponseWriter, r *http.Request) {
	homedir := os.Getenv("HOME")
	currdir := filepath.Join(homedir, r.URL.Path)
	fmt.Fprintf(w, "<b>%s</b>:", currdir)

	files, err := ioutil.ReadDir(currdir)
	checkError(err)
	imgCount := 0
	fmt.Fprint(w, "<ul>")
	for _, file := range files {
		if strings.HasPrefix(file.Name(), ".") {
			continue // skip hidden files
		}

		if file.IsDir() {
			// Directories are links to view the directory.
			path := filepath.Join(currdir, file.Name())
			path = path[len(homedir):]
			fmt.Fprintf(w, "<li><a href=%s>%s</a></li>", path, file.Name())
		} else {
			// We show all files, but only count ones that look like images.
			if strings.HasSuffix(file.Name(), ".heic") || strings.HasSuffix(file.Name(), ".jpeg") || strings.HasSuffix(file.Name(), ".png") {
				imgCount++
			}
			fmt.Fprintf(w, "<li>%s</li>", file.Name())
		}

	}
	fmt.Fprint(w, "</ul>")

	if imgCount > 0 {
		// This folder has images, so it can be chosen as srcdir.
		fmt.Fprint(w, "<form action='/' method='post'>")
		fmt.Fprintf(w, "<input type='submit' name='srcdir' value='%s'/>", currdir)
		fmt.Fprint(w, "</form>")
	}
}

func checkError(err error) {
	if err != nil {
		log.Panicln(err)
	}
}
