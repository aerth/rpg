package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"runtime/pprof"
	"time"

	"golang.org/x/image/colornames"

	//	_ "image/png"

	"github.com/aerth/rpg"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var (
	flagenemies = flag.Int("e", 2, "number of enemies to begin with")
	flaglevel   = flag.String("test", "1", "custom world test (filename)")

	debug = flag.Bool("v", false, "extra logs")
)

const (
	LEFT  = rpg.LEFT
	RIGHT = rpg.RIGHT
	UP    = rpg.UP
	DOWN  = rpg.DOWN
)

var (
	ZV = pixel.ZV
	IM = pixel.IM
	V  = pixel.V
	R  = pixel.R
)

var (
	defaultzoom  = 3.0
	camZoomSpeed = 1.20
)

func run() {
	f, err := os.Create("p.debug")
	if err != nil {
		log.Fatal(err)
	}
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

	rand.Seed(time.Now().UnixNano())

	winbounds := pixel.R(0, 0, 800, 600)

	fmt.Println("which screen resolution?")
	fmt.Println("1. 800x600")
	fmt.Println("2. 1024x768")
	fmt.Println("3. 1280x720")
	var screenres int
	_, err = fmt.Scanf("%d", &screenres)
	if err != nil {
		screenres = 0
	}
	switch screenres {
	default:
	case 2:
		winbounds = pixel.R(0, 0, 1024, 768)
	case 3:
		winbounds = pixel.R(0, 0, 1280, 720)
	}
	// window options
	cfg := pixelgl.WindowConfig{
		Title:  rpg.Version(),
		Bounds: winbounds,
		//Undecorated: true,
		VSync: false,
	}

	// create window
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}
	win.SetSmooth(false)
	buttons := []rpg.Button{
		{Name: "manastorm", Frame: pixel.R(10, 10, 42, 42)},
		{Name: "magicbullet", Frame: pixel.R(42, 10, 42+42, 42)},
	}
	// START
	//world.Char.Rect = world.Char.Rect.Moved(V(33, 33))
	// load world
	//	worldbounds = pixel.R(float64(-4000), float64(-4000), float64(4000), float64(4000))

	world := rpg.NewWorld(*flaglevel, *flagenemies)
	if world == nil {
		return
	}
	// sprite sheet
	spritesheet, spritemap := rpg.LoadSpriteSheet("tileset.png")

	// layers (TODO: slice?)
	// batch sprite drawing
	globebatch := pixel.NewBatch(&pixel.TrianglesData{}, spritesheet)
	animbatch := pixel.NewBatch(&pixel.TrianglesData{}, spritesheet)

	// water world 67 wood, 114 117 182 special, 121 135 dirt, 128 blank, 20 grass
	//	rpg.DrawPattern(batch, spritemap[53], pixel.R(-3000, -3000, 3000, 3000), 100)

	// draw menu bar
	menubatch := pixel.NewBatch(&pixel.TrianglesData{}, spritesheet)
	rpg.DrawPattern(menubatch, spritemap[67], pixel.R(0, 0, win.Bounds().W()+20, 60), 100)
	for _, btn := range buttons {
		spritemap[200].Draw(menubatch, IM.Moved(btn.Frame.Center()))
	}

	redrawWorld := func(w *rpg.World) {
		globebatch.Clear()
		// draw it on to canvasglobe
		for _, v := range w.Tiles {
			v.Draw(globebatch, spritesheet, spritemap)
		}
		for _, v := range w.Blocks {
			v.Draw(globebatch, spritesheet, spritemap)
		}

	}

	// create NPC

	world.NewMobs(*flagenemies)
	l := time.Now()
	var last = &l
	second := time.Tick(time.Second)
	tick := time.Tick(time.Second)

	frames := 0
	var delda float64 = 0.00
	var camZoom = &defaultzoom
	var dt *float64 = &delda
	t1 := time.Now()
	text := rpg.NewText(36)
	// start loop
	imd := imdraw.New(nil)
	rand.Seed(time.Now().UnixNano())
	//var latest string
	redrawWorld(world)
