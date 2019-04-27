package main

import (
	"fmt"
	"log"
	"math/rand"
	"runtime"
	"time"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/markzuber/zgotrace/raytrace"
)

type ebitenRender struct {
	offscreen   *ebiten.Image
	pixelBuffer *raytrace.PixelBuffer
	screenScale float64
}

func newEbitenRender(imageWidth int, imageHeight int, screenScale float64) *ebitenRender {
	rand.Seed(time.Now().UnixNano())
	offscreen, _ := ebiten.NewImage(imageWidth, imageHeight, ebiten.FilterDefault)
	var pixelBuffer = raytrace.NewPixelBuffer(imageWidth, imageHeight)

	return &ebitenRender{offscreen, pixelBuffer, screenScale}
}

func (er *ebitenRender) update(screen *ebiten.Image) error {
	ebImg, _ := ebiten.NewImageFromImage(er.pixelBuffer.GetImage(), ebiten.FilterDefault)

	// fmt.Printf("%v\n", er.pixelBuffer.GetImage().At(0, 0))

	// todo: get rid of offscreen?
	er.offscreen.DrawImage(ebImg, nil)

	if ebiten.IsDrawingSkipped() {
		return nil
	}

	screen.DrawImage(er.offscreen, nil)
	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f", ebiten.CurrentTPS()))
	return nil
}

func maxParallelism() int {
	maxProcs := runtime.GOMAXPROCS(0)
	numCPU := runtime.NumCPU()
	if maxProcs < numCPU {
		return maxProcs
	}
	return numCPU
}

func (er *ebitenRender) startDisplay() {
	if err := ebiten.Run(er.update, er.pixelBuffer.Width(), er.pixelBuffer.Height(), er.screenScale, "Zubes Ray Tracer"); err != nil {
		log.Fatal(err)
	}
}

func (er *ebitenRender) doRender(rayTracer raytrace.RayTracer, renderConfig *raytrace.RenderConfig, scene raytrace.Scene) {

	var render = raytrace.NewRenderer(renderConfig)
	render.Render(rayTracer, er.pixelBuffer, scene, renderConfig)
}
