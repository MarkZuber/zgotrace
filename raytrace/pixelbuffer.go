package raytrace

import (
	"image"
)

type PixelBuffer struct {
	width  int
	height int
	img    *image.RGBA
}

func NewPixelBuffer(width int, height int) *PixelBuffer {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	return &PixelBuffer{width, height, img}
}

func (buf *PixelBuffer) Width() int {
	return buf.width
}

func (buf *PixelBuffer) Height() int {
	return buf.height
}

func (buf *PixelBuffer) SetPixelColor(x int, y int, clr ColorVector) {
	// r, g, b, a := clr.RGBA()
	// fmt.Printf("SETTING PIXEL COLOR %v, %v, %v, %v\n", r, g, b, a)
	buf.img.Set(x, y, clr)
	// fmt.Printf("Just set: %v\n", buf.img.At(x, y))
}

func (buf *PixelBuffer) GetImage() *image.RGBA {
	return buf.img
}