MainLoop:
	for !win.Closed() {
		world.Reset()
		rpg.TitleMenu(win)
		world.Char.Health = 255
		world.Char.Rect = rpg.DefaultPhys.Rect.Moved(rpg.FindRandomTile(world.Tiles))
	GameLoop:
		for !win.Closed() {

			if world.Char.Health < 1 {
				log.Println("GAME OVER")
				break GameLoop
			}
			*dt = time.Since(*last).Seconds()
			*last = time.Now()
			// zoom with mouse scroll

			*camZoom *= math.Pow(camZoomSpeed, win.MouseScroll().Y)
			if *camZoom > 6.5 {
				*camZoom = 6.5
			}
			if *camZoom < 2 {
				*camZoom = 2
			}
			// drawing
			//win.Clear(rpg.RandomColor())
			win.Clear(colornames.Black)
			animbatch.Clear()
			// if key
			if win.JustPressed(pixelgl.KeyQ) && win.Pressed(pixelgl.KeyLeftControl) {
				break MainLoop
			}
			// teleport random
			if win.JustPressed(pixelgl.Key8) {
				world.Char.Rect = rpg.DefaultSpriteRectangle.Moved(rpg.TileNear(world.Tiles, world.Char.Rect.Center()).Loc)

			}

			// move all enemies (debug)
			if win.JustPressed(pixelgl.Key9) {
				for _, v := range world.Entities {
					v.Rect = rpg.DefaultEntityRectangle.Moved(rpg.TileNear(world.Tiles, v.Rect.Center()).Loc)
				}
			}
			if win.JustReleased(pixelgl.KeyI) {
				rpg.InventoryLoop(win, world)
			}

			if win.JustPressed(pixelgl.KeyEqual) {
				*debug = !*debug
			}

			dir := controlswitch(dt, world, win, buttons, win)
			world.Char.Update(*dt, dir, world)
			world.Update(*dt)
			world.Clean()
			cam := pixel.IM.Scaled(pixel.ZV, *camZoom).Moved(win.Bounds().Center()).Moved(world.Char.Rect.Center().Scaled(-*camZoom))
			win.SetMatrix(cam)

			// draw map (tiles and blocks) (never updated for now)
			globebatch.Draw(win)

			if *debug {
				world.HighlightPaths(win)
			}
			// draw entities and objects (not tiles and blocks)
			world.Draw(win)

			// animations such as magic spells
			imd.Clear()
			world.CleanAnimations()
			world.ShowAnimations(imd)
			imd.Draw(win)

			if *debug {
				for _, o := range world.Tile(world.Char.Rect.Center()).PathNeighbors() {
					ob := o.(rpg.Object)
					ob.W = world
					ob.Highlight(win, rpg.TransparentPurple)
				}
			}

			// back to window cam
			win.SetMatrix(pixel.IM)
			world.Char.Matrix = pixel.IM.Scaled(pixel.ZV, *camZoom).Scaled(pixel.ZV, 0.5).Moved(pixel.V(0, 0)).Moved(win.Bounds().Center())
			world.Char.Draw(win)
			// draw score board
			text.Clear()
			rpg.DrawScore(winbounds, text, win,
				"[%vHP·%vMP·%sGP] [LVL%v] [%v/%vXP] [%vKills]", world.Char.Health, world.Char.Mana, world.Char.CountGold(), world.Char.Level, world.Char.Stats.XP, world.Char.NextLevel(), world.Char.Stats.Kills)

			// draw menubar
			menubatch.Draw(win)
			if win.JustPressed(pixelgl.Key6) {
				redrawWorld(world)
			}

			// draw health, mana, xp bars
			world.Char.DrawBars(win, win.Bounds())

			if win.JustPressed(pixelgl.MouseButtonLeft) {
				mouseloc := win.MousePosition()
				if b, f, ok := world.IsButton(buttons, mouseloc); ok {
					log.Println(mouseloc)
					log.Printf("Clicked button: %q", b.Name)
					f(win, world)

				} else {

					mcam := pixel.IM.Moved(win.Bounds().Center())
					mouseloc1 := mcam.Unproject(mouseloc)
					unit := mouseloc1.Unit()
					//				log.Println("unit:", unit)
					dir := rpg.UnitToDirection(unit)
					//				log.Println("direction:", dir)
					if dir == rpg.OUT || dir == rpg.IN {
						dir = world.Char.Dir
					}
					if world.Char.Mana > 0 {
						world.Action(world.Char, world.Char.Rect.Center(), rpg.MagicBullet, dir)
					} else {
						log.Println("Not enough mana")
					}
				}
			}
			select {
			default: //
			case <-tick:
			}

			//spritemap[20].Draw(menubar, pixel.IM.Scaled(ZV, 10).Moved(pixel.V(30, 30)))
			//menubar.Draw(win, pixel.IM)
			win.Update()

			// fps, gps
			frames++
			gps := world.Char.Rect.Center()
			select {
			default: //keep going
			case <-second:
				str := fmt.Sprintf(""+
					"FPS: %d | GPS: (%v,%v) | VEL: (%v) | HP: (%v) ",
					frames, int(gps.X), int(gps.Y), int(world.Char.Phys.Vel.Len()), world.Char.Health)
				win.SetTitle(str)
				frames = 0
			}

		}
		log.Printf("You survived for %s.\nYou acquired %s gold", time.Now().Sub(t1), world.Char.CountGold())
		log.Println("Inventory:", rpg.FormatItemList(world.Char.Inventory))
		log.Printf("Skeletons killed: %v", world.Char.Stats.Kills)
	}

}

