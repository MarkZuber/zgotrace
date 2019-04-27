package raytrace

import (
	"fmt"
	"math"
	"math/rand"

	"github.com/davecgh/go-spew/spew"
)

type MonteCarloTracer struct {
	imageWidth     int
	imageHeight    int
	camera         Camera
	world          Hitable
	lightHitable   Hitable
	renderConfig   *RenderConfig
	backgroundFunc BackgroundFunc
}

func NewMonteCarloTracer() RayTracer {
	return &MonteCarloTracer{}
}

func (t *MonteCarloTracer) Configure(
	imageWidth int,
	imageHeight int,
	camera Camera,
	world Hitable,
	lightHitable Hitable,
	renderConfig *RenderConfig,
	backgroundFunc BackgroundFunc) {
	t.imageWidth = imageWidth
	t.imageHeight = imageHeight
	t.camera = camera
	t.world = world
	t.lightHitable = lightHitable
	t.renderConfig = renderConfig
	t.backgroundFunc = backgroundFunc
}

func (t *MonteCarloTracer) GetPixelColor(x int, y int) ColorVector {

	fmt.Printf("IN GetPixelColor(%v, %v)\n", x, y)
	color := NewColorVector(0, 0, 0)
	xfloat := float64(x)
	yfloat := float64(y)

	if t.renderConfig.samplesPerPixel > 1 {
		for sample := 0; sample < t.renderConfig.samplesPerPixel; sample++ {
			u := (xfloat + rand.Float64()) / float64(t.imageWidth)
			v := (yfloat + rand.Float64()) / float64(t.imageHeight)
			r := t.camera.GetRay(u, v)

			color = color.Add(t.getRayColor(r, t.world, 0))
		}

		color = color.DivScalar(float64(t.renderConfig.samplesPerPixel))

	} else {
		color = t.getRayColor(t.camera.GetRay(xfloat/float64(t.imageWidth), yfloat/float64(t.imageHeight)), t.world, 0)
	}

	color = color.ApplyGamma2()
	return color
}

func (t *MonteCarloTracer) getRayColor(ray *Ray, world Hitable, depth int) ColorVector {
	spew.Dump(ray)

	hr := world.Hit(ray, 0.001, math.MaxFloat64)
	if hr != nil {
		fmt.Printf("WE GOT A HIT!\n")
		emitted := hr.Material().Emitted(ray, hr, hr.UvCoords(), hr.P())

		if depth < t.renderConfig.maxDepth {
			scatterResult := hr.Material().Scatter(ray, hr)
			if scatterResult.IsScattered() {
				if scatterResult.IsSpecular() {
					return scatterResult.Attenuation().Mul(t.getRayColor(scatterResult.SpecularRay(), world, depth+1))
				}
				p0 := NewHitablePdf(t.lightHitable, hr.P())
				p := NewMixturePdf(p0, scatterResult.Pdf())
				scattered := NewRay(hr.P(), p.Generate())
				pdfValue := p.GetValue(scattered.Direction())

				scatteringPdf := hr.Material().ScatteringPdf(ray, hr, scattered)
				if scatteringPdf < 0.01 {
					scatteringPdf = 0.01
				}

				depthRayColor := t.getRayColor(scattered, world, depth+1)
				recurseColor := ((scatterResult.Attenuation().Mul(depthRayColor).MulScalar(scatteringPdf)).DivScalar(pdfValue))
				return emitted.Add(recurseColor)
			}
		}

		return emitted
	}

	return t.backgroundFunc() // todo: have this take Ray as a parameter?
}
