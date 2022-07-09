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

	rpg "github.com/aerth/rpc/librpg"
	"github.com/aerth/rpc/librpg/common"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var (
	flagenemies = flag.Int("e", 2, "number of enemies to begin with")
	flaglevel   = flag.String("test", "", "custom world test (filename)")
	flagseed    = flag.String("seed", fmt.Sprintf("%d", rand.Int63()), "new random world")

	debug = flag.Bool("v", false, "extra logs")
)

const (
	LEFT      = rpg.LEFT
	RIGHT     = rpg.RIGHT
	UP        = rpg.UP
	DOWN      = rpg.DOWN
	UPLEFT    = rpg.UPLEFT
	UPRIGHT   = rpg.UPRIGHT
	DOWNLEFT  = rpg.DOWNLEFT
	DOWNRIGHT = rpg.DOWNRIGHT
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
	if *debug {
		log.SetFlags(log.Lshortfile)
	} else {
		log.SetFlags(log.Lmicroseconds)

	}
	f, err := os.Create("p.debug")
	if err != nil {
		log.Fatal(err)
	}
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

	rand.Seed(time.Now().UnixNano())

	winbounds := pixel.R(0, 0, 800, 600)
	fmt.Println("Welcome to", rpg.Version())
	fmt.Println("Source code: https://github.com/aerth/rpg")
	fmt.Println("Please select screen resolution:")
	fmt.Println("1. 800x600")
	fmt.Println("2. 1024x768")
	fmt.Println("3. 1280x800")
	fmt.Println("4. 1280x800 undecorated")
	fmt.Println("Hit CTRL+F during normal gameplay for full screen toggle")
	var screenres int
	_, err = fmt.Scanf("%d", &screenres)
	if err != nil {
		fmt.Println("... choosing 800x600")
		screenres = 0
	}

	// window options
	cfg := pixelgl.WindowConfig{
		Title:       rpg.Version(),
		Bounds:      winbounds,
		Undecorated: false,
		VSync:       true,
	}

	switch screenres {
	default:
	case 2:
		winbounds = pixel.R(0, 0, 1024, 768)
	case 3:
		winbounds = pixel.R(0, 0, 1280, 800)
	case 4:
		log.Println("undecorated!")
		winbounds = pixel.R(0, 0, 1280, 800)
		cfg.Undecorated = true
	}

	cfg.Bounds = winbounds

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
	cursorsprite := rpg.GetCursor(1)

	// world generate
	world := rpg.NewWorld(*flaglevel, *flagenemies, *flagseed)
	if world == nil {
		return
	}
	// world tile sprite sheet
	spritesheet, spritemap := rpg.LoadSpriteSheet("tileset.png")

	// layers (TODO: slice?)
	// batch sprite drawing
	animbatch := pixel.NewBatch(&pixel.TrianglesData{}, spritesheet)

	// load loot sprite
	goldsheet, err := rpg.LoadPicture("sprites/loot.png")
	if err != nil {
		panic("need sprites/loot.png")
	}

	// full sprite
	goldsprite := pixel.NewSprite(goldsheet, goldsheet.Bounds())

	// loot batch
	lootbatch := pixel.NewBatch(&pixel.TrianglesData{}, goldsheet)

	// Fill in with water:
	// water world 67 wood, 114 117 182 special, 121 135 dirt, 128 blank, 20 grass
	//	rpg.DrawPattern(batch, spritemap[53], pixel.R(-3000, -3000, 3000, 3000), 100)

	// draw menu bar
	menubatch := pixel.NewBatch(&pixel.TrianglesData{}, spritesheet)
	common.DrawPattern(menubatch, spritemap[67], pixel.R(0, 0, win.Bounds().W()+20, 60), 100)
	for _, btn := range buttons {
		spritemap[200].Draw(menubatch, IM.Moved(btn.Frame.Center()))
	}

	// create Mobs
	world.NewMobs(*flagenemies)
	l := time.Now()
	var last = &l
	second := time.Tick(time.Second)
	frames := 0
	var camZoom = new(float64)
	var dt = new(float64)
	t1 := time.Now()
	fontsize := 36.00
	if win.Bounds().Max.X < 1100 {
		fontsize = 24.00
	}
	win.SetCursorVisible(false)
	text := rpg.NewText(fontsize)
	// start loop
	imd := imdraw.New(nil)
	rand.Seed(time.Now().UnixNano())
	var fullscreen = false
	//var latest string
	//redrawWorld(world)
	log.Println("Reticulating Splines...")
	globebatch := world.Init()

	var statustxt string
MainLoop:
	for !win.Closed() {
		// show title menu
		rpg.TitleMenu(win)

		// reset world
		world.Reset()
		var message = "Welcome to the world!\nCtrl+F Full Screen\nCtrl+Q Quit\n\nInventory: 'i'\nSPACE to continue"
		message = ""
	GameLoop:
		for !win.Closed() {

			*dt = time.Since(*last).Seconds()
			*last = time.Now()

			// check if ded
			if world.Char.Health < 1 {
				log.Println("GAME OVER")
				log.Printf("You survived for %s.\nYou acquired %s gold", time.Now().Sub(t1), world.Char.CountGold())
				log.Printf("Skeletons killed: %v", world.Char.Stats.Kills)
				log.Println(world.Char.StatsReport())

				break GameLoop
			}

			if message != "" {
				world.TextBox(win, message)
				message = ""
			}

			// zoom with mouse scroll, limit when not in debug mode
			*camZoom *= math.Pow(camZoomSpeed, win.MouseScroll().Y)
			if !*debug && *camZoom > 6.5 {
				*camZoom = 6.5
			}
			if !*debug && *camZoom < 2 {
				*camZoom = 2
			}
			if *debug {
				if *camZoom == 0 {
					*camZoom = 1
				}
			}
			cam := pixel.IM.Scaled(pixel.ZV, *camZoom).Moved(win.Bounds().Center()).Moved(world.Char.Rect.Center().Scaled(-*camZoom))

			// drawing
			win.Clear(colornames.Blue)
			animbatch.Clear()
			// if key
			if win.JustPressed(pixelgl.KeyQ) && win.Pressed(pixelgl.KeyLeftControl) {
				break MainLoop
			}
			if win.JustPressed(pixelgl.KeyF) && win.Pressed(pixelgl.KeyLeftControl) {
				fullscreen = !fullscreen
				if fullscreen {
					win.SetMonitor(pixelgl.PrimaryMonitor())
				} else {
					win.SetMonitor(nil)
				}
				menubatch.Clear()
				win.Update()
				newbounds := win.Bounds()
				common.DrawPattern(menubatch, spritemap[67], pixel.R(0, 0, newbounds.Max.X+20, 60), 100)
				for _, btn := range buttons {
					spritemap[200].Draw(menubatch, IM.Moved(btn.Frame.Center()))
				}
				menubatch.Draw(win)

			}
			// teleport random

			if win.JustPressed(pixelgl.Key8) {
				world.Char.Rect = common.DefaultSpriteRectangle.Moved(common.TileNear(world.Tiles, world.Char.Rect.Center()).Loc)
			}
			if win.JustPressed(pixelgl.KeyM) {
				world.NewMobs(1)
			}
			if win.Pressed(pixelgl.KeyLeftControl) && win.JustPressed(pixelgl.KeyL) {
				world.RandomLootSomewhere()
			}
			if win.Pressed(pixelgl.KeyLeftControl) && win.JustPressed(pixelgl.KeyK) {
				world.NewLoot(cam.Unproject(win.MousePosition()), []rpg.Item{rpg.RandomMagicItem()})
			}

			// move all enemies (debug)
			if win.JustPressed(pixelgl.Key9) {
				for _, v := range world.Entities {
					v.Rect = rpg.DefaultEntityRectangle.Moved(common.TileNear(world.Tiles, v.Rect.Center()).Loc)
				}
			}
			if win.JustReleased(pixelgl.KeyI) {
				rpg.InventoryLoop(win, world)
			}

			if win.JustPressed(pixelgl.KeyEqual) {
				*debug = !*debug
				if *debug {
					log.SetFlags(log.Lshortfile)
				} else {
					log.SetFlags(0)
				}
			}

			// get direction, velocity
			dir := controlswitch(dt, world, win)

			// check collision, move character
			world.Char.Update(*dt, dir, world)

			// update everything else (entities, DObjects)
			world.Update(*dt)
			world.Clean()

			// camera mode (center on player)
			win.SetMatrix(cam)

			// draw map to win (tiles and blocks)
			globebatch.Draw(win)

			// highlight enemy paths
			if *debug {
				world.HighlightPaths(win)
			}

			// draw entities and objects (not tiles and blocks)
			world.Draw(win)

			// draw animations such as magic spells
			imd.Clear()
			world.CleanAnimations()
			world.ShowAnimations(imd)
			imd.Draw(win)

			// highlight player tiles (left right up down and center)
			if *debug {
				for _, o := range world.Tile(world.Char.Rect.Center()).PathNeighbors() {
					ob := o.(common.Object)
					ob.W = world
					ob.Highlight(win, common.TransparentPurple)
				}
			}

			// draw all groundscores
			lootbatch.Clear()
			for _, dob := range world.DObjects {
				dob.Object.Draw(lootbatch, goldsheet, []*pixel.Sprite{goldsprite}, 0.5)
			}
			lootbatch.Draw(win)

			// back to window matrix
			win.SetMatrix(pixel.IM)

			// draw player in center of screen
			world.Char.Matrix = pixel.IM.Scaled(pixel.ZV, *camZoom).Scaled(pixel.ZV, 0.5).Moved(pixel.V(0, 0)).Moved(win.Bounds().Center())
			world.Char.Draw(win)

			// draw score board
			text.Clear()
			rpg.DrawScore(win.Bounds(), text, win,
				"%v HP · %v MP · %s GP · LVL %v · %v/%v XP · %v Kills %s", world.Char.Health, world.Char.Mana, world.Char.CountGold(), world.Char.Level, world.Char.Stats.XP, world.Char.NextLevel(), world.Char.Stats.Kills, statustxt)

			// draw menubar
			menubatch.Draw(win)
			if win.JustPressed(pixelgl.Key6) {
				//redrawWorld(world)
			}

			// draw health, mana, xp bars
			world.Char.DrawBars(win, win.Bounds())

			cursorsprite.Draw(win, pixel.IM.Scaled(pixel.ZV, 4).Moved(win.MousePosition()).Moved(pixel.V(0, -16)))

			// done drawing
			if !win.Pressed(pixelgl.KeyLeftControl) && win.Pressed(pixelgl.MouseButtonRight) {
				mouseloc := win.MousePosition()

				mcam := pixel.IM.Moved(win.Bounds().Center())
				mouseloc1 := mcam.Unproject(mouseloc)
				unit := mouseloc1.Unit()
				//                              log.Println("unit:", unit)
				dirmouse := rpg.UnitToDirection(unit)

				//                              log.Println("direction:", dir)

				switch dirmouse {

				case LEFT:
					world.Char.Phys.Vel.X = -world.Char.Phys.RunSpeed * (1 + *dt)
					world.Char.Dir = LEFT
				case RIGHT:
					world.Char.Phys.Vel.X = +world.Char.Phys.RunSpeed * (1 + *dt)
					world.Char.Dir = RIGHT
				case UP:
					world.Char.Phys.Vel.Y = +world.Char.Phys.RunSpeed * (1 + *dt)
					world.Char.Dir = UP
				case DOWN:
					world.Char.Phys.Vel.Y = -world.Char.Phys.RunSpeed * (1 + *dt)
					world.Char.Dir = DOWN
				case UPLEFT:
					world.Char.Phys.Vel.Y = +world.Char.Phys.RunSpeed * (1 + *dt)
					world.Char.Phys.Vel.X = -world.Char.Phys.RunSpeed * (1 + *dt)
					world.Char.Dir = UPLEFT
				case UPRIGHT:
					world.Char.Phys.Vel.X = +world.Char.Phys.RunSpeed * (1 + *dt)
					world.Char.Phys.Vel.Y = +world.Char.Phys.RunSpeed * (1 + *dt)
					world.Char.Dir = UPRIGHT
				case DOWNLEFT:
					world.Char.Phys.Vel.Y = -world.Char.Phys.RunSpeed * (1 + *dt)
					world.Char.Phys.Vel.X = -world.Char.Phys.RunSpeed * (1 + *dt)
					world.Char.Dir = DOWNLEFT
				case DOWNRIGHT:
					world.Char.Phys.Vel.X = +world.Char.Phys.RunSpeed * (1 + *dt)
					world.Char.Phys.Vel.Y = -world.Char.Phys.RunSpeed * (1 + *dt)
					world.Char.Dir = DOWNRIGHT
				default:
				}
			}
			if !win.Pressed(pixelgl.KeyLeftControl) && win.JustPressed(pixelgl.MouseButtonLeft) {
				mouseloc := win.MousePosition()
				if b, f, ok := world.IsButton(buttons, mouseloc); ok {
					log.Println(mouseloc)
					log.Printf("Clicked button: %q", b.Name)
					f(win, world)

				} else {

					// pick up loot
					if loot, ok := world.IsLoot(cam.Unproject(mouseloc)); ok {
						statustxt = world.Char.PickUp(loot)

					} else {
						mcam := pixel.IM.Moved(win.Bounds().Center())
						mouseloc1 := mcam.Unproject(mouseloc)
						// magic bullet
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

				if *debug {
					log.Println(frames, "frames per second")
					log.Println(len(world.DObjects), "dynamic objects")
					log.Println(len(world.Animations), "animations")
					log.Println(len(world.Entities), "living entities")
				}
				frames = 0
			}

		}
	}
	log.Printf("You survived for %s.\nYou acquired %s gold", time.Now().Sub(t1), world.Char.CountGold())
	log.Println("Inventory:", rpg.FormatItemList(world.Char.Inventory))
	log.Printf("Skeletons killed: %v", world.Char.Stats.Kills)
	log.Println(world.Char.StatsReport())

}

func controlswitch(dt *float64, w *rpg.World, win *pixelgl.Window) rpg.Direction {
	// Manastorm
	if win.JustPressed(pixelgl.KeySpace) || win.JustPressed(pixelgl.MouseButtonMiddle) {
		if w.Char.Mana > 0 {
			w.Action(w.Char, w.Char.Rect.Center(), rpg.ManaStorm, w.Char.Dir)
		} else {
			log.Println("Not enough mana")
		}
	}

	// Magic bullet
	if win.JustPressed(pixelgl.KeyB) {
		if w.Char.Mana > 0 {
			w.Action(w.Char, w.Char.Rect.Center(), rpg.MagicBullet, w.Char.Dir)
		} else {
			log.Println("Not enough mana")
		}
	}

	// Slow Cheat
	if win.Pressed(pixelgl.KeyTab) {
		*dt /= 8
	}
	// Fast Cheat
	if win.Pressed(pixelgl.KeyLeftShift) {
		*dt *= 4
	}

	// MP Cheat
	if win.Pressed(pixelgl.Key1) {
		w.Char.Mana += 1
		if w.Char.Mana > 255 {
			w.Char.Mana = 255
		}
	}

	// HP Cheat
	if win.Pressed(pixelgl.Key2) {
		w.Char.Health += 1
		if w.Char.Health > 255 {
			w.Char.Health = 255
		}
	}

	// XP Cheat
	if win.Pressed(pixelgl.Key3) {
		w.Char.Stats.XP += 10
	}

	// Fly mode
	if win.Pressed(pixelgl.KeyCapsLock) {
		w.Char.Phys.CanFly = !w.Char.Phys.CanFly
	}

	// control character with direction and velocity (collision check elsewhere)
	dir := w.Char.Dir
	if !win.Pressed(pixelgl.KeyLeftControl) && (win.Pressed(pixelgl.KeyLeft) || win.Pressed(pixelgl.KeyH) || win.Pressed(pixelgl.KeyA)) {
		w.Char.Phys.Vel.X = -w.Char.Phys.RunSpeed * (1 + *dt)
		dir = LEFT
	}
	if !win.Pressed(pixelgl.KeyLeftControl) && (win.Pressed(pixelgl.KeyRight) || win.Pressed(pixelgl.KeyL) || win.Pressed(pixelgl.KeyD)) {
		w.Char.Phys.Vel.X = +w.Char.Phys.RunSpeed * (1 + *dt)
		dir = RIGHT
	}
	if !win.Pressed(pixelgl.KeyLeftControl) && (win.Pressed(pixelgl.KeyDown) || win.Pressed(pixelgl.KeyJ) || win.Pressed(pixelgl.KeyS)) {
		w.Char.Phys.Vel.Y = -w.Char.Phys.RunSpeed * (1 + *dt)
		dir = DOWN

	}
	if !win.Pressed(pixelgl.KeyLeftControl) && (win.Pressed(pixelgl.KeyUp) || win.Pressed(pixelgl.KeyK) || win.Pressed(pixelgl.KeyW)) {
		w.Char.Phys.Vel.Y = +w.Char.Phys.RunSpeed * (1 + *dt)
		dir = UP
	}

	// TODO: fix double velocity on diag
	if !win.Pressed(pixelgl.KeyLeftControl) && (win.Pressed(pixelgl.KeyUp) && win.Pressed(pixelgl.KeyLeft)) {
		dir = rpg.UPLEFT
	}
	if !win.Pressed(pixelgl.KeyLeftControl) && (win.Pressed(pixelgl.KeyUp) && win.Pressed(pixelgl.KeyRight)) {
		dir = rpg.UPRIGHT
	}
	if !win.Pressed(pixelgl.KeyLeftControl) && (win.Pressed(pixelgl.KeyDown) && win.Pressed(pixelgl.KeyLeft)) {
		dir = rpg.DOWNLEFT
	}
	if !win.Pressed(pixelgl.KeyLeftControl) && (win.Pressed(pixelgl.KeyDown) && win.Pressed(pixelgl.KeyRight)) {
		dir = rpg.DOWNRIGHT
	}
	return dir
}
func main() {
	flag.Parse()
	pixelgl.Run(run)
}
