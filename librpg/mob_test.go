package rpg

import (
	"log"
	"testing"

	"github.com/faiface/pixel"
)

func TestMobCenter(t *testing.T) {
	w := NewWorld("world name", 1, "1")
	mob := w.NewEntity(SKELETON)
	mob.Rect = mob.Rect.Moved(pixel.V(100, 100))
	if mob.Center() != pixel.V(100, 100) {
		t.FailNow()
	}
	log.Println(mob.Center())
}
