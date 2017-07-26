package main

import (
	"image"
	"image/color"
	"image/draw"
	"math"
	"path/filepath"

	"github.com/200sc/go-dist/colorrange"
	"github.com/200sc/go-dist/floatrange"
	"github.com/200sc/go-dist/intrange"

	"github.com/oakmound/oak/alg"
	"github.com/oakmound/oak/alg/intgeom"
	"github.com/oakmound/oak/event"
	"github.com/oakmound/oak/physics"

	"github.com/oakmound/oak"
	"github.com/oakmound/oak/render"
)

var (
	daySprite   *render.Sprite
	daySpeed    = -.08
	bikesPerDay = 1
	// todo: fix this color range issue, where 255-255 will result in alpha 0
	tireColor   = colorrange.NewLinear(color.RGBA{0, 0, 0, 254}, color.RGBA{255, 255, 255, 254})
	inRimColor  = colorrange.NewLinear(color.RGBA{0, 0, 0, 254}, color.RGBA{255, 255, 255, 254})
	outRimColor = colorrange.NewLinear(color.RGBA{0, 0, 0, 254}, color.RGBA{255, 255, 255, 254})
)

func main() {
	oak.AddScene("day", sceneStart, sceneLoop, sceneEnd)
	oak.Init("day")
}

func sceneStart(prev string, input interface{}) {
	bkg := render.LoadSprite(filepath.Join("raw", "room.png"))
	render.Draw(bkg, 1)
	daySprite = render.NewSprite(154, 112, dayGradient())
	render.Draw(daySprite, 0)
	// Oh man this initialization is gonna suck
	bk := &Bike{}
	bk.Sprite = render.NewEmptySprite(0, 0, 640, 480)
	bk.frontWheel = Wheel{}
	bk.frontWheel.Sprite = render.NewEmptySprite(0, 0, 1, 1)
	bk.frontWheel.needsRedraw = true
	bk.frontWheel.Rim = Rim{
		outerThickness: 3.0,
		innerThickness: 2.0,
		outerColor:     outRimColor.Poll(),
		innerColor:     inRimColor.Poll(),
		radius:         50,
	}
	bk.frontWheel.Tire = Tire{
		thickness: 5.0,
		color:     tireColor.Poll(),
	}
	bk.backWheel = bk.frontWheel
	bk.backWheel.Sprite = render.NewEmptySprite(0, 0, 1, 1)
	bk.backWheel.Rim.innerColor = inRimColor.Poll()
	// This negative thing is weird probably need to blame it on shiny
	// todo: wheels shouldn't control their positioning, the frame should
	bk.backWheel.SetPos(-150, 0)
	bk.SetPos(300, 300)
	event.GlobalBind(func(int, interface{}) int {
		if oak.IsDown("W") {
			bk.ShiftY(-2)
		}
		if oak.IsDown("S") {
			bk.ShiftY(2)
		}
		if oak.IsDown("A") {
			bk.ShiftX(-2)
		}
		if oak.IsDown("D") {
			bk.ShiftX(2)
		}
		return 0
	}, "EnterFrame")
	render.Draw(bk, 2)
}

func sceneLoop() bool {
	daySprite.ShiftY(daySpeed)
	return daySprite.Y() > -1*((255*6)-212)
}

func sceneEnd() (string, *oak.SceneResult) {
	return "day", nil
}

type Bike struct {
	*render.Sprite
	Frame
	Gears
	Pedals
	Handlebars
	Seat
	frontWheel, backWheel Wheel
	accessories           []Addon
}

func (b *Bike) Draw(buff draw.Image) {
	b.DrawOffset(buff, 0, 0)
}

func (b *Bike) DrawOffset(buff draw.Image, xOff, yOff float64) {
	rgba := b.GetRGBA()
	// Draw frame
	// Draw wheels
	fwrgba := b.frontWheel.BuildRGBA()
	draw.Draw(rgba,
		rgba.Bounds(),
		fwrgba,
		image.Point{int(b.frontWheel.X()), int(b.frontWheel.Y())},
		draw.Over)
	bwrgba := b.backWheel.BuildRGBA()
	draw.Draw(rgba,
		rgba.Bounds(),
		bwrgba,
		image.Point{int(b.backWheel.X()), int(b.backWheel.Y())},
		draw.Over)
	// Draw gears
	// Draw seat
	// Draw pedals
	// Draw handlebars
	// Draw other necessary things?
	// Draw addons
	b.Sprite.DrawOffset(buff, xOff, yOff)
}

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
}

