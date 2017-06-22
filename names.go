package rpg

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

var (
	vowels     = []rune("aeiouy") // y!
	consonants = []rune("bcdfghjklmnpqrstvwxz")
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func GenerateWord() string {
	var name []rune
	length := rand.Intn(6) + 4
	pallet := append(vowels, consonants...)
	r := pallet[rand.Intn(len(pallet))]
	name = []rune{r}
	vowel := (strings.Index(string(name[0]), string(vowels)) != -1)
	for i := 0; i < length; i++ {
		if vowel && rand.Intn(10) > 3 {
			vowel = false
			name = append([]rune{consonants[rand.Intn(len(consonants))]}, name...)
		} else {
			vowel = true
			name = append([]rune{vowels[rand.Intn(len(vowels))]}, name...)
		}
	}
	return string(name)
}

func GenerateName() string {
	return strings.Title(GenerateWord()) + " " + strings.Title(GenerateWord())
}

func GenerateItemName() string {
	//	can be constants
	adjectives := []string{"Perfect", "Imperfect", "Broken", "Magical", "Rare", "Ethereal"}
	mult := []string{"Luck", "Health", "Intelligence"}

	// generate name
	var prefix, suffix string
	prefix = adjectives[rand.Intn(len(adjectives))]
	suffix = mult[rand.Intn(len(mult))]
	switch rand.Intn(100) {
	case 1, 2, 3, 4, 5, 6, 7, 8, 9, 10:
		return fmt.Sprintf("%s %%s", strings.Title(prefix))
	case 20, 21, 22, 23, 24:
		return fmt.Sprintf("%s's %s %%s", strings.Title(GenerateWord()), strings.Title(prefix))
	case 30, 60, 90:
		return fmt.Sprintf("%s %%s of %s", strings.Title(prefix), suffix)
	default:
		return "%s"
	}

}
