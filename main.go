package main

import (
	"fmt"
	"log"
	"os"

	"github.com/markzuber/zgotrace/raytrace"
	"github.com/markzuber/zgotrace/raytrace/scenes"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	verbose   = kingpin.Flag("debug", "Enable debug mode.").Short('v').Bool()
	outputDir = kingpin.Flag("outpath", "Output file path.").Short('o').ExistingDir()
	showfps   = kingpin.Flag("showfps", "Show FPS").Short('f').Bool()

	isTwoPhase   = kingpin.Flag("twophase", "Two Phase Rendering").Default("true").Bool()
	maxDepth     = kingpin.Flag("maxdepth", "Max ray recursion depth").Default("50").Int()
	numSamples   = kingpin.Flag("numsamples", "Num samples per pixel").Default("10").Int()
	imageWidth   = kingpin.Flag("imagewidth", "Width of image").Default("300").Int()
	imageHeight  = kingpin.Flag("imageHeight", "Height of image").Default("300").Int()
	displayScale = kingpin.Flag("displayscale", "Scale from image to screen").Default("2.0").Float64()
)

// Set the working directory to the root of Go package, so that its assets can be accessed.
func init() {
	dir, err := importPathToDir("github.com/markzuber/zgotrace/images")
	if err != nil {
		log.Fatalln("Unable to find Go package in your GOPATH, it's needed to load assets:", err)
	}
	err = os.Chdir(dir)
	if err != nil {
		log.Panicln("os.Chdir:", err)
	}
}

func getFileOutputDir() string {
	var fileOutputDir string

	if *outputDir == "" {
		fileOutputDir, _ = os.UserHomeDir()
	} else {
		fileOutputDir = *outputDir
	}
	return fileOutputDir
}

func main() {
	kingpin.Parse()

	fileOutputDir := getFileOutputDir()
	fmt.Printf("fileoutputdir: %v\n", fileOutputDir)

	ebRen := newEbitenRender(*imageWidth, *imageHeight, *displayScale, fileOutputDir, *showfps)

	var renderConfig = raytrace.NewRenderConfig(*maxDepth, *numSamples, *isTwoPhase)
	var rayTracer = raytrace.NewMonteCarloTracer()

	go ebRen.DoRenderMulti(rayTracer, renderConfig, []SceneInfo{
		{scenes.CreateCornellBoxScene(), "cornellbox.png"},
		{scenes.CreateCornellBoxWithSmokeScene(), "cornellbox_withsmoke.png"},
		{scenes.CreateManySpheresScene(), "manyspheres.png"},
		{scenes.CreateImageTextureScene("globetex.jpg"), "imagetex.png"}})

	ebRen.StartDisplay("Zube's RayTracer")
}
