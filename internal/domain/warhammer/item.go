package warhammer

import (
	"strconv"
	"strings"
)

type WhItemType int

func (input WhItemType) Copy() WhItemType {
	return input
}

func getAllowedItemType() string {
	var list = map[string]WhItemType{
		"melee":      0,
		"ranged":     1,
		"ammunition": 2,
		"armour":     3,
		"container":  4,
		"other":      5,
		"grimoire":   6,
	}

	values := make([]string, 0, len(list))
	for _, v := range list {
		values = append(values, strconv.Itoa(int(v)))
	}
	return strings.Join(values, " ")
}
