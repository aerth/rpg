package rpg

import (
	"log"
	"math/rand"
	"time"

	"golang.org/x/image/colornames"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
)

type Animation struct {
	Name      string
	Type      ActionType
	loc       pixel.Vec
	rect      pixel.Rect
	radius    float64
	step      float64
	counter   float64
	cols      []pixel.RGBA
	until     time.Time
	start     time.Time
	damage    float64
	direction Direction
	ticker    <-chan time.Time
	level     uint
}
type ActionType int

const (
	Talk ActionType = iota
	Slash
	ManaStorm
	MagicBullet
	Arrow
)

func (w *World) Action(char *Character, loc pixel.Vec, t ActionType, dir Direction) {
	switch t {
	case Talk:
		log.Println("nothing to say yet")
	case Slash:
		log.Println("no weapon yet")
	case ManaStorm:
		cost := uint(2.5 * float64(char.Level))
		if char.Mana < cost {
			w.Message("not enough mana")
			return
		}
		char.Mana -= cost
		w.NewAnimation(char.Rect.Center(), "manastorm", OUT)
	case MagicBullet:
		cost := uint(1)
		if char.Mana < cost {
			w.Message("not enough mana")
			return
		}

		char.Mana -= cost
		w.NewAnimation(char.Rect.Center(), "magicbullet", dir)
	default: //
	}
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func (a *Animation) update(dt float64) {
	if time.Since(a.until) > time.Millisecond {
		a = nil
		return
	}
	switch a.Type {
	default:
		log.Println("nil animation?")
		return
	case MagicBullet:
		a.counter += dt
		for a.counter > a.step {
			if len(a.cols) == 0 {
				a.cols = []pixel.RGBA{RandomColor(), RandomColor(), RandomColor(), RandomColor(), RandomColor()}
			}
			a.counter -= a.step
			for i := len(a.cols) - 2; i >= 0; i-- {
				a.cols[i+1] = a.cols[i]
			}
			a.cols[0] = RandomColor().Scaled(0.3)
			a.cols[1] = RandomColor().Scaled(0.3)
		}
		if a.direction != OUT && a.direction != IN {
			a.loc = a.loc.Add((a.direction.V().Scaled(100 * dt)))
			a.rect = pixel.R(-a.radius, -a.radius, a.radius, a.radius).Moved(a.loc)
		}
	case ManaStorm:
		a.counter += dt
		for a.counter > a.step {
			if len(a.cols) == 0 {
				a.cols = []pixel.RGBA{RandomColor(), RandomColor(), RandomColor(), RandomColor(), RandomColor()}
			}
			a.counter -= a.step
			for i := len(a.cols) - 2; i >= 0; i-- {
				a.cols[i+1] = a.cols[i]
			}
			a.cols[0] = RandomColor().Scaled(0.3)
			a.cols[1] = RandomColor().Scaled(0.3)
		}
	case Arrow:
		a.counter += dt
		for a.counter > a.step {
			a.counter -= a.step
			a.direction = Direction(rand.Intn(5))
		}
		a.loc = a.loc.Add((a.direction.V().Scaled(100 * dt)))
		a.rect = pixel.R(-100, -5, 100, 5).Moved(a.loc)

	}
}

func (a *Animation) draw(imd *imdraw.IMDraw) {
	if a == nil || time.Since(a.start) < time.Millisecond {
		return
	}
	switch a.Type {
	case MagicBullet, ManaStorm:
		for i := len(a.cols) - 1; i >= 0; i-- {
			imd.Color = a.cols[i]
			imd.Push(a.loc)
			imd.Circle(float64(i+2)*a.radius/float64(len(a.cols)), 0)
		}
	case Arrow:
		imd.Color = RandomColor()
		b := a.loc.Add(a.direction.V().Scaled(10))
		imd.Push(a.loc, b)
		imd.Line(float64(a.level))

	default:
		log.Println("bad animation?")
	}
}
func (w *World) NewAnimation(loc pixel.Vec, kind string, direction Direction) {
	switch kind {
	default: //
		log.Println("invalid animation type")
		return
	case "magicbullet":
		a := new(Animation)
		a.Type = MagicBullet
		a.loc = loc
		a.radius = 10
		a.step = 1.0 / 7
		a.rect = pixel.R(-a.radius, -a.radius, a.radius, a.radius).Moved(a.loc)
		a.cols = []pixel.RGBA{}
		a.start = time.Now()
		a.until = time.Now().Add(time.Second * 2)
		a.direction = direction
		a.damage = w.Char.Stats.Intelligence * 0.5
		w.Animations = append(w.Animations, a)

	case "manastorm":
		a := new(Animation)
		a.Type = ManaStorm
		a.loc = loc
		a.radius = (w.Char.Stats.Intelligence / 20) * 4
		a.step = 1.0 / 7
		a.damage = w.Char.Stats.Intelligence * 0.2
		a.rect = pixel.R(-a.radius, -a.radius, a.radius, a.radius).Moved(a.loc)
		a.cols = []pixel.RGBA{}
		a.start = time.Now()
		a.direction = direction
		dur := time.Duration(w.Char.Stats.Intelligence * float64(time.Millisecond) * 18)
		//	log.Println(dur)
		a.until = time.Now().Add(dur)
		a.ticker = time.Tick(time.Millisecond * 200)
		w.Animations = append(w.Animations, a)

	case "arrow":
		a := new(Animation)
		a.Type = Arrow
		a.loc = loc
		a.damage = w.Char.Stats.Strength * 10
		a.rect = pixel.R(-100, -5, 100, 5)
		a.cols = []pixel.RGBA{pixel.ToRGBA(colornames.Red)}
		a.direction = direction
		dur := time.Second * 3
		a.until = time.Now().Add(dur)
		a.ticker = time.Tick(time.Millisecond * 500)
		a.level = w.Char.Level
		w.Animations = append(w.Animations, a)

	}

}
