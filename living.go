package rpg

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"strconv"
	"time"

	"github.com/faiface/pixel"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type Entity struct {
	Name      string
	Type      EntityType
	CanFly    bool
	Friendly  bool
	Rate      float64
	State     animState
	Frame     pixel.Rect
	Matrix    pixel.Matrix
	Rect      pixel.Rect
	SpriteNum int `json:"Sprite"`
	Program   EntityState
	P         EntityProperties
	Phys      ePhys
	Dir       Direction
	counter   float64
	paths     []pixel.Vec
	w         *World
}

type EntityType int
type EntityState int

const (
	SKELETON EntityType = iota
	SKELETON_GUARD
	DOBJECT
)

type EntityProperties struct {
	Health float64
	Mana   float64
	Loot   []Item
	IsDead bool
}

const (
	S_IDLE EntityState = iota
	S_RUN
	S_WANDER
	S_GUARD
	S_SUSPECT
	S_HUNT
)

func (e *Entity) String() string {

	return fmt.Sprintf("%s at %v,%v", e.Name, int(e.Rect.Center().X), int(e.Rect.Center().Y))
}

func (w *World) NewEntity(t EntityType) *Entity {
	n := len(w.Entities)
	var e *Entity
	switch t {
	default: // no default
	case SKELETON, SKELETON_GUARD:
		if w.Sheets[t] == nil || w.Anims[t] == nil {
			sheet, anims, err := LoadEntitySheet("sprites/"+t.String()+".png", 13, 21)
			if err != nil {
				panic(fmt.Errorf("error loading skeleton sheet: %v", err))
			}
			w.Sheets[t] = sheet
			w.Anims[t] = anims
			log.Printf("New Skeleton Animation Frames: %v", len(anims[S_RUN]))
		}

		e = &Entity{
			Name: fmt.Sprintf("%s #%v", t, n),
			w:    w,
			Type: t,
			P: EntityProperties{
				Health: float64(rand.Intn(255)),
				Mana:   float64(rand.Intn(255)),
			},
			Rect:  pixel.R(-16, -16, 16, 16),
			State: Running,
			Frame: w.Anims[t][S_RUN][DOWN][0],
			Phys:  DefaultMobPhys,
			Rate:  0.1,
		}
	}

	if e == nil {
		return nil
	}

	e.P.Loot = RandomLoot()

	w.Entities = append(w.Entities, e)

	return e
}

type Item struct {
	Name       string
	Type       ItemType
	Properties ItemProperties
	Quantity   uint64
}

