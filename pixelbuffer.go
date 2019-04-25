package main

type pixelBuffer struct {
	width  int
	height int
}

func newPixelBuffer(width int, height int) *pixelBuffer {
	return &pixelBuffer{width, height}
}

func (buf *pixelBuffer) Width() int {
	return buf.width
}

func (buf *pixelBuffer) Height() int {
	return buf.height
}

func (buf *pixelBuffer) SetPixelColor(x int, y int, color colorVector) {

}