func controlswitch(dt *float64, w *rpg.World, win *pixelgl.Window, buttons []rpg.Button, buf pixel.Target) rpg.Direction {
	if win.JustPressed(pixelgl.KeySpace) || win.JustPressed(pixelgl.MouseButtonRight) {
		if w.Char.Mana > 0 {
			w.Action(w.Char, w.Char.Rect.Center(), rpg.ManaStorm, w.Char.Dir)
		} else {
			log.Println("Not enough mana")
		}
	}
	if win.JustPressed(pixelgl.KeyB) {
		if w.Char.Mana > 0 {
			w.Action(w.Char, w.Char.Rect.Center(), rpg.MagicBullet, w.Char.Dir)
		} else {
			log.Println("Not enough mana")
		}
	}

	// slow motion with tab
	if win.Pressed(pixelgl.KeyTab) {
		*dt /= 8
	}
	// speed motion with tab
	if win.Pressed(pixelgl.KeyLeftShift) {
		*dt *= 4
	}
	if win.Pressed(pixelgl.Key1) {
		w.Char.Mana += 1
		if w.Char.Mana > 255 {
			w.Char.Mana = 255
		}
	}
	if win.Pressed(pixelgl.Key2) {
		w.Char.Health += 1
		if w.Char.Health > 255 {
			w.Char.Health = 255
		}
	}

	if win.Pressed(pixelgl.Key3) {
		w.Char.Stats.XP += 10
	}

	if win.Pressed(pixelgl.KeyCapsLock) {
		w.Char.Phys.CanFly = !w.Char.Phys.CanFly
	}
	dir := w.Char.Dir

	/*if win.JustPressed(pixelgl.MouseButtonLeft) {
		mouseloc := win.MousePosition()
		if b, f, ok := w.IsButton(buttons, mouseloc); ok {
			log.Println(mouseloc)
			log.Printf("Clicked button: %q", b.Name)
			f(win, w)

		}
	} */

	if win.Pressed(pixelgl.KeyLeft) || win.Pressed(pixelgl.KeyH) || win.Pressed(pixelgl.KeyA) {
		dir = LEFT
		w.Char.Phys.Vel.X = -w.Char.Phys.RunSpeed * (1 + *dt)
	}
	if win.Pressed(pixelgl.KeyRight) || win.Pressed(pixelgl.KeyL) || win.Pressed(pixelgl.KeyD) {
		w.Char.Phys.Vel.X = +w.Char.Phys.RunSpeed * (1 + *dt)
		dir = RIGHT
	}
	if win.Pressed(pixelgl.KeyDown) || win.Pressed(pixelgl.KeyJ) || win.Pressed(pixelgl.KeyS) {
		w.Char.Phys.Vel.Y = -w.Char.Phys.RunSpeed * (1 + *dt)
		dir = DOWN

	}
	if win.Pressed(pixelgl.KeyUp) || win.Pressed(pixelgl.KeyK) || win.Pressed(pixelgl.KeyW) {
		w.Char.Phys.Vel.Y = +w.Char.Phys.RunSpeed * (1 + *dt)
		dir = UP
	}

	if win.Pressed(pixelgl.KeyUp) && win.Pressed(pixelgl.KeyLeft) {
		dir = rpg.UPLEFT
	}
	if win.Pressed(pixelgl.KeyUp) && win.Pressed(pixelgl.KeyRight) {
		dir = rpg.UPRIGHT
	}
	if win.Pressed(pixelgl.KeyDown) && win.Pressed(pixelgl.KeyLeft) {
		dir = rpg.DOWNLEFT
	}
	if win.Pressed(pixelgl.KeyDown) && win.Pressed(pixelgl.KeyRight) {
		dir = rpg.DOWNRIGHT
	}
	// restart the level on pressing enter
	//	if win.JustPressed(pixelgl.KeyEnter) {
	//		rpg.InventoryLoop(win, w)
	//	}
	return dir
}
func main() {
	flag.Parse()
	pixelgl.Run(run)
}
