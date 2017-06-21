package rpg

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
)

var SpriteFrame = pixel.R(-100, -100, 100, 100)

// World holds all information about a world
type World struct {
	Name                    string
	Bounds                  pixel.Rect
	DObjects, Tiles, Blocks []Object                    // sorted
	Background              string                      // path to pic, will be repeated xy if not empty
	background              *pixel.Sprite               // sprite to repeat xy
	Batches                 map[EntityType]*pixel.Batch // one batch for every spritemap
	Color                   pixel.RGBA                  // clear window with this color
	Entities                []*Entity
	Sheets                  map[EntityType]pixel.Picture
	Anims                   map[EntityType]map[EntityState]map[Direction][]pixel.Rect // frankenmap
	Char                    *Character
	Animations              []*Animation
	Messages                []string
	Settings                WorldSettings
}

type WorldSettings struct {
	NumEnemy int
}

func NewWorld(name string, difficulty int) *World {
	w := new(World)
	w.Name = name
	w.Color = RandomColor()
	w.Sheets = make(map[EntityType]pixel.Picture)
	w.Anims = make(map[EntityType]map[EntityState]map[Direction][]pixel.Rect)
	char := NewCharacter()
	char.Inventory = []Item{MakeGold(uint64(rand.Intn(7)))} // start with some loot
	char.W = w
	w.Char = char
	w.Batches = map[EntityType]*pixel.Batch{}
	// create sheets, animations , batch for each sprite map
	for _, t := range []EntityType{SKELETON, SKELETON_GUARD} {
		sheet, anims, err := LoadEntitySheet("sprites/"+t.String()+".png", 13, 21)
		if err != nil {
			panic(fmt.Errorf("error loading sheet: %s %v", t, err))
		}

		w.Sheets[t] = sheet
		w.Anims[t] = anims
		w.Batches[t] = pixel.NewBatch(&pixel.TrianglesData{}, w.Sheets[t])

	}

	log.Println("Loading...")
	if e := w.LoadMap("maps/" + name + ".map"); e != nil {
		log.Println(e)
		return nil
	}
	if len(w.Tiles) == 0 {
		log.Println("Invalid map. No objects found")
		return nil
	}
	char.Rect = char.Rect.Moved(FindRandomTile(w.Tiles))
	return w
}

func (w *World) NewSpecial(o Object) {
	o.Type = O_SPECIAL
	w.DObjects = append(w.DObjects, o)
}
func (w World) String() string {
	return w.Name
}
func (w *World) Update(dt float64) {
	w.checkLevel()
	// clean mobs
	entities := []*Entity{}
	for i := range w.Entities {
		if len(w.Entities) < i || w.Entities[i] == nil {
			continue
		}
		w.Entities[i].ChangeMind(dt)
		w.Entities[i].Update(dt)
		if w.Entities[i].P.Health > 0 {
			entities = append(entities, w.Entities[i])
			continue
		}

		// entity is dead, spawn another
		if len(w.Entities) > 10 {
			continue
		}
		npc := w.NewEntity(SKELETON_GUARD)
		npc.Rect = npc.Rect.Moved(FindRandomTile(w.Tiles))
		entities = append(entities, npc)
		npc = w.NewEntity(SKELETON)
		npc.Rect = npc.Rect.Moved(FindRandomTile(w.Tiles))
		entities = append(entities, npc)
	}
	w.Entities = entities

	// update animations
	if len(w.Animations) > 0 {
		for i := range w.Animations {
			w.Animations[i].update(dt)
		}

	}

	// animations effect entities
	for _, a := range w.Animations {
		if a == nil || time.Since(a.start) < time.Millisecond*300 || time.Since(a.until) > time.Millisecond {
			continue
		}

		for i, v := range w.Entities {
			if a.rect.Contains(v.Rect.Center()) {
				w.Message(fmt.Sprintf("%s took %v damage", v.Name, a.damage))
				w.Entities[i].P.Health -= a.damage
				if w.Entities[i].P.Health <= 0 {
					// entity damage should be function
					if w.Entities[i].P.IsDead {
						w.Entities[i].P.Health = 0
						continue
					}
					w.Entities[i].P.Health = 0
					w.Entities[i].P.IsDead = true
					w.Char.Stats.Kills++

					log.Println("Got new loot!:", FormatItemList(v.P.Loot))
					w.Char.Inventory = StackItems(w.Char.Inventory, v.P.Loot)
					w.Char.ExpUp(v.P.XP)

				}

				log.Printf("%s took %v damage, now at %v HP",
					w.Entities[i].Name, a.damage, w.Entities[i].P.Health)
			}
		}
	}

}

