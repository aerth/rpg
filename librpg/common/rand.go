package common

import (
	"math"
	"math/rand"

	"github.com/faiface/pixel"
)

func RandomColor() pixel.RGBA {

	r := rand.Float64()
	g := rand.Float64()
	b := rand.Float64()
	len := math.Sqrt(r*r + g*g + b*b)
	//if len == 0 {
	//	goto again
	//}
	return pixel.RGB(r/len, g/len, b/len)

}
