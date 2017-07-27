package bike

import (
	"image"

	"github.com/oakmound/oak/render"
)

type Wheel struct {
	*render.Sprite
	Rim
	Tire
	needsRedraw bool
}

func NewWheel() Wheel {
	w := Wheel{}
	w.Sprite = render.NewEmptySprite(0, 0, 1, 1)
	w.needsRedraw = true
	w.Rim = NewRim()
	w.Tire = NewTire()
	return w
}

func NewWheelPair() (Wheel, Wheel) {
	w := NewWheel()
	w2 := w
	w2.Sprite = render.NewEmptySprite(0, 0, 1, 1)
	return w, w2
}

func (wh Wheel) Radius() float64 {
	return wh.Rim.radius + wh.Tire.thickness + wh.Rim.thickness()
}

func (wh Wheel) buildRGBA() *image.RGBA {
	// if w's rgba does not need updating, just return it
	// todo: manipulate this needsRedraw field
	// todo: generalize this needsRedraw thing to everything on the bike
	if wh.needsRedraw {
		w := int(wh.Radius() * 2)
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
