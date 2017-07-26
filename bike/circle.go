package bike

import (
	"image"
	"image/color"
	"math"

	"github.com/oakmound/oak/alg"
	"github.com/oakmound/oak/physics"
)

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
