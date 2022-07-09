package rpg

import "github.com/aerth/rpc/librpg/common"

func (w *World) RandomLootSomewhere() {
	loot := createLoot()
	w.NewLoot(common.FindRandomTile(w.Tiles), []Item{loot})
}
