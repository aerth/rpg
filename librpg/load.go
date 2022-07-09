// aerth game
// copyright 2017 aerth <aerth@riseup.net>
package rpg

import (
	"math/rand"
	"time"

	_ "image/png"

	"github.com/aerth/rpc/librpg/common"
	"github.com/faiface/pixel"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// loadCharacterSheet returns an animated spritesheet
// 13W 21H
func LoadEntitySheet(sheetPath string, framesx, framesy uint8) (sheet pixel.Picture, anims map[EntityState]map[Direction][]pixel.Rect, err error) {
	sheet, err = common.LoadPicture(sheetPath)
	frameWidth := float64(int(sheet.Bounds().Max.X / float64(framesx)))
	frameHeight := float64(int(sheet.Bounds().Max.Y / float64(framesy)))
	//log.Println(frameWidth, "width", frameHeight, "height")
	// create a array of frames inside the spritesheet
	var frames = []pixel.Rect{}
	for y := 0.00; y+frameHeight <= sheet.Bounds().Max.Y; y = y + frameHeight {
		for x := 0.00; x+float64(frameWidth) <= sheet.Bounds().Max.X; x = x + float64(frameWidth) {
			frames = append(frames, pixel.R(
				x,
				y,
				x+frameWidth,
				y+frameHeight,
			))
		}
	}

	//log.Println("total skeleton frames", len(frames))

	// 0-5 die
	// BLANK 6-12
	// 13-25 shoot right
	// 26-39 shoot down
	// 6-76 shoot left
	// 7-25 shoot up
	anims = make(map[EntityState]map[Direction][]pixel.Rect)
	anims[S_IDLE] = make(map[Direction][]pixel.Rect)
	anims[S_WANDER] = make(map[Direction][]pixel.Rect)
	anims[S_RUN] = make(map[Direction][]pixel.Rect)
	anims[S_GUARD] = make(map[Direction][]pixel.Rect)
	anims[S_SUSPECT] = make(map[Direction][]pixel.Rect)
	anims[S_HUNT] = make(map[Direction][]pixel.Rect)
	anims[S_DEAD] = make(map[Direction][]pixel.Rect)

	// spritesheet is right down left up
	anims[S_DEAD][LEFT] = frames[0:5]
	anims[S_DEAD][RIGHT] = frames[0:5]
	anims[S_DEAD][UP] = frames[0:5]
	anims[S_DEAD][DOWN] = frames[0:5]
	anims[S_IDLE][LEFT] = frames[143:144]
	anims[S_IDLE][UP] = frames[156:157]
	anims[S_IDLE][RIGHT] = frames[169:170]
	anims[S_IDLE][DOWN] = frames[182:183]
	anims[S_RUN][LEFT] = frames[143:152]
	anims[S_RUN][UP] = frames[156:165]
	anims[S_RUN][RIGHT] = frames[169:178]
	anims[S_RUN][DOWN] = frames[182:191]
	return sheet, anims, nil
}

// loadCharacterSheet returns an animated spritesheet
func LoadCharacterSheet(sheetPath string, numframes uint8) (sheet pixel.Picture, anims map[Direction][]pixel.Rect, err error) {
	sheet, err = common.LoadPicture("sprites/char.png")
	if err != nil {
		panic(err)
	}
	frameWidth := int(sheet.Bounds().Max.X/float64(numframes)) * 2
	//log.Println(frameWidth, "width")
	// create a array of frames inside the spritesheet
	var frames = new([16]pixel.Rect)

	for i, x := 0, 0.0; x+float64(frameWidth) <= sheet.Bounds().Max.X; i, x = i+1, x+float64(frameWidth) {
		if i > 15 {
			break
		}
		frames[i] = pixel.R(
			x,
			0,
			x+float64(frameWidth),
			sheet.Bounds().H(),
		)
	}

	anims = make(map[Direction][]pixel.Rect)
	anims[LEFT] = frames[:4]
	anims[UPLEFT] = frames[:4]
	anims[RIGHT] = frames[4:8]
	anims[DOWNRIGHT] = frames[4:8]
	anims[DOWN] = frames[8:12]
	anims[DOWNLEFT] = frames[8:12]
	anims[UP] = frames[12:]
	anims[UPRIGHT] = frames[12:]
	return sheet, anims, nil
}

func LoadNewCharacterSheet(sheetPath string) (sheet pixel.Picture, anims map[Direction][]pixel.Rect, err error) {
	numframes := 16.00
	sheet, err = common.LoadPicture(sheetPath)
	if err != nil {
		panic(err)
	}
	var frames = new([16]pixel.Rect)
	frameWidth := sheet.Bounds().Max.X / float64(len(frames))

	for i := 0.00; i < numframes; i++ {

		frames[int(i)] = pixel.R(
			(i+1.00)*frameWidth,
			0,
			((i+1.00)*frameWidth)+float64(frameWidth),
			sheet.Bounds().H(),
		)

	}

	anims = make(map[Direction][]pixel.Rect)
	anims[LEFT] = frames[:4]
	anims[UPLEFT] = frames[:4]
	anims[RIGHT] = frames[4:8]
	anims[DOWNRIGHT] = frames[4:8]
	anims[DOWN] = frames[8:12]
	anims[DOWNLEFT] = frames[8:12]
	anims[UP] = frames[12:]
	anims[UPRIGHT] = frames[12:]
	return sheet, anims, nil

}

func LoadSpriteSheet(path string) (pixel.Picture, []*pixel.Sprite) {
	spritesheet, err := common.LoadPicture("sprites/" + path)
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
