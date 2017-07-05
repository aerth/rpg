package rpg

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"strconv"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/text"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type Character struct {
	Phys      charPhys // properties
	Stats     Stats
	Sprite    *pixel.Sprite              // current stamp
	Matrix    pixel.Matrix               // location in canvas/map
	Frame     pixel.Rect                 // size (for animation)
	Rect      pixel.Rect                 // size (for collision)
	Dir       Direction                  // Running direction (Idle down)
	Sheet     pixel.Picture              // all frames of animation (4 for each 4 direction, total 16)
	Anims     map[Direction][]pixel.Rect // animation
	Rate      float64                    // animation
	counter   float64                    // in animation
	State     animState                  // Idle or Running
	Inventory []Item                     // inventory
	Health    uint                       // hp
	Mana      uint                       // mp
	Invisible bool                       // hidden from enemies
	Level     uint
	tick      time.Time
	textbuf   *text.Text
	W         *World
}

type charPhys struct {
	RunSpeed float64
	Rect     pixel.Rect
	Vel      pixel.Vec
	Gravity  float64
	CanFly   bool
	Rate     float64
}

var DefaultStats = Stats{
	Intelligence: 60,
	Strength:     60,
	Wisdom:       60,
	Vitality:     60,
}

// DefaultPhys character
var DefaultPhys = charPhys{
	RunSpeed: 60.5,
	Rect:     pixel.R(-8, -8, 8, 8),
	Gravity:  50.00,
	Rate:     2,
}

func (c *Character) StatsReport() string {
	s := (c.Stats.String())
	s += fmt.Sprintf("Level %v\nHealth: %v\nMana: %v\nXP: %v/%v", c.Level, c.Health, c.Mana, c.Stats.XP, c.NextLevel())
	for _, item := range c.Inventory {
		if item.Effect != nil {
			s += fmt.Sprintf("\n Effects: %q", item.Name)

		}
	}
	return s

}

func NewCharacter(skin string) *Character {
	// get main character asset
	sheet, anims, err := LoadCharacterSheet("sprites/"+skin+".png", 32)
	if err != nil {
		panic(fmt.Errorf("error loading character sheet: %v", err))
	}
	c := new(Character)
	c.Sheet = sheet
	c.Anims = anims
	//log.Printf("Anims: %v", len(anims))
	c.Sprite = pixel.NewSprite(nil, pixel.Rect{})
	c.Rect = DefaultPhys.Rect
	c.State = Idle
	c.Frame = c.Anims[DOWN][0]
	c.Phys = DefaultPhys
	c.Rate = 0.1
	c.Health = 255
	c.Mana = 255
	c.Stats = DefaultStats
	return c
}

func (char *Character) Draw(t pixel.Target) {
	if char.Sprite == nil {
		char.Sprite = pixel.NewSprite(nil, pixel.Rect{})
	}
	// draw the correct frame with the correct position and direction
	char.Sprite.Set(char.Sheet, char.Frame)
	char.Sprite.Draw(t, char.Matrix)

}

func (char *Character) Update(dt float64, dir Direction, world *World) {
	for _, i := range char.Inventory {
		if i.Effect != nil {
			char.Stats = i.Effect(char.Stats)
		}

	}
	if time.Since(char.tick) >= time.Second*1 {
		if char.Mana < 255 {
			char.Mana++
		}
		char.tick = time.Now()
	}
	char.counter += dt
	// determine the new animation state

	var newState animState

	switch {
	default:
		newState = char.State
	case -2 < char.Phys.Vel.Len() && char.Phys.Vel.Len() < 2:
		char.Phys.Vel = pixel.ZV
		newState = Idle
	case char.Phys.Vel.Len() > 2:
		newState = Running
	case char.Phys.Vel.Len() < -2:
		newState = Running
	}

	// reset the time counter if the state changed
	if char.State != newState || char.Dir != dir {
		char.State = newState
		char.Dir = dir
		//log.Println(char.State, char.Dir)
		char.counter = 0
	}

	// determine the correct animation frame

	if char.State == Idle {
		char.Frame = char.Anims[dir][0]
	} else if char.State == Running {
		// count 0 1 2 3 0 1 2 3...
		i := int(math.Floor(char.counter / char.Rate))

		char.Frame = char.Anims[dir][i%len(char.Anims[dir])]

		// gradually lose momentum
		char.Phys.Vel = pixel.Lerp(char.Phys.Vel, pixel.ZV, 1-math.Pow(1.0/char.Phys.Gravity, dt))
		next := char.Rect.Moved(char.Phys.Vel.Scaled(dt))

		f := func(nextblock pixel.Rect) bool {
			for _, c := range world.Blocks {
				area := c.Rect.Norm().Intersect(nextblock.Norm()).Norm().Area()
				if area != 0 {
					//				log.Printf("%f %s %s blocked by: %v at rect: %s", area, char.Rect, nextblock, c.SpriteNum, c.Rect)
					return false
				}
			}
			return true
		}
		if char.Phys.CanFly {
			char.Rect = next
			return
		}
		// only walk on tiles
		f2 := func(nexttile pixel.Rect) bool {

			for _, c := range world.Tiles {
				if c.Type == O_TILE && c.Rect.Intersect(nexttile).Norm().Area() != 0 {
					return true
				}
			}
			// out of map
			// log.Println("no tile to step on", dot)
			return false
		}

		if f(next) && f2(next) {
			//			log.Println("passed:", next)
			char.Rect = next

		} else {
			char.Phys.Vel = pixel.ZV
		}

	}
}

func (char *Character) Damage(n uint, from string) {
	if from != "" {
		from = fmt.Sprintf("from %s", from)
	}
	if char.Health < n {
		char.Health = 0
		log.Printf("Player took critical hit %s!", from)
		return
	}
	char.Health -= n
	char.W.Message(fmt.Sprintf("ouch! took %v damage, now at %v", n, char.Health))
}

func (char *Character) ResetLocation() {
	char.Rect = DefaultPhys.Rect
	char.Phys.Vel = pixel.ZV

}

func (char *Character) CountGold() string {
	var madlootyo uint64
	for _, item := range char.Inventory {
		if item.Type == GOLD {
			madlootyo += item.Quantity
		}
	}
	return strconv.FormatInt(int64(madlootyo), 10)
}

func (char *Character) ExpUp(amount uint64) {
	log.Println("Gained experience:", amount)
	char.Stats.XP += amount
	char.Stats.Score += amount

}

func (w *World) checkLevel() {

	if w.Char.Stats.XP == 0 {
		return
	}
	nextlvl := w.Char.NextLevel()
	if w.Char.Stats.XP > nextlvl {
		w.Char.Level++
		w.Char.Health = 255
		w.Message("LVL UP")
		log.Printf("level up (%v)! next lvl at %v xp", w.Char.Level, nextlvl)
		if xp := w.Char.Stats.XP - nextlvl; xp > 0 {
			w.Char.Stats.XP = xp
		} else {
			w.Char.Stats.XP = 0

		}
		switch w.Char.Level {
		default:
			w.Char.Stats.Intelligence += float64(10 * w.Char.Level)
		}
		log.Println(w.Char.Stats)
	}
}

func (c *Character) NextLevel() uint64 {
	return uint64(150 * c.Level)
}

func (c *Character) MaxHealth() uint64 {
	return uint64(c.Health * c.Level)
}

func (c *Character) PickUp(items []Item) string {

	c.Inventory = StackItems(c.Inventory, items)
	return fmt.Sprint(items)

}
