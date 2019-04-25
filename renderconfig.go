package main

type renderConfig struct {
	numThreads      int
	maxDepth        int
	samplesPerPixel int
	isTwoPhase      bool
}

func NewRenderConfig(numThreads int, maxDepth int, samplesPerPixel int, isTwoPhase bool) *renderConfig {
	return &renderConfig{numThreads, maxDepth, samplesPerPixel, isTwoPhase}
}

func (c *renderConfig) IsTwoPhase() bool {
	return c.isTwoPhase
}