func (w *World) DrawEntity(n int) {
	w.Entities[n].Draw(w.Batches[w.Entities[n].Type], w)

}

// Tile scans tiles and returns the first tile located at dot
func (w *World) Tile(dot pixel.Vec) Object {
	if w.Tiles == nil {
		log.Println("nil tiles")
		return Object{W: w}
	}

	if len(w.Tiles) == 0 {
		log.Println("no tiles")
		return Object{W: w}
	}
	for i := len(w.Tiles) - 1; i >= 0; i-- {
		if w.Tiles[i].Rect.Contains(dot) {
			ob := w.Tiles[i]
			ob.W = w
			return ob
		}
	}
	return Object{W: w}
}

// Block scans blocks and returns the first block located at dot
func (w *World) Block(dot pixel.Vec) Object {
	for i := range w.Blocks {
		if w.Blocks[i].Rect.Contains(dot) {
			return w.Blocks[i]
		}
	}
	return Object{}
}

// Object at location
func (w *World) Object(dot pixel.Vec) Object {
	if w.Blocks != nil {
		for _, v := range w.Blocks {
			if v.Rect.Contains(dot) {
				return v
			}
		}
	}
	for _, v := range w.Tiles {
		if v.Rect.Contains(dot) {
			return v
		}
	}
	return Object{Type: O_NONE}

}

/*
// Object returns the object at dot, very expensive
func (w *World) Object(dot pixel.Vec) Object {

	var ob Object
	for i := range w.Objects {
		if w.Objects[i].Rect.Contains(dot) {
			ob = w.Objects[i]
			if ob.Type == O_BLOCK { // prefer block over tile
				return ob
			}
		}
	}
	return ob

}
*/
func (w *World) Draw(target pixel.Target) {
	for i := range w.Entities {
		w.DrawEntity(i)

	}
	for i := range w.Batches {
		w.Batches[i].Draw(target)
	}

}

func (w *World) ShowAnimations(imd *imdraw.IMDraw) {
	for i := range w.Animations {
		if w.Animations[i] != nil {
			w.Animations[i].draw(imd)
		}
	}

}

func (w *World) HighlightPaths(target pixel.Target) {
	imd := imdraw.New(nil)
	for i := range w.Entities {
		color := TransparentRed
		if len(w.Entities[i].paths) != 0 {
			for _, vv := range w.Entities[i].paths {
				//color = color.Scaled(0.3)
				imd.Color = color
				v := w.Tile(vv)
				imd.Push(v.Rect.Min, v.Rect.Max)
				imd.Rectangle(4)
			}
		}
	}
	imd.Draw(target)

}

func (w *World) Clean() {
	for i := range w.Batches {
		w.Batches[i].Clear()
	}
}

func (w *World) CleanAnimations() {
	anims := []*Animation{}
	now := time.Now().UnixNano()
	for i := range w.Animations {
		if w.Animations[i] != nil && w.Animations[i].until.UnixNano() > now {
			anims = append(anims, w.Animations[i])
		} else {
			//log.Println("removing animation", i)
		}
	}
	w.Animations = anims
}
func (w *World) Reset() {
	w.Char.Health = 255
	w.Char.Stats = DefaultStats
	w.Char.Level = 0
	w.Char.Mana = 0
	w.Char.Inventory = []Item{createLoot()}
	w.Char.Rect = DefaultPhys.Rect.Moved(FindRandomTile(w.Tiles))
	w.Char.Phys.Vel = pixel.ZV
	npc := w.NewEntity(SKELETON_GUARD)
	npc.Rect = npc.Rect.Moved(FindRandomTile(w.Tiles))
	npc2 := w.NewEntity(SKELETON)
	npc2.Rect = npc.Rect.Moved(FindRandomTile(w.Tiles))
	w.Entities = []*Entity{npc, npc2}
	w.Animations = nil

}
