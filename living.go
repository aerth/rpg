package rpg

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"time"

	"golang.org/x/image/colornames"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

//var DefaultEntityRectangle = pixel.R(-16, -8, 16, 24)

var DefaultEntityRectangle = pixel.R(-16, -16, 16, 16)

//var DefaultSpriteRectangle = pixel.R(-16, 0, 16, 32)

type Entity struct {
	Name       string
	Type       EntityType
	CanFly     bool
	Friendly   bool
	Rate       float64
	State      animState
	Frame      pixel.Rect
	Matrix     pixel.Matrix
	Rect       pixel.Rect
	SpriteNum  int `json:"Sprite"`
	Program    EntityState
	P          EntityProperties
	Phys       ePhys
	Dir        Direction
	counter    float64
	paths      []pixel.Vec
	w          *World
	calculated time.Time
	imd        *imdraw.IMDraw
	ticker     <-chan time.Time
}

type EntityType int
type EntityState int

const (
	SKELETON EntityType = iota
	SKELETON_GUARD
	DOBJECT
)

type EntityProperties struct {
	XP                uint64
	Health, MaxHealth float64
	Mana              float64
	Loot              []Item
	IsDead            bool
	Strength          float64
	AttackSpeed       uint64
}

const (
	S_IDLE EntityState = iota
	S_RUN
	S_WANDER
	S_GUARD
	S_SUSPECT
	S_HUNT
	S_DEAD
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
			log.Println("New sheet:", t)
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
				Health:      255,
				Mana:        255,
				Strength:    2,
				XP:          10,
				MaxHealth:   255,
				AttackSpeed: 550,
			},
			Rect:  DefaultEntityRectangle,
			State: Running,
			Frame: w.Anims[t][S_RUN][DOWN][0],
			Phys:  DefaultMobPhys,
			Rate:  0.1,
		}
		e.ticker = time.Tick(time.Millisecond * time.Duration(e.P.AttackSpeed))
	}

	if e == nil {
		return nil
	}

	e.P.Loot = RandomLoot()

	w.Entities = append(w.Entities, e)

	return e
}

type ePhys struct {
	RunSpeed float64
	Vel      pixel.Vec
	Gravity  float64
	Rate     float64
}

// DefaultPhys character
var DefaultMobPhys = ePhys{
	RunSpeed: 40.5,
	Gravity:  50.00,
	Rate:     2,
}

func (e *Entity) Draw(t pixel.Target, w *World) {

	sprite := pixel.NewSprite(nil, pixel.Rect{})
	// draw the correct frame with the correct position and direction

	scaling := 0.5
	if e.Type == SKELETON_GUARD {
		scaling = 0.7
	}
	sprite.Set(w.Sheets[e.Type], e.Frame)
	sprite.Draw(t, pixel.IM.Scaled(pixel.ZV, scaling).Moved(e.Rect.Center()))
	//sprite.Draw(t, pixel.IM.Scaled(pixel.ZV, 0.5).Moved(e.Rect.Center()))

	// HP bars
	if e.imd == nil {
		e.imd = imdraw.New(nil)
	}
	e.imd.Clear()
	rect := e.Rect.Norm()
	rect.Max.Y = rect.Min.Y + 2
	rect.Max.X = rect.Min.X + 30
	if e.P.Health > 0 {
		DrawBar(e.imd, colornames.Red, e.P.Health, e.P.MaxHealth, rect)
		e.imd.Draw(t)
	}
	/* good debug square
	e.imd.Color = colornames.Green
	e.imd.Push(e.Rect.Min, e.Rect.Max)
	e.imd.Rectangle(1)
	e.imd.Draw(t)
	*/
}

func (e *Entity) Center() pixel.Vec {
	return e.Rect.Center()
}

func (e *Entity) ChangeMind(dt float64) {
	if t := e.w.Tile(e.Center()); t.Type != O_TILE {
		e.Phys.Vel = pixel.ZV
		return
	}
	if e.w.Char.Invisible {
		e.Phys.Vel = pixel.ZV
		return
	}

	r := pixel.Rect{e.Rect.Center(), e.w.Char.Rect.Center()}
	if r.Size().Len() < e.Rect.Size().Len()/2 {
		//	log.Println("in attack range", r.Size().Len())
		e.Phys.Vel = e.Rect.Center().Sub(e.w.Char.Rect.Center()).Unit().Scaled(-e.Phys.RunSpeed)
		select {

		case <-e.ticker:
			e.w.Char.Damage(uint(rand.Intn(10*int(e.P.Strength))), e.Name)
		default:
		}

		return
	}

	if e.CanFly {
		log.Println("FLYING", e.Name)
		if !e.w.Char.Invisible {

			e.Phys.Vel = e.Rect.Center().Sub(e.w.Char.Rect.Center()).Unit().Scaled(-e.Phys.RunSpeed)
		} else {
			e.Phys.Vel = pixel.ZV
		}

		return
	}
	e.pathcalc(e.w.Char.Rect.Center())
	if len(e.paths) > 2 {
		e.Phys.Vel = e.Rect.Center().Sub(e.paths[len(e.paths)-2]).Unit().Scaled(-e.Phys.RunSpeed)
	}
}

func (e *Entity) Update(dt float64) {
	blk := e.w.Tile(e.Rect.Center())
	if blk.Type != O_TILE {
		old := e.Rect.Center()
		e.Rect = DefaultEntityRectangle.Moved(TileNear(e.w.Tiles, e.Center()).Loc)
		log.Println("Moved skel:", old, "to", e.Rect.Center())
		return
	}

	e.counter += dt
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
	// choose frame
	e.Frame = w.Anims[e.Type][e.Program][e.Dir][i%len(w.Anims[e.Type][e.Program][e.Dir])]

	// move
	next := e.Rect.Moved(e.Phys.Vel.Scaled(dt))
	t := w.Tile(next.Center())
	if t.Type == O_NONE && !e.CanFly {
		return
	}
	if !e.CanFly && t.Type == O_BLOCK {
		log.Println(e.Type, "got blocked", t.Loc)
		next = e.Rect.Moved(e.Phys.Vel.Scaled(-dt * 10))
		if w.Tile(next.Center()).Type != O_TILE {
			log.Println("returning")
			return
		}
	}

	//	log.Println(e.Name, "wants to go", next.Center(), "from", e.Rect.Center())
	f := func(dot pixel.Vec) bool {
		if e.CanFly {
			return true
		}
		for _, c := range w.Blocks {
			if c.Type == O_BLOCK && c.Rect.Contains(dot) {
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
		for _, c := range w.Tiles {
			if c.Type == O_TILE && c.Rect.Contains(dot) {
				return true
			}

		}
		return false
	}
	if f(next.Center()) && f2(next.Center()) {
		e.Rect = next
	} else {
		//log.Println("cant move", e.Name, "to ", next.Center(), w.Tile(next.Center()), e.paths[0])
		if len(e.paths) > 0 {
			e.paths = e.paths[:len(e.paths)-1]
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
	anims[S_DEAD] = make(map[Direction][]pixel.Rect)

	// spritesheet is right down left up
	anims[S_DEAD][LEFT] = frames[0:5]
	anims[S_DEAD][RIGHT] = frames[0:5]
	anims[S_DEAD][UP] = frames[0:5]
	anims[S_DEAD][DOWN] = frames[0:5]
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
		npc.Rect = npc.Rect.Moved(FindRandomTile(w.Tiles))

		for i := 1; i < n; i++ {
			npc = w.NewEntity(SKELETON)
			npc.Rect = npc.Rect.Moved(FindRandomTile(w.Tiles))
		}

	}

}
