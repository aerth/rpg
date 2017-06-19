package rpg

import (
	"log"
	"math/rand"
	"testing"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func TestCreateLoot(t *testing.T) {
	i := createLoot()
	log.Println(i)

}
