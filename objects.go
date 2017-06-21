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

//var DefaultSpriteRectangle = pixel.R(-16, 0, 16, 32)
//var DefaultSpriteRectangle = pixel.R(-16, 0, 16, 32)

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
	W         *World           `json:"-"`
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

var TransparentBlue = pixel.ToRGBA(colornames.Blue).Scaled(0.8)
var TransparentRed = pixel.ToRGBA(colornames.Red).Scaled(0.8)
var TransparentPurple = pixel.ToRGBA(colornames.Purple).Scaled(0.8)

func (o Object) Highlight(win pixel.Target, color pixel.RGBA) {
	imd := imdraw.New(nil)
	imd.Push(o.Rect.Min, o.Rect.Max)
	imd.Rectangle(2)
	imd.Draw(win)
}
func (o Object) Draw(win pixel.Target, spritesheet pixel.Picture, sheetFrames []*pixel.Sprite) {
	//	r := pixel.Rect{o.Loc, o.w.Char.Rect.Center()}
	//	sz := r.Size()
	//	if sz.X > 1000 || sz.Y > 1000 {
	//		return
	//	}
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
	//	if o.Loc == pixel.ZV && o.Rect.Size().Y != 32 {
	//		log.Println(o.Rect.Size(), "cool rectangle", o.SpriteNum)
	//		DrawPattern(win, o.Sprite, o.Rect, 0)
	//	} else {
	o.Sprite.Draw(win, pixel.IM.Moved(o.Loc))
	//	}

}
func (w *World) LoadMapFile(path string) error {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	return w.loadmap(b)
}
func (w *World) LoadMap(path string) error {
	b, err := assets.Asset(path)
	if err != nil {
		return err
	}
	return w.loadmap(b)
}
func (w *World) loadmap(b []byte) error {
	var things = []Object{}
	err := json.Unmarshal(b, &things)
	if err != nil {
		return fmt.Errorf("invalid map: %v", err)
	}
	total := len(things)
	for i, t := range things {
		t.W = w
		t.Rect = DefaultSpriteRectangle.Moved(t.Loc)
		switch t.SpriteNum {
		case 53: // water
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
	if len(w.Blocks) == 0 && len(w.Tiles) == 0 {
		return fmt.Errorf("invalid map")
	}
	return nil
}

// assumes only tiles are given
func FindRandomTile(os []Object) pixel.Vec {
	if len(os) == 0 {
		panic("no objects")
	}
	ob := os[rand.Intn(len(os))]
	if ob.Loc != pixel.ZV && ob.SpriteNum != 0 && ob.Type == O_TILE {
		return ob.Rect.Center()
	}
	return FindRandomTile(os)
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

func TilesAt(objects []Object, position pixel.Vec) []Object {
	var good []Object
	all := GetObjects(objects, position)
	if len(all) > 0 {
		for _, o := range all {
			if DefaultSpriteRectangle.Moved(o.Loc).Contains(position) && o.Type == O_TILE {
				good = append(good, o)
			}

		}
	}
	return good

}
func GetObjectsAt(objects []Object, position pixel.Vec) []Object {
	var good []Object
	all := GetObjects(objects, position)
	if len(all) > 0 {
		for _, o := range all {
			if DefaultSpriteRectangle.Moved(o.Loc).Contains(position) {
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
		if n := o.W.Tile(pixel.V(o.Rect.Center().X+offset[0], o.Rect.Center().Y+offset[1])); n.Type == o.Type {
			neighbors = append(neighbors, n)
		}
	}
	return neighbors

}
func (w *World) drawTiles(path string) error {
	spritesheet, spritemap := LoadSpriteSheet(path)
	// layers (TODO: slice?)
	// batch sprite drawing
	globebatch := pixel.NewBatch(&pixel.TrianglesData{}, spritesheet)
	// water world 67 wood, 114 117 182 special, 121 135 dirt, 128 blank, 20 grass
	//      rpg.DrawPattern(batch, spritemap[53], pixel.R(-3000, -3000, 3000, 3000), 100)

	globebatch.Clear()
	// draw it on to canvasglobe
	for _, o := range w.Tiles {
		o.Draw(globebatch, spritesheet, spritemap)
	}
	for _, o := range w.Blocks {
		o.Draw(globebatch, spritesheet, spritemap)
	}
	w.Batches[EntityType(-1)] = globebatch
	return nil
}

func TileNear(all []Object, loc pixel.Vec) Object {
	tile := TilesAt(all, loc)
	radius := 1.00
	oloc := loc
	if len(tile) > 0 {
		oloc = tile[0].Loc
	}
	log.Println("looking for loc:", loc)
	for i := 0; i < len(all); i++ {
		if loc == oloc {
			loc.X += 16
			continue
		}
		log.Println("Checking loc", loc)
		os := TilesAt(all, loc)
		if len(os) > 0 {
			if os[0].Loc == pixel.ZV || os[0].Loc == oloc {
				continue
			}
			return os[0]
		}
		os = TilesAt(all, loc.Scaled(-2))
		if len(os) > 0 {
			if os[0].Loc == pixel.ZV || os[0].Loc == oloc {
				continue
			}
			return os[0]
		}

		loc.X += radius * 16
		loc.Y += radius * 16
		if i%4 == 1 {
			radius++
		}
	}
	return Object{Type: O_NONE}

}
