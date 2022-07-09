package common

import (
	"log"
	"math/rand"

	"github.com/faiface/pixel"
)

//var DefaultSpriteRectangle = pixel.R(-16, 0, 16, 32)
//var DefaultSpriteRectangle = pixel.R(-16, 0, 16, 32)

// assumes only tiles are given
func FindRandomTile(os []Object) pixel.Vec {
	if len(os) == 0 {
		panic("no objects")
	}
	ob := os[rand.Intn(len(os))]
	if ob.Loc != pixel.ZV && ob.SpriteNum != 0 && ob.Type == O_TILE {
		return ob.Rect.Center()
	}
	return FindRandomTile(os)
}

func GetObjects(objects []Object, position pixel.Vec) []Object {
	var good []Object
	for _, o := range objects {
		if o.Rect.Contains(position) {
			good = append(good, o)
		}
	}
	return good
}

func GetTiles(objects []Object) []Object {
	var tiles []Object
	for _, o := range objects {
		if o.Type == O_TILE {
			tiles = append(tiles, o)
		}
	}
	return tiles
}

func TilesAt(objects []Object, position pixel.Vec) []Object {
	var good []Object
	all := GetObjects(objects, position)
	if len(all) > 0 {
		for _, o := range all {
			if DefaultSpriteRectangle.Moved(o.Loc).Contains(position) && o.Type == O_TILE {
				good = append(good, o)
			}

		}
	}
	return good

}
func GetObjectsAt(objects []Object, position pixel.Vec) []Object {
	var good []Object
	all := GetObjects(objects, position)
	if len(all) > 0 {
		for _, o := range all {
			if DefaultSpriteRectangle.Moved(o.Loc).Contains(position) {
				good = append(good, o)
			}

		}
	}
	return good

}

func GetBlocks(objects []Object, position pixel.Vec) []Object {
	var bad []Object
	all := GetObjects(objects, position)
	if len(all) > 0 {
		for _, o := range all {
			if o.Type == O_BLOCK {
				bad = append(bad, o)
			}

		}
	}
	return bad
}

// GetNeighbors gets the neighboring tiles
func (o Object) GetNeighbors() []Object {
	neighbors := []Object{}
	of := 32.0
	for _, offset := range [][]float64{
		{-of, 0},
		{of, 0},
		{0, -of},
		{0, of},
	} {
		if n := o.W.Tile(pixel.V(o.Rect.Center().X+offset[0], o.Rect.Center().Y+offset[1])); n.Type == o.Type {
			neighbors = append(neighbors, n)
		}
	}
	return neighbors

}
func TileNear(all []Object, loc pixel.Vec) Object {
	tile := TilesAt(all, loc)
	snap := 32.00
	loc.X = float64(int(loc.X/snap)) * snap
	loc.Y = float64(int(loc.Y/snap)) * snap
	radius := 1.00
	oloc := loc
	if len(tile) > 0 {
		oloc = tile[0].Loc
	}
	log.Println("looking for loc:", loc)
	for i := 0; i < len(all); i++ {
		loc.X += radius * 16
		loc.Y += radius * 16
		if loc == oloc {
			continue
		}
		log.Println("Checking loc", loc)
		os := TilesAt(all, loc)
		if len(os) > 0 {
			if os[0].Loc == pixel.ZV || os[0].Loc == oloc {
				continue
			}
			return os[0]
		}
		loc2 := loc
		loc2.X = -loc.X
		os = TilesAt(all, loc.Scaled(-1))
		if len(os) > 0 {
			if os[0].Loc == pixel.ZV || os[0].Loc == oloc {
				continue
			}
			return os[0]
		}
		os = TilesAt(all, loc2.Scaled(1))
		if len(os) > 0 {
			if os[0].Loc == pixel.ZV || os[0].Loc == oloc {
				continue
			}
			return os[0]
		}

		os = TilesAt(all, loc2.Scaled(-1))
		if len(os) > 0 {
			if os[0].Loc == pixel.ZV || os[0].Loc == oloc {
				continue
			}
			return os[0]
		}

		if i%4 == 0 {
			radius++
			loc.X += 16
		}
	}
	log.Println("not found")
	return Object{Type: O_NONE}

}
