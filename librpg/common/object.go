package common

import (
	"fmt"
	"log"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"golang.org/x/image/colornames"
)

var DefaultSpriteRectangle = pixel.R(-16, -16, 16, 16)

type ObjectType int

const (
	O_NONE ObjectType = iota
	O_TILE
	O_BLOCK
	O_INVISIBLE
	O_SPECIAL
	O_WIN
	O_DYNAMIC // loot, doors
)

type Object struct {
	Loc       pixel.Vec        `json:"L"`
	Rect      pixel.Rect       `json:"-"`
	Type      ObjectType       `json:"T"`
	P         ObjectProperties `json:",omitempty"`
	SpriteNum int              `json:"S,omitempty"`
	Sprite    *pixel.Sprite    `json:"-"`
	W         interface {
		Tile(dot pixel.Vec) Object
	} `json:"-"`
	//	Contains  []Item           `json:"-"`
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

var TransparentBlue = pixel.ToRGBA(colornames.Blue).Scaled(0.4)
var TransparentRed = pixel.ToRGBA(colornames.Red).Scaled(0.4)
var TransparentPurple = pixel.ToRGBA(colornames.Purple).Scaled(0.4)

func (o Object) Highlight(win pixel.Target, color pixel.RGBA) {
	imd := imdraw.New(nil)
	imd.Color = color
	imd.Push(o.Rect.Min, o.Rect.Max)
	imd.Rectangle(1)
	imd.Draw(win)
}
func (o Object) Draw(win pixel.Target, spritesheet pixel.Picture, sheetFrames []*pixel.Sprite, scaled float64) {
	//	r := pixel.Rect{o.Loc, o.w.Char.Rect.Center()}
	//	sz := r.Size()
	//	if sz.X > 1000 || sz.Y > 1000 {
	//		return
	//	}
	if o.P.Invisible {
		return
	}
	if o.Sprite == nil && o.Type != O_DYNAMIC {
		if 0 > o.SpriteNum || o.SpriteNum > len(sheetFrames) {
			log.Printf("unloadable sprite: %v/%v", o.SpriteNum, len(sheetFrames))
			return
		}
		o.Sprite = sheetFrames[o.SpriteNum]
	}
	if o.Sprite == nil && o.Type == O_DYNAMIC {
		o.Sprite = sheetFrames[0]
	}
	//	if o.Loc == pixel.ZV && o.Rect.Size().Y != 32 {
	//		log.Println(o.Rect.Size(), "cool rectangle", o.SpriteNum)
	//		DrawPattern(win, o.Sprite, o.Rect, 0)
	//	} else {
	if scaled != 0.00 {
		o.Sprite.Draw(win, pixel.IM.Scaled(pixel.ZV, scaled).Moved(o.Loc))
		return
	}

	o.Sprite.Draw(win, pixel.IM.Moved(o.Loc))
	//	}

}

const _ObjectType_name = "O_NONEO_TILEO_BLOCKO_INVISIBLEO_SPECIALO_WINO_DYNAMIC"

var _ObjectType_index = [...]uint8{0, 6, 12, 19, 30, 39, 44, 53}

func (i ObjectType) String() string {
	if i < 0 || i >= ObjectType(len(_ObjectType_index)-1) {
		return fmt.Sprintf("ObjectType(%d)", i)
	}
	return _ObjectType_name[_ObjectType_index[i]:_ObjectType_index[i+1]]
}
