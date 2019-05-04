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

	ebRen := newEbitenRender(600, 600, 2.0)

	isTwoPhase := true
	maxDepth := 50
	numSamplesPerPixel := 100

	var renderConfig = raytrace.NewRenderConfig(maxDepth, numSamplesPerPixel, isTwoPhase)
	// var scene = scenes.CreateCornellBoxScene()
	var scene = scenes.CreateCornellBoxWithSmokeScene()
	// var scene = scenes.CreateManySpheresScene()
	var rayTracer = raytrace.NewMonteCarloTracer()

	go ebRen.doRender(rayTracer, renderConfig, scene)

	ebRen.startDisplay()
}
