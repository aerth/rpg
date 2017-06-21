package rpg

import (
	"image/color"
	"log"
	"math/rand"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type Animation struct {
	Name      string
	Type      ActionType
	loc       pixel.Vec
	rect      pixel.Rect
	radius    float64
	step      float64
	counter   float64
	cols      [5]pixel.RGBA
	until     time.Time
	start     time.Time
	damage    float64
	direction Direction
}

func (a *Animation) update(dt float64) {
	if time.Since(a.until) > time.Millisecond {
		a = nil
		return
	}
	a.counter += dt
	for a.counter > a.step {
		a.counter -= a.step
		for i := len(a.cols) - 2; i >= 0; i-- {
			a.cols[i+1] = a.cols[i]
		}
		a.cols[0] = RandomColor().Scaled(0.3)
		a.cols[1] = RandomColor().Scaled(0.3)
	}
	if a.direction != OUT {
		a.loc = a.loc.Add(a.direction.V().Scaled(100 * dt))
		a.rect = pixel.R(-a.radius, -a.radius, a.radius, a.radius).Moved(a.loc)
	}
}

func (a *Animation) draw(imd *imdraw.IMDraw) {
	if a == nil || time.Since(a.start) < time.Millisecond {
		return
	}

	for i := len(a.cols) - 1; i >= 0; i-- {
		imd.Color = a.cols[i]
		imd.Push(a.loc)
		imd.Circle(float64(i+2)*a.radius/float64(len(a.cols)), 0)
	}
}
func (w *World) NewAnimation(loc pixel.Vec, kind string, direction Direction) {
	switch kind {
	default: //
		log.Println("invalid animation type")
		return
	case "magicbullet":
		a := new(Animation)
		a.loc = loc
		a.radius = 5
		a.step = 1.0 / 7
		a.rect = pixel.R(-a.radius, -a.radius, a.radius, a.radius).Moved(a.loc)
		a.cols = [5]pixel.RGBA{}
		a.start = time.Now()
		a.until = time.Now().Add(time.Second * 2)
		a.direction = direction
		a.damage = w.Char.Stats.Intelligence * 1.3
		w.Animations = append(w.Animations, a)

	case "manastorm":
		a := new(Animation)
		a.loc = loc
		a.radius = 140 * w.Char.Stats.Intelligence / 100
		a.step = 1.0 / 7
		a.damage = w.Char.Stats.Intelligence * 1.3
		a.rect = pixel.R(-a.radius, -a.radius, a.radius, a.radius).Moved(a.loc)
		a.cols = [5]pixel.RGBA{}
		a.start = time.Now()
		a.direction = direction
		a.until = time.Now().Add(
			time.Duration(
				w.Char.Stats.Intelligence * float64(time.Millisecond) * 10))
		w.Animations = append(w.Animations, a)
	}

}

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
	imd.Color = color
	percent := cur / max
	one := rect.Min
	one.Y++
	imd.Push(rect.Min, rect.Max)
	imd.Rectangle(1)
	pt := pixel.V(rect.Max.X*percent, rect.Max.Y)
	imd.Push(rect.Min, pt)
	imd.Rectangle(0)

}
