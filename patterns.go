package rpg

import (
	"log"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

type Animation struct {
	Name    string
	Type    ActionType
	loc     pixel.Vec
	rect    pixel.Rect
	radius  float64
	step    float64
	counter float64
	cols    [5]pixel.RGBA
	until   time.Time
	start   time.Time
	damage  float64
}

func (a *Animation) update(dt float64) {
	if time.Since(a.until) > time.Second {
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
}

func (a *Animation) draw(imd *imdraw.IMDraw) {
	if a == nil || time.Since(a.start) > time.Millisecond {
		return
	}

	for i := len(a.cols) - 1; i >= 0; i-- {
		imd.Color = a.cols[i]
		imd.Push(a.loc)
		imd.Circle(float64(i+1)*a.radius/float64(len(a.cols)), 0)
	}
}
func (w *World) NewAnimation(loc pixel.Vec, kind string) {
	switch kind {
	default: //
		log.Println("invalid animation type")
		return
	case "magic":
		a := new(Animation)
		a.loc = loc
		a.radius = 180
		a.step = 1.0 / 7
		a.damage = w.Char.Stats.Intelligence * 1.3
		a.rect = pixel.R(-180, -180, 180, 180).Moved(a.loc)
		a.cols = [5]pixel.RGBA{}
		a.start = time.Now().Add(1 * time.Second)
		a.until = time.Now().Add(3 * time.Second)
		w.Animations = append(w.Animations, a)
		for i, v := range w.Entities {
			if a.rect.Contains(v.Rect.Center()) {
				w.Entities[i].P.Health -= a.damage
				if w.Entities[i].P.Health <= 0 {
					// damage func should be function
					if w.Entities[i].P.IsDead {
						w.Entities[i].P.Health = 0
						continue
					}
					w.Entities[i].P.Health = 0
					w.Entities[i].P.IsDead = true
					w.Char.Stats.Kills++

					log.Println("Got new loot!:", FormatItemList(v.P.Loot))
					w.Message(" Loot: " + FormatItemList(v.P.Loot))

					w.Char.Inventory = StackItems(w.Char.Inventory, v.P.Loot)
					log.Println("New inventory:", w.Char.Inventory)
					w.Char.ExpUp(1)
					w.checkLevel()

				}

				log.Printf("%s took %v damage, now at %v HP",
					w.Entities[i].Name, a.damage, w.Entities[i].P.Health)
			}
		}

	}

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
