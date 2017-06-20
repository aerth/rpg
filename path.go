package rpg

import (
	"log"
	"time"

	astar "github.com/beefsack/go-astar"
	"github.com/faiface/pixel"
)

func (e *Entity) pathcalc(target pixel.Vec) {
	var (
		maxcost = 200.00 // 100 regular tiles away
	)
	if len(e.paths) == 0 && time.Since(e.calculated) < time.Second*3 {
		return
	}
	e.calculated = time.Now()
	/*	t1 := time.Now()
		defer func(t time.Time) {
			log.Println("TILE path calc took:", time.Since(t))
		}(t1)
	*/tile := e.w.Tile(e.Rect.Center())
	targett := e.w.Tile(target)
	if tile.Type == O_NONE || targett.Type == O_NONE {
		if tile.Type != O_TILE { // we spawned on bad tile
			e.P.Health = 0
			e.P.IsDead = true
			log.Println("killing bad entity")
			return
		}
	}
	/*	t2 := time.Now()
		defer func(t time.Time) {
			log.Println("PATH path calc took:", time.Since(t))
		}(t2)
	*/path, distance, found := astar.Path(tile, targett)
	if found {
		if distance > maxcost { // cost path
			e.calculated = time.Now().Add(-time.Minute)
			log.Println("too far")
			e.paths = nil
			return
		}
		//log.Println("distance:", distance)
		var paths []pixel.Vec

		for _, p := range path {

			//log.Println(p)
			center := p.(Object).Rect.Center()
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
		n := o.w.Object(pixel.V(o.Rect.Center().X+offset[0], o.Rect.Center().Y+offset[1]))
		if n.Type != O_NONE {
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

type ObjectType int

const (
	O_NONE ObjectType = iota
	O_TILE
	O_BLOCK
	O_INVISIBLE
	O_SPECIAL
	O_WIN
)

// KindCosts map tile kinds to movement costs.

var tileCosts = map[ObjectType]float64{
	O_NONE:      2.00,
	O_BLOCK:     30.00,
	O_INVISIBLE: 3.00,
	O_SPECIAL:   0.00,
	O_TILE:      1.00,
	O_WIN:       0.00,
}
