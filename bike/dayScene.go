package bike

import (
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
	bk := NewBike()

	bk.SetPos(200, 200)
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
