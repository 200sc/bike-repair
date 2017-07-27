package bike

import (
	"image"
	"image/color"

	"github.com/200sc/go-dist/colorrange"
	"github.com/200sc/go-dist/floatrange"
	"github.com/200sc/go-dist/intrange"
	"github.com/oakmound/oak/alg"
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
	gearsIndex  int
	w, h        int
	needsRedraw bool
}

var (
	frameNodes             = intrange.NewLinear(5, 10)
	frameWidth             = intrange.NewLinear(150, 500)
	frameHeight            = intrange.NewLinear(150, 350)
	frameThickness         = floatrange.NewLinear(2, 14)
	frameColor             = colorrange.NewLinear(color.RGBA{0, 0, 0, 254}, color.RGBA{255, 255, 255, 254})
	frameExcessConnections = intrange.NewLinear(1, 4)
)

func NewFrame() Frame {
	f := Frame{}
	f.thickness = frameThickness.Poll()
	f.color = frameColor.Poll()
	f.w, f.h = frameWidth.Poll(), frameHeight.Poll()
	f.Sprite = render.NewEmptySprite(0, 0, f.w, f.h)
	xPositions := intrange.NewLinear(0, f.w-1)
	yPositions := intrange.NewLinear(0, f.h-1)
	f.nodes = make([]intgeom.Point, frameNodes.Poll())
	for i := range f.nodes {
		// This will produce atrocious bicycles
		f.nodes[i] = intgeom.NewPoint(xPositions.Poll(), yPositions.Poll())
	}
	// The bottom left most point is backWheel
	// The bottom right most point is frontWheel
	// The top left most point is seat
	// The top right most point is handlebars
	// gears is any other point
	botLeft := intgeom.NewPoint(f.w, -1)
	botRight := intgeom.NewPoint(-1, -1)
	topLeft := intgeom.NewPoint(f.w, f.h)
	topRight := intgeom.NewPoint(-1, f.h)
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
	f.connections = ConnectGraph(f.nodes)
	// Add a few more additional connections
	f.connections = AddRandomConnections(frameExcessConnections.Poll(), f.connections)
	return f
}

func (f Frame) buildRGBA() *image.RGBA {
	if f.needsRedraw {
		rgba := image.NewRGBA(image.Rect(0, 0, f.w, f.h))
		for i, list := range f.connections {
			n1 := f.nodes[i]
			for _, c := range list {
				n2 := f.nodes[c]
				render.DrawThickLine(rgba, n1.X, n1.Y, n2.X, n2.Y, f.color, alg.RoundF64(f.thickness))
			}
		}
		f.SetRGBA(rgba)
	}
	return f.GetRGBA()
}
