package main

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
	color colorVector
}

type renderer struct {
	jobs    chan pixelJob
	results chan pixelResult
	config  *renderConfig
}

func newRenderer(config *renderConfig) *renderer {
	var jobs = make(chan pixelJob, 1000)
	var results = make(chan pixelResult, 1000)
	return &renderer{jobs, results, config}
}

func (r *renderer) Render(pixelBuffer *pixelBuffer, scene *scene, renderConfig *renderConfig) {

	// fire progress callback that we've started

	camera := scene.GetCamera(pixelBuffer.Width(), pixelBuffer.Height())
	world := scene.GetWorld()
	lightHitable := scene.GetLightHitable()
	backgroundFunc := scene.GetBackgroundFunc()

	if renderConfig.IsTwoPhase() {
		r.renderMultithread(
			pixelBuffer,
			camera,
			world,
			lightHitable,
			NewRenderConfig(renderConfig.numThreads, 5, 1, false),
			backgroundFunc)
	}

	r.renderMultithread(
		pixelBuffer,
		camera,
		world,
		lightHitable,
		renderConfig,
		backgroundFunc)
}

func (r *renderer) enqueuePixels(width int, height int) {
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			job := pixelJob{x, y}
			r.jobs <- job
		}
	}
	close(r.jobs)
}

func (r *renderer) calcResults(pixelBuffer *pixelBuffer, done chan bool) {
	for result := range r.results {
		fmt.Printf("SetPixelColor (%v, %v) -> (%v, %v, %v)\n", result.job.x, result.job.y, result.color.r, result.color.g, result.color.b)
		pixelBuffer.SetPixelColor(result.job.x, result.job.y, result.color)
	}
	done <- true
}

func (r *renderer) pixelWorker(rayTracer *simpleTracer, wg *sync.WaitGroup) {
	for job := range r.jobs {
		color := rayTracer.GetPixelColor(job.x, job.y)
		result := pixelResult{job, color}
		r.results <- result
	}
	wg.Done()
}

func (r *renderer) createWorkerPool(rayTracer *simpleTracer, numWorkers int) {
	var wg sync.WaitGroup
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go r.pixelWorker(rayTracer, &wg)
	}
	wg.Wait()
	close(r.results)
}

func (r *renderer) renderMultithread(
	pixelBuffer *pixelBuffer,
	camera *camera,
	world *hitable,
	lightHitable *hitable,
	renderConfig *renderConfig,
	backgroundFunc backgroundFunc) {

	startTime := time.Now()

	rayTracer := newSimpleTracer()

	go r.enqueuePixels(pixelBuffer.Width(), pixelBuffer.Height())
	done := make(chan bool)
	go r.calcResults(pixelBuffer, done)
	r.createWorkerPool(rayTracer, renderConfig.numThreads)

	<-done

	endTime := time.Now()
	diff := endTime.Sub(startTime)
	fmt.Println("total time ", diff.Seconds(), "seconds")
}
