package raytrace

import (
	"math"

	"github.com/go-gl/mathgl/mgl64"
)

type ColorVector struct {
	colors mgl64.Vec3
}

func NewColorVector(r float64, g float64, b float64) ColorVector {
	return ColorVector{mgl64.Vec3{r, g, b}}
}

func (c ColorVector) RGBA() (r, g, b, a uint32) {
	return uint32(c.R() * 0xffff), uint32(c.G() * 0xffff), uint32(c.B() * 0xffff), 0xffff
}

func (c ColorVector) R() float64 {
	return c.colors.X()
}

func (c ColorVector) G() float64 {
	return c.colors.Y()
}

func (c ColorVector) B() float64 {
	return c.colors.Z()
}

func (c ColorVector) Mul(other ColorVector) ColorVector {
	return NewColorVector(c.R()*other.R(), c.G()*other.G(), c.B()*other.B())
}

func (c ColorVector) MulScalar(other float64) ColorVector {
	return NewColorVector(c.R()*other, c.G()*other, c.B()*other)
}

func (c ColorVector) DivScalar(other float64) ColorVector {
	return NewColorVector(c.R()/other, c.G()/other, c.B()/other)
}

func (c ColorVector) Add(other ColorVector) ColorVector {
	return NewColorVector(c.R()+other.R(), c.G()+other.G(), c.B()+other.B())
}

func (c ColorVector) ApplyGamma2() ColorVector {
	return NewColorVector(math.Sqrt(c.R()), math.Sqrt(c.G()), math.Sqrt(c.B()))
}
