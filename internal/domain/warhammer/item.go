package warhammer

type WhItemType int

func (input WhItemType) Copy() WhItemType {
	return input
}

func getAllowedItemType() string {
	var list = map[string]int{
		"melee":      0,
		"ranged":     1,
		"ammunition": 2,
		"armour":     3,
		"container":  4,
		"other":      5,
		"grimoire":   6,
	}

	return formatAllowedIntTypes(list)
}

type WhItemHands int

func (input WhItemHands) Copy() WhItemHands {
	return input
}

func getAllowedItemHands() string {
	var list = map[string]int{
		"one": 1,
		"two": 2,
	}
	return formatAllowedIntTypes(list)
}

type WhItemMeleeReach int

func (input WhItemMeleeReach) Copy() WhItemMeleeReach {
	return input
}

func getAllowedItemMeleeReach() string {
	var list = map[string]int{
		"personal":   0,
		"very_short": 1,
		"short":      2,
		"average":    3,
		"long":       4,
		"very_long":  5,
		"massive":    6,
	}
	return formatAllowedIntTypes(list)
}

type WhItemMeleeGroup int

func (input WhItemMeleeGroup) Copy() WhItemMeleeGroup {
	return input
}

func getAllowedItemMeleeGroup() string {
	var list = map[string]int{
		"basic":      0,
		"cavalry":    1,
		"fencing":    2,
		"brawling":   3,
		"flail":      4,
		"parry":      5,
		"polearm":    6,
		"two_handed": 7,
	}
	return formatAllowedIntTypes(list)
}
