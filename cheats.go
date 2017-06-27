package rpg

func (w *World) RandomLootSomewhere() {
	loot := createLoot()
	w.NewLoot(FindRandomTile(w.Tiles), []Item{loot})
}
