package rpg

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/text"
)

func init() {

	rand.Seed(time.Now().UnixNano())
	log.SetFlags(log.Lshortfile)
	if len(os.Args) == 1 {
		log.SetFlags(0)
	}
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
	RunSpeed: 50.5,
	//Rect:     pixel.R(-8, -8, 8, 8),
	Gravity: 50.00,
	Rate:    2,
}

func NewCharacter() *Character {
	// get main character asset
	sheet, anims, err := LoadCharacterSheet("sprites/char.png", 32)
	if err != nil {
		panic(fmt.Errorf("error loading character sheet: %v", err))
	}
	c := new(Character)
	c.Sheet = sheet
	c.Anims = anims
	//log.Printf("Anims: %v", len(anims))
	c.Sprite = pixel.NewSprite(nil, pixel.Rect{})
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

		f := func(dot pixel.Vec) bool {
			for _, c := range world.Blocks {
				if c.Rect.Contains(dot) && c.Type == O_BLOCK {
					//log.Printf("blocked by: %v at rect: %s, dot: %s", c.SpriteNum, c.Rect, dot)
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
		f2 := func(dot pixel.Vec) bool {
			for _, c := range world.Tiles {
				if c.Type == O_TILE && c.Rect.Contains(dot) {
					return true
				}
			}
			// out of map
			// log.Println("no tile to step on", dot)
			return false
		}

		if f(next.Center()) && f2(next.Center()) {
			//			log.Println("passed:", next)
			char.Rect = next

		} else {
			char.Phys.Vel = pixel.ZV
		}

	}
}

func (char *Character) Damage(n uint, from string) {
	if from != "" {
		from = fmt.Sprintln("from", from)
	}

	if char.Health < n {
		char.Health = 0
		log.Printf("Player took critical hit %s!", from)
		return
	}
	//log.Printf("Player took %v damage %s!", n, from)
	char.Health -= n
	char.W.Message(fmt.Sprintf("ouch! took %v damage, now at %v", n, char.Health))
}

func (char *Character) ResetLocation() {
	char.Rect = DefaultPhys.Rect
	char.Phys.Vel = pixel.ZV

}

type ActionType int

const (
	Talk ActionType = iota
	Slash
	ManaStorm
	MagicBullet
)

func (w *World) Action(char *Character, loc pixel.Vec, t ActionType, dir Direction) {
	switch t {
	case Talk:
		log.Println("nothing to say yet")
	case Slash:
		log.Println("no weapon yet")
	case ManaStorm:
		cost := uint(2)
		if char.Mana < cost {
			w.Message("not enough mana")
			return
		}
		char.Mana -= cost
		w.NewAnimation(char.Rect.Center(), "manastorm", OUT)
	case MagicBullet:
		cost := uint(1)
		if char.Mana < cost {
			w.Message("not enough mana")
			return
		}

		char.Mana -= cost
		w.NewAnimation(char.Rect.Center(), "magicbullet", dir)
	default: //
	}
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
	char.Stats.XP += amount

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
