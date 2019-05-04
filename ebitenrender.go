package main

import (
	"fmt"
	"log"
	"math/rand"
	"path/filepath"
	"time"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/markzuber/zgotrace/raytrace"
)

type ebitenRender struct {
	offscreen      *ebiten.Image
	pixelBuffer    *raytrace.PixelBuffer
	screenScale    float64
	outputFilePath string
	showFps        bool
}

func newEbitenRender(imageWidth int, imageHeight int, screenScale float64, outputFilePath string, showFps bool) *ebitenRender {
	rand.Seed(time.Now().UnixNano())
	offscreen, _ := ebiten.NewImage(imageWidth, imageHeight, ebiten.FilterDefault)
	var pixelBuffer = raytrace.NewPixelBuffer(imageWidth, imageHeight)

	return &ebitenRender{offscreen, pixelBuffer, screenScale, outputFilePath, showFps}
}

func (er *ebitenRender) update(screen *ebiten.Image) error {
	ebImg, _ := ebiten.NewImageFromImage(er.pixelBuffer.GetImage(), ebiten.FilterDefault)

	if ebiten.IsDrawingSkipped() {
		return nil
	}

	screen.DrawImage(ebImg, nil)

	if er.showFps {
		ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f", ebiten.CurrentTPS()))
	}

	return nil
}

func (er *ebitenRender) StartDisplay(windowTitle string) {
	ebiten.SetRunnableInBackground(true)
	if err := ebiten.Run(er.update, er.pixelBuffer.Width(), er.pixelBuffer.Height(), er.screenScale, windowTitle); err != nil {
		log.Fatal(err)
	}
}

func (er *ebitenRender) DoRender(rayTracer raytrace.RayTracer, renderConfig *raytrace.RenderConfig, scene raytrace.Scene, imageOutputFilename string) {
	var render = raytrace.NewRenderer(renderConfig)

	fmt.Printf("Rendering for %v started at: %v\n", imageOutputFilename, time.Now())
	render.Render(rayTracer, er.pixelBuffer, scene, renderConfig)
	outFilePath := filepath.Join(er.outputFilePath, imageOutputFilename)
	fmt.Printf("Saving image to %v ", outFilePath)
	er.pixelBuffer.SavePng(outFilePath)
	fmt.Println("...done")
}

type SceneInfo struct {
	scene               raytrace.Scene
	imageOutputFilename string
}

func (er *ebitenRender) DoRenderMulti(rayTracer raytrace.RayTracer, renderConfig *raytrace.RenderConfig, sceneInfos []SceneInfo) {
	for _, si := range sceneInfos {
		er.DoRender(rayTracer, renderConfig, si.scene, si.imageOutputFilename)
	}
}
