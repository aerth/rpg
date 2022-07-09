package rpg

type StatusEffect int

const (
	_ StatusEffect = iota
	E_POISON
	E_PARALYSIS
	E_FROZEN
	E_BURNED
	E_SLEEP
	E_CONFUSED
	E_TIRED
	E_DRUNK
)