type Gears struct {
	*render.Sprite
}

type Pedals struct {
	*render.Sprite
}

type Seat struct {
	*render.Sprite
}

type Handlebars struct {
	*render.Sprite
}

type Addon struct {
	*render.Sprite
}

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

type Tire struct {
	flat      bool
	thickness float64
	color     color.Color
	// todo: style / tread
}

type Wheel struct {
	*render.Sprite
	Rim
	Tire
	needsRedraw bool
}

func (wh *Wheel) BuildRGBA() *image.RGBA {
	// if w's rgba does not need updating, just return it
	// todo: manipulate this needsRedraw field
	// todo: generalize this needsRedraw thing to everything on the bike
	if wh.needsRedraw {
		w := int((wh.Rim.radius * 2) + (wh.Tire.thickness * 2) + (wh.Rim.thickness() * 2))
		h := w
		rgba := image.NewRGBA(image.Rect(0, 0, w, h))
		// Draw a bunch of circles
		// This is what oak should be able to do for us
		radius := wh.Rim.radius + wh.Rim.innerThickness
		inset := wh.Rim.outerThickness + wh.Tire.thickness
		drawCircle(rgba, wh.Rim.innerColor, radius, wh.Rim.innerThickness, inset, inset)
		radius += wh.Rim.outerThickness
		inset -= wh.Rim.outerThickness
		drawCircle(rgba, wh.Rim.outerColor, radius, wh.Rim.outerThickness, inset, inset)
		radius += wh.Tire.thickness
		inset -= wh.Tire.thickness
		drawCircle(rgba, wh.Tire.color, radius, wh.Tire.thickness, inset, inset)
		wh.SetRGBA(rgba)
	}
	return wh.GetRGBA()
}

// drawCircle draws from radius inward by thickness. A thickness of zero draws
// nothing
func drawCircle(rgba *image.RGBA, c color.Color, radius, thickness float64, offsets ...float64) {
	offX := 0.0
	offY := 0.0
	if len(offsets) > 0 {
		offX = offsets[0]
		if len(offsets) > 1 {
			offY = offsets[1]
		}
	}
	rVec := physics.NewVector(radius+offX, radius+offY)
	delta := physics.AngleVector(0)
	circum := 2 * radius * math.Pi
	rotation := 360 / circum
	for j := 0.0; j < circum; j++ {
		// Todo: determine angle increment needed
		// uh duh its the circumference
		delta.Rotate(rotation)
		// We add rVec to move from -1->1 to 0->2 in terms of radius scale
		start := delta.Copy().Scale(radius).Add(rVec)
		//fmt.Println("Start", start.X(), start.Y(), "Delta", delta.X(), delta.Y())
		for i := 0.0; i < thickness; i++ {
			// this pixel is radius minus the delta, to move inward
			cur := start.Add(delta.Copy().Scale(-1))
			//fmt.Println("Setting pixel", cur.X(), cur.Y())
			rgba.Set(alg.RoundF64(cur.X()), alg.RoundF64(cur.Y()), c)
		}
	}
}

func dayGradient() *image.RGBA {
	h := 255 * 6
	w := 52
	timeBkg := image.NewRGBA(image.Rect(0, 0, w, h))
	r := 0
	g := 255
	b := 255
	dr := 0
	dg := -1
	db := 0
	for y := 0; y < h; y++ {
		c := color.RGBA{uint8(r), uint8(g), uint8(b), 255}
		for x := 0; x < w; x++ {
			timeBkg.Set(x, y, c)
		}
		if r == 0 && g == 0 && b == 255 {
			dr = 1
			dg = 0
			db = 0
		} else if r == 255 && g == 0 && b == 255 {
			dr = 0
			dg = 0
			db = -1
		} else if r == 255 && g == 0 && b == 0 {
			dr = 0
			dg = 1
			db = 0
		} else if r == 255 && g == 255 && b == 0 {
			dr = 0
			dg = 0
			db = 1
		} else if r == 255 && g == 255 && b == 255 {
			dr = -1
			db = 0
			dg = 0
		}
		r += dr
		g += dg
		b += db
	}
	return timeBkg
}
