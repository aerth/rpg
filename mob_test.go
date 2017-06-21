package rpg

import (
	"log"
	"testing"

	"github.com/faiface/pixel"
)

func TestMobCenter(t *testing.T) {
	w := NewWorld("1", 1)
	mob := w.NewEntity(SKELETON)
	mob.Rect = mob.Rect.Moved(pixel.V(100, 100))
	log.Println(mob.Center())
}
