package bike

import "github.com/oakmound/oak/alg/intgeom"

func MinComponents(a, b intgeom.Point) intgeom.Point {
	if b.X < a.X {
		a.X = b.X
	}
	if b.Y < a.Y {
		a.Y = b.Y
	}
	return a
}

func MaxComponents(a, b intgeom.Point) intgeom.Point {
	if b.X > a.X {
		a.X = b.X
	}
	if b.Y > a.Y {
		a.Y = b.Y
	}
	return a
}
