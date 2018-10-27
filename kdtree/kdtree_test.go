package kdtree_test

import (
    "image/color"
    "image/color/palette"
    "testing"
    
    "github.com/edjmore/mosaic/kdtree"
)

func TestKdtree(t *testing.T) {
    // Init tree with palette (216 colors).
    k := kdtree.New()
    var pal color.Palette = palette.WebSafe
    for _, c := range pal {
        k.Add(c)
    }
    
    // Generate some colors to use as input for Nearest(). 
    // The Kdtree result should match palette's Convert() method.
    for r := 0; r < 256; r += 10 {
        for g := 0; g < 256; g += 10 {
            for b := 0; b < 256; b += 10 {
                c := color.RGBA{uint8(r), uint8(g), uint8(b), 0xff}
                expected := pal.Convert(c)
                actual := k.Nearest(c)
                if actual != expected {
                    t.Fatalf("Nearest(%v): expected %v, but got %v", c, expected, actual)
                }
            }
        }
    }
}
