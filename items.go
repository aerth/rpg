package rpg

import (
	"math/rand"
	"strings"
)

func createLoot() Item {

	item := Item{
		Type:     ItemType(rand.Intn(5)) + 1,
		Quantity: 1,
	}

	return item
}

func createItemLoot() Item {
	item := createLoot()
	if item.Type == GOLD || item.Type == FOOD {
		return createItemLoot()
	}
	return item
}

// just stack gold potions and food for now
func (i Item) Stack(items []Item) []Item {
	var stacked []Item
	var foods, potions, goldlvl uint64

	//log.Println(items)
	//log.Println("Totalnum:", len(items))
	for _, item := range items {
		//log.Println(item)
		switch item.Type {
		case FOOD:
			foods += item.Quantity
		case POTION:
			potions += item.Quantity
		case GOLD:
			//log.Printf("adding %v and %v", goldlvl, item.Quantity)
			goldlvl += item.Quantity
		default:
			//log.Printf("stacking item %s", item)
			stacked = append(stacked, item)
			//log.Printf("items are now %v", len(stacked))
		}
	}

	switch i.Type {
	default:
		stacked = append(stacked, i)
	case GOLD:
		goldlvl += i.Quantity
	case POTION:
		potions += i.Quantity
	case FOOD:
		foods += i.Quantity
	}
	if potions > 0 {
		stacked = append(stacked, Item{Type: POTION, Quantity: potions})
	}
	if foods > 0 {

		stacked = append(stacked, Item{Type: FOOD, Quantity: foods})
	}
	if goldlvl > 0 {
		stacked = append(stacked, MakeGold(goldlvl))
	}

	return stacked

}

func StackItems(itemsets ...[]Item) []Item {
	var backpack []Item
	for _, inventory := range itemsets {
		for _, item := range inventory {
			backpack = item.Stack(backpack)
		}
	}
	return backpack
}

func FormatItemList(items []Item) string {
	if len(items) == 0 {
		return "none"
	}
	var s string
	for i, item := range items {
		if i%5 == 0 {
			s += "\n"
		}
		s += item.String() + ", "
	}

	return strings.TrimSuffix(s, ", ")
}

func RandomLoot() []Item {
	switch rand.Intn(5) {
	default: // 0
		return []Item{} // no loot
	case 1, 2:
		return []Item{MakeGold(uint64(rand.Intn(255) + 1))}
	case 3, 4:
		return []Item{createLoot()}

	}
}
