package util

import (
	"image"
	"image/color"
)

func Pixelate(im image.Image, k int) image.Image {
	// Init output image with same dimensions as input.
	b := im.Bounds()
	bOut := image.Rect(0, 0, b.Max.X-b.Min.X, b.Max.Y-b.Min.Y)
	out := image.NewRGBA(bOut)

	for x := b.Min.X; x < b.Max.X; x += k {
		for y := b.Min.Y; y < b.Max.Y; y += k {

			// Compute average color of kxk grid.
			var rs, gs, bs, pxCount uint32
			for xx := x; xx < x+k && xx < b.Max.X; xx++ {
				for yy := y; yy < y+k && yy < b.Max.Y; yy++ {
					r, g, b, _ := im.At(xx, yy).RGBA()
					rs += r
					gs += g
					bs += b
					pxCount++
				}
			}
			avgColor := color.RGBA{
				uint8(rs / pxCount / 0x101),
				uint8(gs / pxCount / 0x101),
				uint8(bs / pxCount / 0x101),
				0xff,
			}

			// Set output image pixels in grid to average color.
			for xx := x; xx < x+k && xx < b.Max.X; xx++ {
				for yy := y; yy < y+k && yy < b.Max.Y; yy++ {
					out.Set(xx-b.Min.X, yy-b.Min.Y, avgColor)
				}
			}
		}
	}
	return out
}
