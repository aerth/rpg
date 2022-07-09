package rpg

import "fmt"

type Stats struct {
	Strength     float64
	Wisdom       float64
	Intelligence float64
	Vitality     float64
	XP           uint64
	Kills        uint64
	Score        uint64
}

func (s Stats) String() string {
	f := `
=== CHARACTER STATS ===
STR %v
WIS %v
INT %v
VIT %v
XP %v
Kills %v
Score %v
`

	return fmt.Sprintf(f, int(s.Strength), int(s.Wisdom), int(s.Intelligence), int(s.Vitality),
		s.XP, s.Kills, s.Score)
}
