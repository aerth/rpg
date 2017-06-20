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

var DefaultSpriteRectangle = pixel.R(-16, -16, 16, 16)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type Object struct {
	Loc       pixel.Vec        `json:"L"`
	Rect      pixel.Rect       `json:"-"`
	Type      ObjectType       `json:"T,ObjectType"`
	P         ObjectProperties `json:",omitempty"`
	SpriteNum int              `json:"S,omitempty"`
	Sprite    *pixel.Sprite    `json:"-"`
	w         *World           `json:"-"`
}

func (o Object) String() string {
	return fmt.Sprintf("%s %s %s %v", o.Loc, o.Rect, o.Type, o.SpriteNum)
}

type ObjectProperties struct {
	Invisible bool `json:",omitempty"`
	//	Tile      bool `json:",omitempty"`
	//	Block     bool `json:",omitempty"`
	Special bool `json:",omitempty"`
}

func NewTile(loc pixel.Vec) Object {
	return Object{
		Loc:  loc,
		Rect: pixel.Rect{loc.Sub(pixel.V(16, 16)), loc.Add(pixel.V(16, 16))},
		Type: O_TILE,
	}
}
func NewBlock(loc pixel.Vec) Object {
	return Object{
		Loc:  loc,
		Rect: pixel.Rect{loc.Sub(pixel.V(16, 16)), loc.Add(pixel.V(16, 16))},
		Type: O_BLOCK,
	}
}

func NewTileBox(rect pixel.Rect) Object {
	return Object{
		Rect: rect,
		Type: O_TILE,
	}
}
func NewBlockBox(rect pixel.Rect) Object {
	return Object{
		Rect: rect,
		Type: O_BLOCK,
	}
}
func (o Object) Highlight(win pixel.Target) {
	imd := imdraw.New(nil)
	color := pixel.ToRGBA(colornames.Red)
	if o.Type == O_TILE {
		color = pixel.ToRGBA(colornames.Blue)
	}
	imd.Color = color.Scaled(0.3)
	imd.Push(o.Rect.Min, o.Rect.Max)
	imd.Rectangle(4)
	imd.Draw(win)
}
func (o Object) Draw(win pixel.Target, spritesheet pixel.Picture, sheetFrames []*pixel.Sprite) {
	if o.Type != O_BLOCK && o.Type != O_TILE {
		log.Println("UNKNOWN TILE", o)
	}
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
	total := len(things)
	for i, t := range things {
		t.w = w
		t.Rect = DefaultSpriteRectangle.Moved(t.Loc)
		switch t.SpriteNum {
		case 53:
			t.Type = O_BLOCK
		default:
		}

		switch t.Type {
		case O_BLOCK:
			//log.Printf("%v/%v block object: %s %v %s", i, total, t.Loc, t.SpriteNum, t.Type)
			w.Blocks = append(w.Blocks, t)
		case O_TILE:
			//log.Printf("%v/%v tile object: %s %v %s", i, total, t.Loc, t.SpriteNum, t.Type)
			w.Tiles = append(w.Tiles, t)

		default: //
			log.Printf("%v/%v skipping bad object: %s %v %s", i, total, t.Loc, t.SpriteNum, t.Type)
		}
	}
	log.Printf("map has %v blocks, %v tiles", len(w.Blocks), len(w.Tiles))
	//log.Println(w.Blocks, w.Tiles)
	return
}

/*
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
*/
// assumes only tiles are given
func FindRandomTile(os []Object) pixel.Vec {
	if len(os) == 0 {
		panic("no objects")
	}
	return os[rand.Intn(len(os))].Rect.Center()
}

func GetObjects(objects []Object, position pixel.Vec) []Object {
	var good []Object
	for _, o := range objects {
		if o.Rect.Contains(position) {
			good = append(good, o)
		}
	}
	return good
}

func GetTiles(objects []Object) []Object {
	var tiles []Object
	for _, o := range objects {
		if o.Type == O_TILE {
			tiles = append(tiles, o)
		}
	}
	return tiles
}

func GetTilesAt(objects []Object, position pixel.Vec) []Object {
	var good []Object
	all := GetObjects(objects, position)
	if len(all) > 0 {
		for _, o := range all {
			if o.Type == O_TILE {
				good = append(good, o)
			}

		}
	}
	return good

}
func GetBlocks(objects []Object, position pixel.Vec) []Object {
	var bad []Object
	all := GetObjects(objects, position)
	if len(all) > 0 {
		for _, o := range all {
			if o.Type == O_BLOCK {
				bad = append(bad, o)
			}

		}
	}
	return bad
}

// GetNeighbors gets the neighboring tiles
func (o Object) GetNeighbors() []Object {
	neighbors := []Object{}
	of := 32.0
	for _, offset := range [][]float64{
		{-of, 0},
		{of, 0},
		{0, -of},
		{0, of},
	} {
		if n := o.w.Tile(pixel.V(o.Rect.Center().X+offset[0], o.Rect.Center().Y+offset[1])); n.Type == o.Type {
			neighbors = append(neighbors, n)
		}
	}
	return neighbors

}