func (i Item) String() (s string) {

	if i.Quantity > 1 {
		s += strconv.FormatInt(int64(i.Quantity), 10)
		s += " "
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

type ItemProperties struct {
	Weight uint8
}

type ePhys struct {
	RunSpeed float64
	Rect     pixel.Rect
	Vel      pixel.Vec
	Gravity  float64
	Rate     float64
}

// DefaultPhys character
var DefaultMobPhys = ePhys{
	RunSpeed: 100.5,
	//Rect:     pixel.R(-8, -8, 8, 8),
	Rect:    pixel.R(98, 98, 108, 108),
	Gravity: 50.00,
	Rate:    2,
}

func (e *Entity) Draw(t pixel.Target, w *World) {

	sprite := pixel.NewSprite(nil, pixel.Rect{})
	// draw the correct frame with the correct position and direction
	sprite.Set(w.Sheets[e.Type], e.Frame)
	sprite.Draw(t, pixel.IM.Moved(e.Rect.Center()))
}

func (e *Entity) ChangeMind(dt float64) {
	tile := e.w.Tile(e.Rect.Center())
	if tile == nil || tile.Type != O_TILE {
		e.Rect = DefaultMobPhys.Rect.Moved(FindRandomTile(e.w.Objects))
	}

	r := pixel.Rect{e.Rect.Center(), e.w.Char.Rect.Center()}
	if r.Size().Len() < 48 {
		e.w.Char.Damage(uint8(rand.Intn(10)), e.Name)
		return
	}

	if e.CanFly {
		if !e.w.Char.Invisible {
			e.Phys.Vel = e.Rect.Center().Sub(e.w.Char.Rect.Center()).Unit().Scaled(e.Phys.RunSpeed)
		} else {
			e.Phys.Vel = pixel.ZV
		}

		return
	}
	//log.Println("finding path", e.Name)

	if len(e.paths) > 2 {
		e.Phys.Vel = e.Rect.Center().Sub(e.paths[len(e.paths)-2]).Unit().Scaled(e.Phys.RunSpeed)
		e.paths = e.paths[:len(e.paths)-1]
		//		e.Phys.Vel = e.paths[0].Unit().Scaled(e.Phys.RunSpeed)
		//log.Println("got vel:", e.Phys.Vel, e.Phys.Vel.Len())
	} else {
		if !e.w.Char.Invisible {
			char := e.w.Char
			if char == nil {
				log.Println("nil character")
				return
			}

			target := e.w.Tile(char.Rect.Center())
			if target == nil {
				log.Println("nil target tile")
				return
			}

			e.pathcalc(target.Rect.Center())

			if len(e.paths) == 0 {
				log.Println("no paths?!")
			}
		} else {
			e.Phys.Vel = pixel.ZV

		}
	}
}

func (e *Entity) Update(dt float64) {
	for {
		tile := e.w.Tile(e.Rect.Center())
		if tile == nil || tile.Type == O_BLOCK {
			e.Rect = DefaultMobPhys.Rect.Moved(FindRandomTile(e.w.Objects))
		} else {
			break
		}
	}

	e.counter += dt
	collide := append(e.w.Objects, e.w.DObjects...)
	w := e.w
	i := int(math.Floor(e.counter / e.Rate))
	//frame := i % len(e.Anims[e.Program][e.Dir])
	if e.Phys.Vel.X != 0 || e.Phys.Vel.Y != 0 {
		e.Program = S_RUN
	} else {
		e.Program = S_IDLE
	}
	if len(w.Anims[e.Type][e.Program][e.Dir]) == 0 {
		log.Println("bad sprite:", e.Name, e.Program, e.Dir)
		return
	}
	e.Frame = w.Anims[e.Type][e.Program][e.Dir][i%len(w.Anims[e.Type][e.Program][e.Dir])]
	//log.Println(e.Name, "Frame#", frame, "out of", len(e.Anims[e.Program][e.Dir]))

	next := e.Rect.Moved(e.Phys.Vel.Scaled(-dt))
	t := w.Tile(next.Center())
	if t == nil && !e.CanFly {

		return
	}
	if !e.CanFly && t.Type == O_BLOCK {
		next = e.Rect.Moved(e.Phys.Vel.Scaled(dt))
	}

	if e.Phys.Vel == pixel.ZV {
		return
	}
	//	log.Println(e.Name, "wants to go", next.Center(), "from", e.Rect.Center())
	f := func(dot pixel.Vec) bool {
		if e.CanFly {
			return true
		}
		for _, c := range collide {
			if c.P.Block && c.Rect.Contains(dot) {
				return false
			}
		}
		return true
	}
	// only walk on tiles
	f2 := func(dot pixel.Vec) bool {
		if e.CanFly {
			return true
		}
		for _, c := range collide {
			if c.Type == O_TILE && c.Rect.Contains(dot) {
				return true
			}

		}
		return false
	}
	if f(next.Center()) && f2(next.Center()) {
		//log.Println(e.Name, "passed:", next)
		e.Rect = next
	} else {
		log.Println("cant move", e.Name, "to ", next.Center(), w.Tile(next.Center()), e.paths[0])
		if len(e.paths) > 1 {
			e.paths = e.paths[:len(e.paths)-1]
		} else {
			e.pathcalc(w.Char.Rect.Center())
		}
	}

}

// loadCharacterSheet returns an animated spritesheet
// 13W 21H
func LoadEntitySheet(sheetPath string, framesx, framesy uint8) (sheet pixel.Picture, anims map[EntityState]map[Direction][]pixel.Rect, err error) {
	sheet, err = LoadPicture(sheetPath)
	frameWidth := float64(int(sheet.Bounds().Max.X / float64(framesx)))
	frameHeight := float64(int(sheet.Bounds().Max.Y / float64(framesy)))
	//log.Println(frameWidth, "width", frameHeight, "height")
	// create a array of frames inside the spritesheet
	var frames = []pixel.Rect{}
	for y := 0.00; y+frameHeight <= sheet.Bounds().Max.Y; y = y + frameHeight {
		for x := 0.00; x+float64(frameWidth) <= sheet.Bounds().Max.X; x = x + float64(frameWidth) {
			frames = append(frames, pixel.R(
				x,
				y,
				x+frameWidth,
				y+frameHeight,
			))
		}
	}

	//log.Println("total skeleton frames", len(frames))

	// 0-5 die
	// BLANK 6-12
	// 13-25 shoot right
	// 26-39 shoot down
	// 6-76 shoot left
	// 7-25 shoot up
	anims = make(map[EntityState]map[Direction][]pixel.Rect)
	anims[S_IDLE] = make(map[Direction][]pixel.Rect)
	anims[S_WANDER] = make(map[Direction][]pixel.Rect)
	anims[S_RUN] = make(map[Direction][]pixel.Rect)
	anims[S_GUARD] = make(map[Direction][]pixel.Rect)
	anims[S_SUSPECT] = make(map[Direction][]pixel.Rect)
	anims[S_HUNT] = make(map[Direction][]pixel.Rect)

	// spritesheet is right down left up
	// why is inverted?
	anims[S_IDLE][LEFT] = frames[143:144]
	anims[S_IDLE][UP] = frames[156:157]
	anims[S_IDLE][RIGHT] = frames[169:170]
	anims[S_IDLE][DOWN] = frames[182:183]
	anims[S_RUN][LEFT] = frames[143:152]
	anims[S_RUN][UP] = frames[156:165]
	anims[S_RUN][RIGHT] = frames[169:178]
	anims[S_RUN][DOWN] = frames[182:191]
	return sheet, anims, nil
}

func (w *World) NewMobs(n int) {
	if w.Settings.NumEnemy == 0 {
		w.Settings.NumEnemy = n
	}
	if n != 0 {
		npc := w.NewEntity(SKELETON_GUARD)
		npc.Phys.RunSpeed = 10
		npc.P.Health = 2000
		// npc.CanFly = true
		npc.Rect = npc.Rect.Moved(FindRandomTile(w.Objects))

		for i := 1; i < n; i++ {
			npc = w.NewEntity(SKELETON)
			npc.Rect = npc.Rect.Moved(FindRandomTile(w.Objects))
		}

	}

}
