package raytrace

import (
	"fmt"
	"sync"
	"time"
)

type pixelJob struct {
	x int
	y int
}

type pixelResult struct {
	job   pixelJob
	color ColorVector
}

type Renderer struct {
	jobs    chan pixelJob
	results chan pixelResult
	config  *RenderConfig
}

func NewRenderer(config *RenderConfig) *Renderer {
	var jobs = make(chan pixelJob, 1000)
	var results = make(chan pixelResult, 1000)
	return &Renderer{jobs, results, config}
}

func (r *Renderer) Render(rayTracer RayTracer, pixelBuffer *PixelBuffer, scene Scene, renderConfig *RenderConfig) {

	// fire progress callback that we've started

	camera := scene.GetCamera(pixelBuffer.Width(), pixelBuffer.Height())
	world := scene.GetWorld()
	lightHitable := scene.GetLightHitable()
	backgroundFunc := scene.GetBackgroundFunc()

	if renderConfig.IsTwoPhase() {
		r.renderMultithread(
			rayTracer,
			pixelBuffer,
			camera,
			world,
			lightHitable,
			NewRenderConfig(renderConfig.numThreads, 5, 1, false),
			backgroundFunc)
	}

	r.renderMultithread(
		rayTracer,
		pixelBuffer,
		camera,
		world,
		lightHitable,
		renderConfig,
		backgroundFunc)
}

func (r *Renderer) enqueuePixels(width int, height int) {
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			job := pixelJob{x, y}
			r.jobs <- job
		}
	}
	close(r.jobs)
}

func (r *Renderer) calcResults(pixelBuffer *PixelBuffer, done chan bool) {
	for result := range r.results {
		// fmt.Printf("SetPixelColor (%v, %v) -> (%v, %v, %v)\n", result.job.x, result.job.y, result.color.R(), result.color.G(), result.color.B())
		pixelBuffer.SetPixelColor(result.job.x, result.job.y, result.color)
	}
	done <- true
}

func (r *Renderer) pixelWorker(rayTracer RayTracer, wg *sync.WaitGroup) {
	for job := range r.jobs {
		color := rayTracer.GetPixelColor(job.x, job.y)
		result := pixelResult{job, color}
		r.results <- result
	}
	wg.Done()
}

func (r *Renderer) createWorkerPool(rayTracer RayTracer, numWorkers int) {
	var wg sync.WaitGroup
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go r.pixelWorker(rayTracer, &wg)
	}
	wg.Wait()
	close(r.results)
}

func (r *Renderer) renderMultithread(
	rayTracer RayTracer,
	pixelBuffer *PixelBuffer,
	camera Camera,
	world Hitable,
	lightHitable Hitable,
	renderConfig *RenderConfig,
	backgroundFunc BackgroundFunc) {

	startTime := time.Now()

	rayTracer.Configure(pixelBuffer.Width(), pixelBuffer.Height(), camera, world, lightHitable, renderConfig, backgroundFunc)

	go r.enqueuePixels(pixelBuffer.Width(), pixelBuffer.Height())
	done := make(chan bool)
	go r.calcResults(pixelBuffer, done)
	r.createWorkerPool(rayTracer, renderConfig.numThreads)

	<-done

	endTime := time.Now()
	diff := endTime.Sub(startTime)
	fmt.Println("total time ", diff.Seconds(), "seconds")
}
