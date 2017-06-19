package rpg

import (
	"fmt"
	"time"

	"golang.org/x/image/colornames"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
)

var SpriteFrame = pixel.R(-100, -100, 100, 100)

// World holds all information about a world
type World struct {
	Name       string
	Bounds     pixel.Rect
	Objects    []*Object
	Background string        // path to pic, will be repeated xy if not empty
	background *pixel.Sprite // sprite to repeat xy
	Batches    map[EntityType]*pixel.Batch
	Color      pixel.RGBA
	Entities   []*Entity
	Sheets     map[EntityType]pixel.Picture
	Anims      map[EntityType]map[EntityState]map[Direction][]pixel.Rect // frankenmap
	Char       *Character
	Animations []*Animation
	Messages   []string
}

func NewWorld(name string, bounds pixel.Rect, testing string) *World {

	w := new(World)
	w.Name = name
	w.Color = RandomColor()
	w.Sheets = make(map[EntityType]pixel.Picture)
	w.Anims = make(map[EntityType]map[EntityState]map[Direction][]pixel.Rect)

	for _, t := range []EntityType{SKELETON, SKELETON_GUARD} {
		sheet, anims, err := LoadEntitySheet("sprites/"+t.String()+".png", 13, 21)
		if err != nil {
			panic(fmt.Errorf("error loading skeleton sheet: %v", err))
		}

		w.Sheets[t] = sheet
		w.Anims[t] = anims
		//log.Printf("New Skeleton Animation Frames: %v", len(w.Anims[t][S_RUN]))
	}

	batchskel := pixel.NewBatch(&pixel.TrianglesData{}, w.Sheets[SKELETON])
	batchguard := pixel.NewBatch(&pixel.TrianglesData{}, w.Sheets[SKELETON_GUARD])

	w.Batches = map[EntityType]*pixel.Batch{
		SKELETON:       batchskel,
		SKELETON_GUARD: batchguard,
	}
	if testing != "" {
		w.LoadMapFile(testing)

	} else {
		w.LoadMap("maps/" + name + ".map")
	}
	w.Messages = []string{"welcome"}
	return w
}

func (w *World) NewSpecial(o *Object) {
	o.Type = O_SPECIAL
	w.Objects = append(w.Objects, o)
}
func (w World) String() string {
	return w.Name
}
func (w *World) Update(dt float64) {
	entities := []*Entity{}
	for i := range w.Entities {
		if w.Entities[i] == nil {
			continue
		}
		w.Entities[i].ChangeMind(dt)
		w.Entities[i].Update(dt)
		if w.Entities[i].P.Health > 0 {
			entities = append(entities, w.Entities[i])
		} else {
			if len(w.Entities) > 64 {
				continue
			}
			npc := w.NewEntity(SKELETON_GUARD)
			npc.Rect = npc.Rect.Moved(FindRandomTile(w.Objects))
			entities = append(entities, npc)
			npc = w.NewEntity(SKELETON)
			npc.Rect = npc.Rect.Moved(FindRandomTile(w.Objects))
			entities = append(entities, npc)

		}
	}
	w.Entities = entities
	tile := w.GetSpecial(w.Char.Rect.Center())
	if tile != nil {
		w.Message("invisible")
		w.Message("full HP")
		w.Char.Health = 255
		w.Char.Invisible = true
	} else {
		w.Char.Invisible = false
	}
	if len(w.Animations) > 0 {
		for i := range w.Animations {
			w.Animations[i].update(dt)

			if w.Animations[i].Type == Magic {
				for _, enemy := range w.Entities {
					w.Animations[i].rect.Contains(enemy.Rect.Center())
				}
			}
		}
	}
}

func (w *World) DrawEntity(n int) {
	w.Entities[n].Draw(w.Batches[w.Entities[n].Type], w)
}
func (w *World) GetSpecial(dot pixel.Vec) *Object {
	for i := range w.Objects {
		if w.Objects[i].Rect.Contains(dot) && w.Objects[i].Type == O_SPECIAL {
			//log.Println("found", w.Objects[i])
			return w.Objects[i]
		}

	}
	return nil

}

func (w *World) Tile(dot pixel.Vec) *Object {
	for i := range w.Objects {
		if w.Objects[i].Rect.Contains(dot) {
			//log.Println("found", w.Objects[i])
			return w.Objects[i]
		}
	}
	return nil

}

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
	color := pixel.ToRGBA(colornames.Red)
	imd.Color = color.Mul(pixel.Alpha(0.2))
	for i := range w.Entities {
		if len(w.Entities[i].paths) != 0 {
			for _, vv := range w.Entities[i].paths {
				v := w.Tile(vv)
				imd.Push(v.Rect.Min, v.Rect.Max)
				imd.Rectangle(0)
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
