package tifig_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/edjmore/mosaic/tifig"
	"github.com/edjmore/mosaic/util"
)

func TestConvertAndResize(t *testing.T) {
	wd, err := os.Getwd()
	checkError(err)
	inpath := filepath.Join(wd, "testdata", "lake_michigan.heic")

	// Create tempdir for the converted file.
	outdir, err := ioutil.TempDir("", "tifig_test_")
	checkError(err)
	defer os.RemoveAll(outdir)
	outpath := filepath.Join(outdir, "lake_michigan.jpeg")

	if err = tifig.ConvertAndResize(inpath, outpath, 200, 200); err != nil {
		t.Error(err)
	}

	// Verify that the output image matches pre-computed checksum.
	checksum, err := util.ComputeChecksum(outpath)
	checkError(err)
	expectedChecksum := "\x19\x15\xd2f\x92@ݙ\x0f\xf4M\xcb\x01ݖ\r\xbbo\x13\xa2Ow\xde\xe8\v\x96\xff\xe6a\xfd\x82Z"
	if checksum != expectedChecksum {
		t.Errorf("checksums don't match: expected %q, but was %q", expectedChecksum, checksum)
	}
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
