package common

import (
	astar "github.com/beefsack/go-astar"
	"github.com/faiface/pixel"
)

func (o Object) PathNeighbors() []astar.Pather {
	neighbors := []astar.Pather{}
	of := 32.0
	//of = 24.0

	for _, offset := range [][]float64{
		{-of, 0},
		{of, 0},
		{0, -of},
		{0, of},
		//{of, -of},
		//{-of, -of},
		//{of, of},
		//{-of, of},
	} {
		n := o.W.Tile(pixel.V(o.Rect.Center().X+offset[0], o.Rect.Center().Y+offset[1]))
		if n.Type == O_TILE {
			neighbors = append(neighbors, n)
		}
	}
	return neighbors
}

func (o Object) PathEstimatedCost(to astar.Pather) float64 {
	toT := to.(Object)
	absX := toT.Rect.Center().X - o.Rect.Center().X
	if absX < 0 {
		absX = -absX
	}
	absY := toT.Rect.Center().Y - o.Rect.Center().Y
	if absY < 0 {
		absY = -absY
	}
	return float64(absX + absY)

}

func (o Object) PathNeighborCost(to astar.Pather) float64 {
	toT := to.(Object)
	cost := tileCosts[toT.Type]
	return cost
}

// KindCosts map tile kinds to movement costs.

var tileCosts = map[ObjectType]float64{
	O_NONE:      30.00,
	O_BLOCK:     30.00,
	O_INVISIBLE: 3.00,
	O_SPECIAL:   0.00,
	O_TILE:      1.00,
	O_WIN:       0.00,
	O_DYNAMIC:   3.00,
}
