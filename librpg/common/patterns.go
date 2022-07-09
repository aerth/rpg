package common

import (
	"image/color"
	"log"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

func DrawPatternObject(spritenum int, objecttype ObjectType, bounds pixel.Rect, width float64) []Object {
	var objects []Object
	size := pixel.Rect{pixel.V(-16, -16), pixel.V(16, 16)}
	//size := DefaultSpriteRectangle
	for y := bounds.Min.Y; y < bounds.Max.Y; y = y + size.H() {
		for x := bounds.Min.X; x < bounds.Max.X; x = x + size.W() {
			o := Object{
				Loc:       pixel.V(x, y),
				Rect:      size.Moved(pixel.V(x, y)),
				Type:      objecttype,
				SpriteNum: spritenum,
			}
			objects = append(objects, o)
		}
	}
	return objects
}

func DrawPattern(canvas pixel.Target, sprite *pixel.Sprite, bounds pixel.Rect, width float64) {
	if bounds.Size() == pixel.ZV {
		return
	}
	var i int
	for y := bounds.Min.Y; y < bounds.Max.Y; y = y + sprite.Frame().H() {
		for x := bounds.Min.X; x < bounds.Max.X; x = x + sprite.Frame().W() {
			sprite.Draw(canvas, pixel.IM.Moved(pixel.V(x, y)))
			i++
		}
	}
	//log.Printf("Draw pattern: %v iterations", i)
}

func Drawbg(canvas *pixelgl.Canvas) {
	imd := imdraw.New(nil)
	var (
		bounds     = canvas.Bounds()
		squaresize = 2.00
	)
	for x := 0.00; x < bounds.W(); x = x + squaresize {
		batch := pixel.NewBatch(&pixel.TrianglesData{}, canvas)
		invert := int((x/squaresize))%2 ^ 1
		for y := 0.00; y < bounds.H(); y = y + squaresize {
			colored := float64(int((y/squaresize))%2 ^ invert)
			//xor := float64(int(colored) ^ invert))
			imd.Clear()
			imd.Color = RandomColor()
			imd.Push(pixel.V(x, y+squaresize))
			imd.Push(pixel.V(x, y+squaresize).Add(pixel.V(squaresize, squaresize)))
			imd.Rectangle(colored)
			imd.Draw(batch)

		}
		log.Printf("LOADING: %v", int(100*(x/bounds.W())))
		batch.Draw(canvas)
	}
}

func DrawBase(canvas *pixelgl.Canvas) {
	batch := pixel.NewBatch(&pixel.TrianglesData{}, canvas)
	imd := imdraw.New(nil)
	imd.Clear()
	imd.Color = colornames.Black
	imd.Push(pixel.V(100.00, 100.00))
	imd.Push(pixel.V(-100.00, -100.00))
	imd.Rectangle(0)

	imd.Color = RandomColor()
	imd.Push(pixel.V(100.00, 100.00))
	imd.Push(pixel.V(-100.00, -100.00))
	imd.Rectangle(20)
	imd.Color = RandomColor()
	imd.Push(pixel.V(100.00, 100.00))
	imd.Push(pixel.V(-100.00, -100.00))
	imd.Rectangle(10)
	imd.Draw(batch)
	batch.Draw(canvas)
}
func DrawBar(imd *imdraw.IMDraw, color color.RGBA, cur, max float64, rect pixel.Rect) {
	rect = rect.Norm()
	imd.Color = color
	if max < cur {
		cur, max = max, cur
	}
	percent := cur / max
	//	if rect.Max.Y-rect.Min.Y > 10 {
	//		rect.Max.Y = rect.Min.Y + 10
	//	}
	//one.Y++
	imd.Push(rect.Min, rect.Max)
	imd.Rectangle(1)
	pt := pixel.V(rect.Min.X+((rect.Max.X-rect.Min.X)*percent), rect.Max.Y)

	if pt.X < rect.Min.X {
		pt.X = rect.Min.X
	}
	imd.Push(rect.Min, pt)
	imd.Rectangle(0)

}
