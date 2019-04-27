package raytrace

type RenderConfig struct {
	numThreads      int
	maxDepth        int
	samplesPerPixel int
	isTwoPhase      bool
}

func NewRenderConfig(numThreads int, maxDepth int, samplesPerPixel int, isTwoPhase bool) *RenderConfig {
	return &RenderConfig{numThreads, maxDepth, samplesPerPixel, isTwoPhase}
}

func (c *RenderConfig) IsTwoPhase() bool {
	return c.isTwoPhase
}
