// Code generated by "stringer -output strings.go -type EntityType,EntityState,ItemType,ObjectType,animState,ActionType,StatusEffect"; DO NOT EDIT.

package rpg

import "fmt"

const _EntityType_name = "SKELETONSKELETON_GUARDDOBJECT"

var _EntityType_index = [...]uint8{0, 8, 22, 29}

func (i EntityType) String() string {
	if i < 0 || i >= EntityType(len(_EntityType_index)-1) {
		return fmt.Sprintf("EntityType(%d)", i)
	}
	return _EntityType_name[_EntityType_index[i]:_EntityType_index[i+1]]
}

const _EntityState_name = "S_IDLES_RUNS_WANDERS_GUARDS_SUSPECTS_HUNTS_DEAD"

var _EntityState_index = [...]uint8{0, 6, 11, 19, 26, 35, 41, 47}

func (i EntityState) String() string {
	if i < 0 || i >= EntityState(len(_EntityState_index)-1) {
		return fmt.Sprintf("EntityState(%d)", i)
	}
	return _EntityState_name[_EntityState_index[i]:_EntityState_index[i+1]]
}

const _ItemType_name = "GOLDPOTIONFOODWEAPONARMORSPECIAL"

var _ItemType_index = [...]uint8{0, 4, 10, 14, 20, 25, 32}

func (i ItemType) String() string {
	i -= 1
	if i < 0 || i >= ItemType(len(_ItemType_index)-1) {
		return fmt.Sprintf("ItemType(%d)", i+1)
	}
	return _ItemType_name[_ItemType_index[i]:_ItemType_index[i+1]]
}

const _ObjectType_name = "O_NONEO_TILEO_BLOCKO_INVISIBLEO_SPECIALO_WINO_DYNAMIC"

var _ObjectType_index = [...]uint8{0, 6, 12, 19, 30, 39, 44, 53}

func (i ObjectType) String() string {
	if i < 0 || i >= ObjectType(len(_ObjectType_index)-1) {
		return fmt.Sprintf("ObjectType(%d)", i)
	}
	return _ObjectType_name[_ObjectType_index[i]:_ObjectType_index[i+1]]
}

const _animState_name = "IdleRunning"

var _animState_index = [...]uint8{0, 4, 11}

func (i animState) String() string {
	if i < 0 || i >= animState(len(_animState_index)-1) {
		return fmt.Sprintf("animState(%d)", i)
	}
	return _animState_name[_animState_index[i]:_animState_index[i+1]]
}

const _ActionType_name = "TalkSlashManaStormMagicBulletArrow"

var _ActionType_index = [...]uint8{0, 4, 9, 18, 29, 34}

func (i ActionType) String() string {
	if i < 0 || i >= ActionType(len(_ActionType_index)-1) {
		return fmt.Sprintf("ActionType(%d)", i)
	}
	return _ActionType_name[_ActionType_index[i]:_ActionType_index[i+1]]
}

const _StatusEffect_name = "E_POISONE_PARALYSISE_FROZENE_BURNEDE_SLEEPE_CONFUSEDE_TIREDE_DRUNK"

var _StatusEffect_index = [...]uint8{0, 8, 19, 27, 35, 42, 52, 59, 66}

func (i StatusEffect) String() string {
	i -= 1
	if i < 0 || i >= StatusEffect(len(_StatusEffect_index)-1) {
		return fmt.Sprintf("StatusEffect(%d)", i+1)
	}
	return _StatusEffect_name[_StatusEffect_index[i]:_StatusEffect_index[i+1]]
}
