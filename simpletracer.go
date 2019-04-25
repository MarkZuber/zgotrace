package main

type simpleTracer struct {
}

func newSimpleTracer() *simpleTracer {
	return &simpleTracer{}
}

func (t *simpleTracer) GetPixelColor(x int, y int) colorVector {
	return colorVector{0.5, 0.6, 0.7}
}
