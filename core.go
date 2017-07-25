package main

import (
	"image"
	"image/color"
	"path/filepath"

	"github.com/oakmound/oak"
	"github.com/oakmound/oak/physics"
	"github.com/oakmound/oak/render"
)

var (
	time     physics.Vector
	daySpeed float64 = -3
)

func main() {
	oak.AddScene("day", sceneStart, sceneLoop, sceneEnd)
	oak.Init("day")
}

func sceneStart(prev string, input interface{}) {
	bkg := render.LoadSprite(filepath.Join("raw", "room.png"))
	render.Draw(bkg, 1)

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
	s := render.NewSprite(154, 112, timeBkg)
	time = physics.NewVector(0, 0)
	time = time.Attach(s)
	render.Draw(s, 0)
}

func sceneLoop() bool {
	time.ShiftY(daySpeed)
	return time.Y() > -1*((255*6)-212)
}

func sceneEnd() (string, *oak.SceneResult) {
	return "day", nil
}
