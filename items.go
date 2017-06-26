package rpg

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"

	"github.com/faiface/pixel"
)

type Item struct {
	Name       string
	Type       ItemType
	Properties ItemProperties
	Effect     func(Stats) Stats
	Quantity   uint64
}

type ItemProperties struct {
	Weight uint8
	Size   pixel.Vec
}

func DefaultItemProperties() ItemProperties {
	return ItemProperties{}
}

func (i Item) String() (s string) {
	if i.Quantity > 1 {
		s += strconv.FormatInt(int64(i.Quantity), 10)
		s += " "
	}
	if i.Name != "" {
		s += i.Name
		if i.Effect != nil {
			switch i.Type {
			case ARMOR:
				s += " (+defence)"
			case POTION:
				s += " (+magic)"
			case WEAPON:
				s += " (+attack)"
			}
		}
		return s
	}
	s += i.Type.String()
	return s

}

type ItemType int

const (
	_ ItemType = iota
	GOLD
	POTION
	FOOD
	WEAPON
	ARMOR
	SPECIAL
)

func MakeGold(amount uint64) Item {
	return Item{
		Name:     "gold",
		Type:     GOLD,
		Quantity: amount,
	}
}

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
		if i%4 == 0 {
			s += "\n"
		}
		s += item.String() + ", "
	}

	return strings.TrimSuffix(s, ", ")
}

func RandomLoot() []Item {
	switch rand.Intn(10) {
	default: // 0
		return []Item{} // no loot
	case 1, 2, 3, 4:
		return []Item{MakeGold(uint64(rand.Intn(255) + 1))}
	case 5, 6:
		return []Item{createLoot()}
		//	case 8:
	case 7, 8, 9:
		return []Item{RandomMagicItem()}
	}
}

func RandomMagicItem() Item {
	// special item
	item := createLoot()
	item.Name = fmt.Sprintf(GenerateItemName(), item.Type.String())
	item.Effect = RandomItemEffect(item.Type)
	return item

}

func RandomItemEffect(t ItemType) func(Stats) Stats {
	fn := func(s Stats) Stats {
		return s
	}

	switch t {
	case ARMOR:
		fn = func(s Stats) Stats {
			s.Vitality += 0.1
			return s
		}
	case WEAPON:
		fn = func(s Stats) Stats {
			s.Strength += 0.1
			return s
		}
	case POTION:
		fn = func(s Stats) Stats {
			s.Intelligence += 0.1
			return s
		}
	default:
		fn = func(s Stats) Stats {
			return s
		}

	}
	return fn
}
