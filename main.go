package main

import (
	"github.com/markzuber/zgotrace/raytrace"
	"github.com/markzuber/zgotrace/raytrace/scenes"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	verbose = kingpin.Flag("debug", "Enable debug mode.").Short('v').Bool()
)

func main() {
	kingpin.Parse()
	// fmt.Printf("%v\n", *verbose)

	// glmain()
	// doRender()

	ebRen := newEbitenRender(200, 200, 4.0)

	numThreads := 1 // maxParallelism()
	isTwoPhase := false
	maxDepth := 50
	numSamplesPerPixel := 1000

	var renderConfig = raytrace.NewRenderConfig(numThreads, maxDepth, numSamplesPerPixel, isTwoPhase)
	var scene = scenes.CreateCornellBoxScene()

	// var rayTracer = raytrace.NewSimpleTracer()
	var rayTracer = raytrace.NewMonteCarloTracer()

	go ebRen.doRender(rayTracer, renderConfig, scene)

	ebRen.startDisplay()
}
