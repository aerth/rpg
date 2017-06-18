package rpg

import (
	"log"

	astar "github.com/beefsack/go-astar"
	"github.com/faiface/pixel"
)

func (e *Entity) pathcalc(target pixel.Vec) {
	/*
		if e.paths != nil && len(e.paths) > 0 {
			return e.Phys.Vel
		} */
	tile := e.w.Tile(e.Rect.Center())
	targett := e.w.Tile(target)
	if tile == nil || targett == nil {
		log.Println(e.Name, "nil target, nil tile")
		return
	}

	//log.Println(e.Name, "path cost to target:", tile.PathEstimatedCost(targett))
	path, distance, found := astar.Path(tile, targett)
	if found {
		//		log.Println(e.Name, "found paths:", len(path))
		//for i, v := range path {
		//	log.Println(e.Name, i, v.(*Object).Rect.Center())
		//	}
		e.paths = []pixel.Vec{}
		//log.Printf("Calculated: %s, %v paths", o.Rect.Center(), len(path))
		for _, p := range path {
			e.paths = append(e.paths, p.(*Object).Loc)
		}
		return
	}
	log.Println(e.Name, "no path found, distance:", distance)
}

// PathNeighbors returns the neighbors of the tile, excluding blockers and
// tiles off the edge of the board.
func (o Object) PathNeighbors() []astar.Pather {
	neighbors := []astar.Pather{}
	of := 32.0
	//of = 24.0
	for _, offset := range [][]float64{
		{-of, 0},
		{of, 0},
		{0, -of},
		{0, of},
	} {
		if n := o.w.Tile(pixel.V(o.Rect.Center().X+offset[0], o.Rect.Center().Y+offset[1])); n != nil {
			if n.P.Tile || n.Type == O_TILE {
				neighbors = append(neighbors, n)
			}
		}
	}
	return neighbors
}

func (o *Object) PathEstimatedCost(to astar.Pather) float64 {
	toT := to.(*Object)
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

func (o *Object) PathNeighborCost(to astar.Pather) float64 {
	toT := to.(*Object)
	cost := tileCosts[toT.Type]
	return cost
}

type ObjectType int

const (
	O_TILE ObjectType = iota
	O_BLOCK
	O_INVISIBLE
	O_SPECIAL
	O_WIN
)

// KindCosts map tile kinds to movement costs.

var tileCosts = map[ObjectType]float64{
	O_BLOCK:     30.00,
	O_INVISIBLE: 1.00,
	O_SPECIAL:   0.00,
	O_TILE:      0.00,
	O_WIN:       0.00,
}
