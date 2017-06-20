package rpg

import (
	"log"
	"math/rand"
	"strings"
	"testing"
)

func TestGenerateName(t *testing.T) {
	for i := 0; i < 100; i++ {
		log.Printf("%s from %s", GenerateName(), strings.Title(GenerateWord()))
	}
}

func TestRand(t *testing.T) {
	for i := 0; i < 100; i++ {

		switch rand.Intn(3) {
		case 0:
			log.Println("0")
		case 1:
			log.Println("1")
		case 2:
			log.Println("2")
		case 3: // never hits
			log.Println("3")
		}
	}
}

func TestGenerateItemName(t *testing.T) {
	for i := 0; i < 100; i++ {
		item := createItemLoot()
		log.Printf(GenerateItemName(), item)
	}
}
