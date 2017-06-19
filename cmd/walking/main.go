package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"runtime/pprof"
	"strconv"
	"time"

	"golang.org/x/image/colornames"

	_ "image/png"

	"github.com/aerth/rpg"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var (
	flagenemies   = flag.Int("e", 2, "number of enemies to begin with")
	flaglevel     = flag.Int("lvl", 1, "starting level (1)")
	flagcustomlvl = flag.String("test", "", "custom world test (filename)")
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
	NUMENEMY = 2
	LEVEL    = 1
)
var (
	defaultzoom  = 1.0
	camZoomSpeed = 1.20
)

func run() {
	NUMENEMY := *flagenemies
	LEVEL := *flaglevel
	TESTLVL := *flagcustomlvl
	f, err := os.Create("p.debug")
	if err != nil {
		log.Fatal(err)
	}
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

	rand.Seed(time.Now().UnixNano())
	winbounds := pixel.R(0, 0, 800, 600)
	//	winbounds = pixel.R(0, 0, 1024, 768)
	//	winbounds = pixel.R(0, 0, 1280, 720)
	// window options
	cfg := pixelgl.WindowConfig{
		Title:  "AERPG",
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
	//char.Rect = char.Rect.Moved(V(33, 33))
	// load world
	worldbounds := pixel.R(float64(-3000), float64(-3000), float64(3000), float64(3000))
	//	worldbounds = pixel.R(float64(-4000), float64(-4000), float64(4000), float64(4000))
	world := rpg.NewWorld(strconv.Itoa(LEVEL), worldbounds, TESTLVL)
	char := rpg.NewCharacter()
	char.Inventory = []rpg.Item{rpg.MakeGold(uint64(rand.Intn(7)))} // start with some loot
	char.Rect = char.Rect.Moved(rpg.FindRandomTile(world.Objects))
	world.Char = char

	// sprite sheet
	spritesheet, spritemap := rpg.LoadSpriteSheet("tileset.png")

	// layers (TODO: slice?)
	// batch sprite drawing
	globebatch := pixel.NewBatch(&pixel.TrianglesData{}, spritesheet)
	animbatch := pixel.NewBatch(&pixel.TrianglesData{}, spritesheet)
	menubatch := pixel.NewBatch(&pixel.TrianglesData{}, spritesheet)

	// water world 67 wood, 114 117 182 special, 121 135 dirt, 128 blank, 20 grass
	//	rpg.DrawPattern(batch, spritemap[53], pixel.R(-3000, -3000, 3000, 3000), 100)

	rpg.DrawPattern(menubatch, spritemap[67], pixel.R(0, 0, win.Bounds().W()+20, 60), 100)
	//rpg.DrawPattern(menubatch, spritemap[230], pixel.R(20, 20, win.Bounds().W()+20, 80), 100)
	for _, btn := range buttons {
		spritemap[200].Draw(menubatch, IM.Moved(btn.Frame.Center()))
	}

	redrawWorld := func() {
		globebatch.Clear()
		// draw it on to canvasglobe
		for i := range world.Objects {
			world.Objects[i].Draw(globebatch, spritesheet, spritemap)
		}
	}

	// create NPC

	world.NewMobs(NUMENEMY)
	l := time.Now()
	var last = &l
	second := time.Tick(time.Second)
	tick := time.Tick(time.Second)

	frames := 0
	var delda float64 = 0.00
	var camZoom = &defaultzoom
	var debug bool
	var dt *float64 = &delda
	t1 := time.Now()
	text := rpg.NewText(36)
	texthelp := rpg.NewTextSmooth(18)
	texthelp.Color = colornames.Red
	fmt.Fprint(texthelp, "[tab=slow] [shift=fast] [q=quit] [space=manastorm] [enter=reset] [i=inventory]")
	redrawWorld()
	// start loop
	imd := imdraw.New(nil)
	rand.Seed(time.Now().UnixNano())
	var latest string
	for !win.Closed() {
		rpg.TitleMenu(world, win)
		world.Reset()
		char.Health = 255
		for !win.Closed() {

			if char.Health < 1 {
				log.Println("GAME OVER")
				break
			}
			*dt = time.Since(*last).Seconds()
			*last = time.Now()
			// zoom with mouse scroll

			*camZoom *= math.Pow(camZoomSpeed, win.MouseScroll().Y)

			// drawing
			//win.Clear(rpg.RandomColor())
			win.Clear(colornames.Black)
			animbatch.Clear()
			// if key
			if win.JustPressed(pixelgl.KeyQ) {
				break
			}
			if win.JustReleased(pixelgl.KeyI) {
				rpg.InventoryLoop(win, world)
			}

			if win.JustPressed(pixelgl.KeyEqual) {
				debug = !debug
			}

			dir := controlswitch(dt, world, win, char, buttons, win)
			char.Update(*dt, dir, world)
			world.Update(*dt)

			char.Matrix = pixel.IM.Scaled(pixel.ZV, *camZoom).Moved(win.Bounds().Center())
			cam := char.Matrix.Moved(char.Rect.Center().Scaled(-*camZoom))
			win.SetMatrix(cam)

			// draw map (tiles and blocks) (never updated for now)
			//canvasglobe.Draw(win, pixel.IM)
			globebatch.Draw(win)

			// draw entities and objects (not tiles and blocks)
			world.Draw(win) // was win
			// highlight paths
			if debug {
				world.HighlightPaths(win)
			}

			imd.Clear()
			world.CleanAnimations()
			world.ShowAnimations(imd)
			imd.Draw(win)

			// back to window cam
			win.SetMatrix(pixel.IM)
			char.Draw(win)
			text.Clear()
			rpg.DrawScore(winbounds, text, win,
				"[%vHP·%vMP·%sGP LVL%v %vXP %vKills] %s", char.Health, char.Mana, char.CountGold(), char.Level, char.Stats.XP, char.Stats.Kills, latest)

			menubatch.Draw(win)

			mouseloc := win.MousePosition()
			if win.JustPressed(pixelgl.MouseButtonLeft) {
				if b, f, ok := world.IsButton(buttons, mouseloc); ok {
					if debug {
						log.Printf("Clicked button: %q", b.Name)
					}
					f(win, world)
				}
			}
			select {
			default: //
			case <-tick:
				if len(world.Messages) > 100 {
					log.Println("truncating messages")
					world.Messages = []string{world.Messages[len(world.Messages)-1]}
				}
				if len(world.Messages) > 0 {
					latest = world.Messages[0]
					if len(world.Messages) > 1 {
						world.Messages = world.Messages[1:]
					} else {
						world.Messages = []string{}
					}

				}

			}
			if b, _, ok := world.IsButton(buttons, mouseloc); ok {
				texthelp.Clear()
				texthelp.Dot = texthelp.Orig
				fmt.Fprintf(texthelp, "%q", b.Name)
			}

			texthelp.Draw(win, pixel.IM.Moved(pixel.V(150, 60)))

			//spritemap[20].Draw(menubar, pixel.IM.Scaled(ZV, 10).Moved(pixel.V(30, 30)))
			//menubar.Draw(win, pixel.IM)
			win.Update()
			world.Clean()

			// fps, gps
			frames++
			gps := char.Rect.Center()
			select {
			default: //keep going
			case <-second:
				latest = ""
				str := fmt.Sprintf(""+
					"FPS: %d | GPS: (%v,%v) | VEL: (%v) | HP: (%v) ",
					frames, int(gps.X), int(gps.Y), int(char.Phys.Vel.Len()), char.Health)
				win.SetTitle(str)
				frames = 0
			}

		}
		log.Printf("You survived for %s.\nYou acquired %s gold", time.Now().Sub(t1), char.CountGold())
		log.Println("Inventory:", rpg.FormatItemList(char.Inventory))
		log.Printf("Skeletons killed: %v", char.Stats.Kills)
	}

}

func controlswitch(dt *float64, w *rpg.World, win *pixelgl.Window, char *rpg.Character, buttons []rpg.Button, buf pixel.Target) rpg.Direction {
	if win.JustPressed(pixelgl.KeySpace) {
		if char.Mana > 0 {
			w.Action(char, char.Rect.Center(), rpg.ManaStorm)
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
		*dt *= 15
	}
	if win.Pressed(pixelgl.Key1) {
		char.Mana += 1
		if char.Mana > 255 {
			char.Mana = 255
		}
	}
	if win.Pressed(pixelgl.KeyCapsLock) {
		char.Phys.CanFly = !char.Phys.CanFly
	}
	dir := char.Dir

	// disable momentum
	if win.JustPressed(pixelgl.KeyS) {
		char.Phys.Vel = pixel.ZV
	}

	if win.JustPressed(pixelgl.MouseButtonLeft) {
		mouseloc := win.MousePosition()
		log.Println(mouseloc)
		if b, f, ok := w.IsButton(buttons, mouseloc); ok {
			log.Printf("Clicked button: %q", b.Name)
			f(win, w)

		}
	}

	if win.Pressed(pixelgl.KeyLeft) || win.Pressed(pixelgl.KeyH) || win.Pressed(pixelgl.KeyA) {
		dir = LEFT
		char.Phys.Vel.X = -char.Phys.RunSpeed * (1 + *dt)
	}
	if win.Pressed(pixelgl.KeyRight) || win.Pressed(pixelgl.KeyL) || win.Pressed(pixelgl.KeyD) {
		char.Phys.Vel.X = +char.Phys.RunSpeed * (1 + *dt)
		dir = RIGHT
	}
	if win.Pressed(pixelgl.KeyDown) || win.Pressed(pixelgl.KeyJ) || win.Pressed(pixelgl.KeyS) {
		char.Phys.Vel.Y = -char.Phys.RunSpeed * (1 + *dt)
		dir = DOWN

	}
	if win.Pressed(pixelgl.KeyUp) || win.Pressed(pixelgl.KeyK) || win.Pressed(pixelgl.KeyW) {

		char.Phys.Vel.Y = +char.Phys.RunSpeed * (1 + *dt)
		dir = UP
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
