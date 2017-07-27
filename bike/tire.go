package bike

import (
	"image/color"

	"github.com/200sc/go-dist/colorrange"
	"github.com/200sc/go-dist/floatrange"
)

var (
	tireColor     = colorrange.NewLinear(color.RGBA{0, 0, 0, 254}, color.RGBA{255, 255, 255, 254})
	tireThickness = floatrange.NewLinear(3, 20)
)

type Tire struct {
	flat      bool
	thickness float64
	color     color.Color
	// todo: style / tread
}

func NewTire() Tire {
	return Tire{
		thickness: 5.0,
		color:     tireColor.Poll(),
	}
}
