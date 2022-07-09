package maps

import (
	"log"

	"github.com/aerth/rpc/librpg/common"
	"github.com/faiface/pixel"
)

type MiniWorld struct {
	Tiles []common.Object
}

// Tile scans tiles and returns the first tile located at dot
func (w *MiniWorld) Tile(dot pixel.Vec) common.Object {
	if w.Tiles == nil {
		log.Println("nil tiles!")
		return common.Object{W: w, Type: common.O_BLOCK}
	}

	if len(w.Tiles) == 0 {
		log.Println("no tiles to look in")
		return common.Object{W: w, Type: common.O_BLOCK}
	}
	for i := len(w.Tiles) - 1; i >= 0; i-- {
		if w.Tiles[i].Rect.Contains(dot) {
			ob := w.Tiles[i]
			ob.W = w
			return ob
		}
	}
	//	log.Println("no tiles found at location:", dot)
	//	panic("bug")
	return common.Object{W: w, Type: common.O_BLOCK}
}
