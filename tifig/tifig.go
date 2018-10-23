package tifig

import (
	"fmt"
	"os/exec"
	"strconv"
)

// Convert an HEIF encoded image to JPEG and save to outpath.
// The image will be resized and cropped to width, height.
// If necessary, the result image will be cropped to match provided aspect ratio.
func ConvertAndResize(inpath, outpath string, width, height int) error {
	cmd := exec.Command(
		"tifig",
		inpath,
		outpath,
		"--crop",
		"-w", strconv.Itoa(width),
		"-h", strconv.Itoa(height),
	)
	out, err := cmd.CombinedOutput()
	if err != nil {
		err = fmt.Errorf("%v: %s", err, out)
	}
	return err
}
