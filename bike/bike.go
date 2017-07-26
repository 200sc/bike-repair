package bike

import (
	"image"
	"image/color"
	"image/draw"

	"github.com/200sc/go-dist/colorrange"
	"github.com/oakmound/oak/render"
)

var (
	// todo: fix this color range issue, where 255-255 will result in alpha 0
	tireColor   = colorrange.NewLinear(color.RGBA{0, 0, 0, 254}, color.RGBA{255, 255, 255, 254})
	inRimColor  = colorrange.NewLinear(color.RGBA{0, 0, 0, 254}, color.RGBA{255, 255, 255, 254})
	outRimColor = colorrange.NewLinear(color.RGBA{0, 0, 0, 254}, color.RGBA{255, 255, 255, 254})
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

func (b *Bike) Draw(buff draw.Image) {
	b.DrawOffset(buff, 0, 0)
}

func (b *Bike) DrawOffset(buff draw.Image, xOff, yOff float64) {
	rgba := b.GetRGBA()
	// Draw frame
	frgba := b.Frame.buildRGBA()
	draw.Draw(rgba,
		rgba.Bounds(),
		frgba,
		image.Point{int(b.Frame.X()), int(b.Frame.Y())},
		draw.Over,
	)
	// Draw wheels
	fwrgba := b.frontWheel.buildRGBA()
	draw.Draw(rgba,
		rgba.Bounds(),
		fwrgba,
		image.Point{int(b.frontWheel.X()), int(b.frontWheel.Y())},
		draw.Over)
	bwrgba := b.backWheel.buildRGBA()
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
