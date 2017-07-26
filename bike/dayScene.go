package bike

import (
	"fmt"
	"image"
	"image/color"
	"path/filepath"

	"github.com/oakmound/oak"
	"github.com/oakmound/oak/event"
	"github.com/oakmound/oak/render"
)

var (
	daySprite   *render.Sprite
	daySpeed    = -.08
	bikesPerDay = 1
)

func DayStart(prev string, input interface{}) {
	bkg := render.LoadSprite(filepath.Join("raw", "room.png"))
	render.Draw(bkg, 1)
	daySprite = render.NewSprite(154, 112, dayGradient())
	render.Draw(daySprite, 0)
	// Oh man this initialization is gonna suck
	// Oh man does this initialization suckkkkk
	bk := &Bike{}
	bk.Frame = NewFrame()
	bk.Frame.needsRedraw = true
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
	backWheelPos := bk.Frame.nodes[bk.Frame.backWheelIndex] //.Add(
	//intgeom.NewPoint(int(bk.backWheel.Radius()), int(bk.backWheel.Radius())))
	bk.backWheel.SetPos(-float64(backWheelPos.X), -float64(backWheelPos.Y))

	frontWheelPos := bk.Frame.nodes[bk.Frame.frontWheelIndex] //.Add(
	//intgeom.NewPoint(int(bk.frontWheel.Radius()), int(bk.frontWheel.Radius())))
	bk.frontWheel.SetPos(-float64(frontWheelPos.X), -float64(frontWheelPos.Y))

	fmt.Println(frontWheelPos, backWheelPos)

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

func DayLoop() bool {
	daySprite.ShiftY(daySpeed)
	return daySprite.Y() > -1*((255*6)-212)
}

func DayEnd() (string, *oak.SceneResult) {
	return "day", nil
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
