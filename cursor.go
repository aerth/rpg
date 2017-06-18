package rpg

import "github.com/faiface/pixel"

func GetCursor(num int) *pixel.Sprite {
	//spritesheet, spritemap := LoadSpriteSheet("cursor.png")
	pic, err := LoadPicture("sprites/cursor.png")
	if err != nil {

		panic(err)
	}
	sprite := pixel.NewSprite(pic, pic.Bounds())
	return sprite
}
