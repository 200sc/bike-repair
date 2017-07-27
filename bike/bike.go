package bike

import (
	"image"
	"image/draw"
	"math"

	"github.com/davecgh/go-spew/spew"
	"github.com/oakmound/oak/alg/intgeom"
	"github.com/oakmound/oak/render"
)

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

func NewBike() *Bike {
	bk := &Bike{}
	bk.Frame = NewFrame()
	bk.Frame.needsRedraw = true
	bk.Sprite = render.NewEmptySprite(0, 0, 640, 480)
	bk.frontWheel = NewWheel()
	bk.backWheel = NewWheel()
	// This negative thing is weird probably need to blame it on shiny
	// todo: wheels shouldn't control their positioning, the frame should
	backWheelPos := bk.Frame.nodes[bk.Frame.backWheelIndex].Add(
		intgeom.NewPoint(-int(bk.backWheel.Radius()), -int(bk.backWheel.Radius())))
	bk.backWheel.SetPos(float64(backWheelPos.X), float64(backWheelPos.Y))

	frontWheelPos := bk.Frame.nodes[bk.Frame.frontWheelIndex].Add(
		intgeom.NewPoint(-int(bk.frontWheel.Radius()), -int(bk.frontWheel.Radius())))
	bk.frontWheel.SetPos(float64(frontWheelPos.X), float64(frontWheelPos.Y))

	lowest := intgeom.NewPoint(math.MaxInt32, math.MaxInt32)
	highest := intgeom.NewPoint(-math.MaxInt32, -math.MaxInt32)

	for _, n := range bk.Frame.nodes {
		lowest = MinComponents(lowest, n)
		highest = MaxComponents(highest, n)
	}
	lowest = MinComponents(lowest, backWheelPos)
	highest = MaxComponents(highest, backWheelPos)
	lowest = MinComponents(lowest, frontWheelPos)
	highest = MaxComponents(highest, frontWheelPos)

	if lowest.X < 0 {
		flowestX := float64(lowest.X)
		bk.Frame.ShiftX(-flowestX)
		bk.backWheel.ShiftX(-flowestX)
		bk.frontWheel.ShiftX(-flowestX)
		highest.X -= lowest.X
	}

	if lowest.Y < 0 {
		flowestY := float64(lowest.Y)
		bk.Frame.ShiftY(-flowestY)
		bk.backWheel.ShiftY(-flowestY)
		bk.frontWheel.ShiftY(-flowestY)
		highest.Y -= lowest.Y
	}

	// Todo: add widths to highest points
	bk.Sprite = render.NewEmptySprite(0, 0, highest.X+100, highest.Y+100)
	spew.Config.MaxDepth = 20
	spew.Dump(bk.frontWheel, bk.backWheel)
	return bk
}

func (b *Bike) Draw(buff draw.Image) {
	b.DrawOffset(buff, 0, 0)
}

func (b *Bike) DrawOffset(buff draw.Image, xOff, yOff float64) {
	// Draw frame
	b.DrawComponentTo(b.Frame)
	// Draw wheels
	b.DrawComponentTo(b.frontWheel)
	b.DrawComponentTo(b.backWheel)
	// Draw gears
	// Draw seat
	// Draw pedals
	// Draw handlebars
	// Draw other necessary things?
	// Draw addons
	b.Sprite.DrawOffset(buff, xOff, yOff)
}

func (b *Bike) DrawComponentTo(c Component) {
	rgba := c.buildRGBA()
	render.ShinyDraw(b.GetRGBA(), rgba, int(c.X()), int(c.Y()))
	// draw.Draw(b.GetRGBA(),
	// 	rgba.Bounds(),
	// 	rgba,
	// 	image.Point{int(c.X()), int(c.Y())},
	// 	draw.Over,
	// )
}

type Component interface {
	X() float64
	Y() float64
	buildRGBA() *image.RGBA
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
