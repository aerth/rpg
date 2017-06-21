// aerth game
// copyright 2017 aerth <aerth@riseup.net>
package rpg

import (
	"io/ioutil"
	"math"
	"math/rand"
	"os"
	"time"

	"golang.org/x/image/font"

	_ "image/png"

	"github.com/faiface/pixel"
	"github.com/golang/freetype/truetype"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// Direction LEFT RIGHT DOWN UP
type Direction int

type animState int

const (
	LEFT Direction = iota
	RIGHT
	DOWN
	UP
	IN
	OUT
	WEST  = LEFT
	EAST  = RIGHT
	NORTH = UP
	SOUTH = DOWN
)

const (
	Idle animState = iota
	Running
)

func (d Direction) String() string {
	switch d {
	case LEFT:
		return "left"
	case RIGHT:
		return "right"
	case UP:
		return "up"
	case DOWN:
		return "down"
	default:
		return "invalid direction"
	}
}

func (d Direction) V() pixel.Vec {
	switch d {
	case LEFT:
		return pixel.V(-1, 0)
	case RIGHT:
		return pixel.V(1, 0)
	case UP:
		return pixel.V(0, 1)
	case DOWN:
		return pixel.V(0, -1)
	default:
		return pixel.V(0, 0)
	}
}

func RandomColor() pixel.RGBA {

	r := rand.Float64()
	g := rand.Float64()
	b := rand.Float64()
	len := math.Sqrt(r*r + g*g + b*b)
	//if len == 0 {
	//	goto again
	//}
	return pixel.RGB(r/len, g/len, b/len)

}

func LoadTTF(path string, size float64) (font.Face, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	font, err := truetype.Parse(data)
	if err != nil {
		return nil, err
	}

	return truetype.NewFace(font, &truetype.Options{
		Size:              size,
		GlyphCacheEntries: 1,
	}), nil
}
func LoadSpriteSheet(path string) (pixel.Picture, []*pixel.Sprite) {
	spritesheet, err := LoadPicture("sprites/" + path)
	/* 0 1 2 3 4 5 6 7 8 9 10 11 12 13 14 15 16
	         * 1
		          * 2
			           * ...
				            * 16
	*/
	if err != nil {
		panic(err)
	}
	var sheetFrames []pixel.Rect
	for x := spritesheet.Bounds().Min.X; x < spritesheet.Bounds().Max.X; x += 32 {
		for y := spritesheet.Bounds().Min.Y; y < spritesheet.Bounds().Max.Y; y += 32 {
			sheetFrames = append(sheetFrames, pixel.R(x, y, x+32, y+32))
		}
	}
	var spritemap = []*pixel.Sprite{}
	for i := 0; i < len(sheetFrames); i++ {
		x := i
		spritemap = append(spritemap, pixel.NewSprite(spritesheet, sheetFrames[x]))
	}
	//log.Println(len(spritemap), "sprites loaded")
	return spritesheet, spritemap
}

// Distance between two vectors
func Distance(v1, v2 pixel.Vec) float64 {
	r := pixel.Rect{v1, v2}.Norm()
	v1 = r.Min
	v2 = r.Max
	h := (v1.X - v2.X) * (v1.X - v2.X)
	v := (v1.Y - v2.Y) * (v1.Y - v2.Y)
	return Sqrt(h + v)
}

// ?
func Sqrt(x float64) float64 {
	z := float64(2.)
	s := float64(0)
	for i := 0; i < 10; i++ {
		z = z - (z*z-x)/(2*z)
		if math.Abs(z-s) < 1e-10 {
			break
		}
		s = z
	}
	return z
}
