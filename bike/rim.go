package bike

import (
	"image/color"

	"github.com/200sc/go-dist/colorrange"
	"github.com/200sc/go-dist/floatrange"
)

var (
	// todo: fix this color range issue, where 255-255 will result in alpha 0
	inRimColor      = colorrange.NewLinear(color.RGBA{0, 0, 0, 254}, color.RGBA{255, 255, 255, 254})
	outRimColor     = colorrange.NewLinear(color.RGBA{0, 0, 0, 254}, color.RGBA{255, 255, 255, 254})
	inRimThickness  = floatrange.NewLinear(1, 4)
	outRimThickness = floatrange.NewLinear(2, 8)
)

type Rim struct {
	outerThickness, innerThickness float64
	radius                         float64
	outerColor, innerColor         color.Color
	spokes                         int
	spokeColor                     color.Color
	// todo: style? emblem?
}

func (r Rim) thickness() float64 {
	return r.outerThickness + r.innerThickness
}

func NewRim() Rim {
	return Rim{
		outerThickness: 3.0,
		innerThickness: 2.0,
		outerColor:     outRimColor.Poll(),
		innerColor:     inRimColor.Poll(),
		radius:         50,
	}
}
