package rpg

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"time"

	"golang.org/x/image/colornames"

	"github.com/aerth/rpg/assets"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type Object struct {
	Loc       pixel.Vec        `json:", omitempty"`
	Rect      pixel.Rect       `json:", omitempty"`
	Type      ObjectType       `json:"-"`
	P         ObjectProperties `json:", omitempty"`
	SpriteNum int              `json:"Sprite,omitempty"`
	Sprite    *pixel.Sprite    `json:"-"`
	w         *World           `json:"-"`
}

func (o Object) String() string {
	return fmt.Sprintf("%s %s %s %v", o.Loc, o.Rect, o.Type, o.SpriteNum)
}

type ObjectProperties struct {
	Invisible bool `json:",omitempty"`
	Tile      bool `json:",omitempty"`
	Block     bool `json:",omitempty"`
	Special   bool `json:",omitempty"`
}

func NewTile(loc pixel.Vec) Object {
	return Object{
		Loc:  loc,
		Rect: pixel.Rect{loc.Sub(pixel.V(16, 16)), loc.Add(pixel.V(16, 16))},
		P: ObjectProperties{
			Tile: true,
		},
	}
}
func NewBlock(loc pixel.Vec) Object {
	return Object{
		Loc:  loc,
		Rect: pixel.Rect{loc.Sub(pixel.V(16, 16)), loc.Add(pixel.V(16, 16))},
		P: ObjectProperties{
			Block: true,
		},
	}
}

func NewTileBox(rect pixel.Rect) Object {
	return Object{
		Rect: rect,
		P: ObjectProperties{
			Tile: true,
		},
	}
}
func NewBlockBox(rect pixel.Rect) Object {
	return Object{
		Rect: rect,
		P: ObjectProperties{
			Block: true,
		},
	}
}
func (o Object) Highlight(win pixel.Target) {
	imd := imdraw.New(nil)
	color := pixel.ToRGBA(colornames.Red)
	if o.P.Tile {
		color = pixel.ToRGBA(colornames.Blue)
	}
	imd.Color = color.Scaled(0.2)
	imd.Push(o.Rect.Min, o.Rect.Max)
	imd.Rectangle(0)
	imd.Draw(win)
}
func (o Object) Draw(win pixel.Target, spritesheet pixel.Picture, sheetFrames []*pixel.Sprite) {
	if o.P.Invisible {
		return
	}

	if o.Sprite == nil {
		if 0 > o.SpriteNum && o.SpriteNum > len(sheetFrames) {
			log.Printf("unloadable sprite: %v/%v", o.SpriteNum, len(sheetFrames))
			return
		}
		o.Sprite = sheetFrames[o.SpriteNum]
	}
	if o.Loc == pixel.ZV && o.Rect.Size().Y != 32 {
		log.Println(o.Rect.Size(), "cool")
		DrawPattern(win, o.Sprite, o.Rect, 100)
	} else {
		o.Sprite.Draw(win, pixel.IM.Moved(o.Loc))
	}

}
func (w *World) LoadMapFile(path string) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		log.Println("error loading map:", err)
		w.Exit(111)
	}
	w.loadmap(b)
}
func (w *World) LoadMap(path string) {
	b, err := assets.Asset(path)
	if err != nil {
		log.Println("error loading map:", err)
		w.Exit(111)
	}
	w.loadmap(b)
}
func (w *World) loadmap(b []byte) {
	var things = []Object{}
	err := json.Unmarshal(b, &things)
	if err != nil {
		log.Println("invalid map:", err)
		w.Exit(111)
	}
	for _, thing := range things {
		t := new(Object)
		*t = thing
		t.w = w
		if t.P.Block {
			t.Type = O_BLOCK
		}
		if t.P.Tile {
			t.Type = O_TILE
		}
		w.Objects = append(w.Objects, t)
	}
	return
}

func (o ObjectType) MarshalJSON() ([]byte, error) {
	i := int(o)
	return json.Marshal(i)
}

func (o ObjectType) UnmarshalJSON(b []byte) error {
	var i int
	err := json.Unmarshal(b, &i)
	if err != nil {
		return err
	}
	o = ObjectType(i)
	return nil
}

func FindRandomTile(os []*Object) pixel.Vec {

	for i := 0; i < len(os); i++ {
		o := os[rand.Intn(len(os))]
		if n := len(o.PathNeighbors()); n > 2 {
			log.Println("returning", o.Rect.Center())
			return o.Rect.Center()
		}
	}
	log.Println("no tiles? looking again:", len(os))
	for _, o := range os {
		if o.Type == O_TILE && rand.Intn(100) < 10 {
			return o.Rect.Center()
		}
	}
	return pixel.ZV
}
