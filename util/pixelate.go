package util

import (
	"image"
	"image/color"
	"sync"
)

func Pixelate(im image.Image, k int) image.Image {
	// Output image will have one pixel for each kxk grid from input.
	b := im.Bounds()
	w, h := (b.Max.X-b.Min.X)/k, (b.Max.Y-b.Min.Y)/k
	out := image.NewRGBA(image.Rect(0, 0, w, h))

	var wg sync.WaitGroup
	for x := b.Min.X; x < b.Max.X; x += k {
		for y := b.Min.Y; y < b.Max.Y; y += k {

			// Compute average pixel color for grid and add to output.
			wg.Add(1)
			go func(x, y int) {
				defer wg.Done()

				avgColor := ComputeAvgColor(im, image.Rect(x, y, x+k, y+k))
				out.Set(x/k, y/k, avgColor)
			}(x, y)
		}
	}

	wg.Wait()
	return out
}

// Compute the average pixel color within grid.
func ComputeAvgColor(im image.Image, grid image.Rectangle) color.Color {
	b := im.Bounds()

	var rs, gs, bs uint32
	for x := grid.Min.X; x < grid.Max.X && x < b.Max.X; x++ {
		for y := grid.Min.Y; y < grid.Max.Y && y < b.Max.Y; y++ {

			r, g, b, _ := im.At(x, y).RGBA()
			rs += r
			gs += g
			bs += b
		}
	}

	pxCount := uint32((grid.Max.X - grid.Min.X) * (grid.Max.Y - grid.Min.Y))
	return color.RGBA{
		uint8(rs / pxCount / 0x101),
		uint8(gs / pxCount / 0x101),
		uint8(bs / pxCount / 0x101),
		0xff,
	}
}
