package raytrace

import (
	"image"
)

type PixelBuffer struct {
	width  int
	height int
	isYUp  bool
	img    *image.RGBA
}

func NewPixelBuffer(width int, height int) *PixelBuffer {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	return &PixelBuffer{width, height, true, img}
}

func (buf *PixelBuffer) Width() int {
	return buf.width
}

func (buf *PixelBuffer) Height() int {
	return buf.height
}

func (buf *PixelBuffer) SetPixelColor(x int, y int, clr ColorVector) {
	// Our Y axis is UP (right handed coordinate system)
	// X is right, and positive Z is out of the screen towards
	// the viewer.  So our calculated Y pixels are
	// the opposite direction of the Y in the image buffer.
	// If IsYUp is true then we'll invert Y when setting it into
	// the image.
	if buf.isYUp {
		y = buf.height - 1 - y
	}

	buf.img.Set(x, y, clr)
}

func (buf *PixelBuffer) GetImage() *image.RGBA {
	return buf.img
}
