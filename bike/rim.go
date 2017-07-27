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
	rimRadius       = floatrange.NewLinear(30, 100)
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
		outerThickness: outRimThickness.Poll(),
		innerThickness: inRimThickness.Poll(),
		outerColor:     outRimColor.Poll(),
		innerColor:     inRimColor.Poll(),
		radius:         rimRadius.Poll(),
	}
}
