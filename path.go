package rpg

import (
	"log"
	"time"

	astar "github.com/beefsack/go-astar"
	"github.com/faiface/pixel"
)

func (e *Entity) pathcalc(target pixel.Vec) {
	var (
		maxcost = 1000.00
	)
	if !e.calculated.IsZero() && time.Since(e.calculated) < time.Millisecond {

		return
	}
	e.calculated = time.Now()

	// get tiles, give world
	tile := e.w.Tile(e.Rect.Center())
	tile.W = e.w
	targett := e.w.Tile(target)
	targett.W = e.w

	// check
	if tile.Type == O_NONE {
		// bad spawn, respawn
		e.P.Health = 0
		return
	}
	if targett.Type == O_NONE {
		// player must be flying
		e.calculated = time.Now().Add(3 * time.Second)
		return
	}

	est := tile.PathEstimatedCost(targett)
	if est < 64 {
		//log.Println("direct to char", e, est)
		e.paths = []pixel.Vec{e.w.Char.Rect.Center(), e.w.Char.Rect.Center(), e.w.Char.Rect.Center()}
		return
	}

	if tile.PathEstimatedCost(targett) > 400 {
		// too far
		//log.Println("path too expensive, trying in 3 seconds")
		e.calculated = time.Now().Add(1 * time.Second)
		return
	}

	// calculate path
	path, distance, found := astar.Path(tile, targett)
	if found {
		if distance > maxcost { // cost path
			e.calculated = time.Now().Add(3 * time.Second)
			log.Println("too far")
			e.paths = nil
			return
		}
		//log.Println("distance:", distance)
		var paths []pixel.Vec

		for _, p := range path {

			//log.Println(p)
			center := p.(Object).Loc.Add(DefaultSpriteRectangle.Center())
			paths = append(paths, center)
		}

		e.paths = paths

		return
	}
	//log.Println(e.Name, "no path found, distance:", distance)
}

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
