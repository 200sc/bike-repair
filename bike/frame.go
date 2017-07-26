package bike

import (
	"image/color"

	"github.com/200sc/go-dist/floatrange"
	"github.com/200sc/go-dist/intrange"
	"github.com/oakmound/oak/alg/intgeom"
	"github.com/oakmound/oak/render"
)

type Frame struct {
	*render.Sprite
	thickness float64
	color     color.Color // todo color functions that take in x,y
	// Frames are a connected graph
	nodes       []intgeom.Point
	connections [][]int
	// With at least these five nodes
	// center of wheels
	frontWheelIndex int
	backWheelIndex  int
	// bottom of handlebars
	handlebarsIndex int
	// bottom center of seat
	seatIndex int
	// center of gears
	gearsIndex int
}

var (
	frameNodes     = intrange.NewLinear(5, 10)
	frameWidth     = intrange.NewLinear(150, 350)
	frameHeight    = intrange.NewLinear(150, 500)
	frameThickness = floatrange.NewLinear(2, 20)
)

func NewFrame() Frame {
	f := Frame{}
	f.thickness = frameThickness.Poll()
	w, h := frameWidth.Poll(), frameHeight.Poll()
	xPositions := intrange.NewLinear(0, w-1)
	yPositions := intrange.NewLinear(0, h-1)
	f.nodes = make([]intgeom.Point, frameNodes.Poll())
	for i := range f.nodes {
		f.nodes[i] = intgeom.NewPoint(xPositions.Poll(), yPositions.Poll())
	}
	// The bottom left most point is backWheel
	// The bottom right most point is frontWheel
	// The top left most point is seat
	// The top right most point is handlebars
	// gears is any other point
	botLeft := intgeom.NewPoint(w, 0)
	botRight := intgeom.NewPoint(0, 0)
	topLeft := intgeom.NewPoint(w, h)
	topRight := intgeom.NewPoint(0, h)
	for i, p := range f.nodes {
		if p.X < botLeft.X && p.Y > botLeft.Y {
			botLeft = p
			f.backWheelIndex = i
		}
		if p.X > botRight.X && p.Y > botRight.Y {
			botRight = p
			f.frontWheelIndex = i
		}
		if p.X < topLeft.X && p.Y < topLeft.Y {
			topLeft = p
			f.seatIndex = i
		}
		if p.X > topRight.X && p.Y < topRight.Y {
			topRight = p
			f.handlebarsIndex = i
		}
	}
	gearPossibilities := len(f.nodes) - 4
	gears := intrange.NewLinear(0, gearPossibilities).Poll()
	for i := range f.nodes {
		if i != f.backWheelIndex && i != f.frontWheelIndex &&
			i != f.seatIndex && i != f.handlebarsIndex {
			if gears == 1 {
				f.gearsIndex = i
				break
			}
			gears--
		}
	}
	// Make connections until graph is connected
	return f
}
