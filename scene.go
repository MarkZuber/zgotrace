package main

type backgroundFunc func() *colorVector

type scene struct {
	camera         *camera
	world          *hitable
	lighthitable   *hitable
	backgroundfunc backgroundFunc
}

func (s *scene) GetCamera(width int, height int) *camera {
	return s.camera
}

func (s *scene) GetWorld() *hitable {
	return s.world
}

func (s *scene) GetLightHitable() *hitable {
	return s.lighthitable
}

func (s *scene) GetBackgroundFunc() backgroundFunc {
	return s.backgroundfunc
}
