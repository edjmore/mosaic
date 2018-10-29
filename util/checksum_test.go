package util_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/edjmore/mosaic/util"
)

func TestComputeChecksum(t *testing.T) {
	wd, err := os.Getwd()
	checkError(err)
	path := filepath.Join(wd, "testdata", "yellow_flowers.jpeg")

	checksum, err := util.ComputeChecksum(path)
	if err != nil {
		t.Fatal(err)
	}
	expected := "fp\xd9uU\xc8U\xd0\xe7\xbf^\x9e\xe9\xec\x00\xaa\xeb\xe4\x84\xeb\xf5\xeb$*\x844}V\x8d\xa2\xbc\xe1"
	if checksum != expected {
		t.Fatalf("expected %q, but got %q", expected, checksum)
	}
}
