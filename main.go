package main

import (
	"fmt"
	"runtime"

	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	verbose = kingpin.Flag("debug", "Enable debug mode.").Short('v').Bool()
)

func maxParallelism() int {
	maxProcs := runtime.GOMAXPROCS(0)
	numCPU := runtime.NumCPU()
	if maxProcs < numCPU {
		return maxProcs
	}
	return numCPU
}

func doRender() {
	numThreads := maxParallelism()
	imageWidth := 10
	imageHeight := 5
	isTwoPhase := false
	maxDepth := 50
	numSamplesPerPixel := 10

	var renderConfig = NewRenderConfig(numThreads, maxDepth, numSamplesPerPixel, isTwoPhase)
	var scene = createCornellBoxScene()
	var render = newRenderer(renderConfig)
	var pixelBuffer = newPixelBuffer(imageWidth, imageHeight)
	render.Render(pixelBuffer, scene, renderConfig)
}

func main() {
	kingpin.Parse()
	fmt.Printf("%v\n", *verbose)

	glmain()
	// doRender()
}
