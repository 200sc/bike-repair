package bike

import (
	"image"
	"image/color"

	"github.com/oakmound/oak/render"
)

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

func (wh *Wheel) Radius() float64 {
	return wh.Rim.radius + wh.Tire.thickness + wh.Rim.thickness()
}

func (wh *Wheel) buildRGBA() *image.RGBA {
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
