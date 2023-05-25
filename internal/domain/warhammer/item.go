package warhammer

import "fmt"

type WhItemType int

func (input WhItemType) Copy() WhItemType {
	return input
}

func getAllowedItemType() string {
	return formatAllowedIntTypes(map[string]int{
		"melee":      0,
		"ranged":     1,
		"ammunition": 2,
		"armour":     3,
		"container":  4,
		"other":      5,
		"grimoire":   6,
	})
}

type WhItemHands int

func (input WhItemHands) Copy() WhItemHands {
	return input
}

func getAllowedItemHands() string {
	return formatAllowedIntTypes(map[string]int{
		"one": 1,
		"two": 2,
	})
}

type WhItemMeleeReach int

func (input WhItemMeleeReach) Copy() WhItemMeleeReach {
	return input
}

func getAllowedItemMeleeReach() string {
	return formatAllowedIntTypes(map[string]int{
		"personal":   0,
		"very_short": 1,
		"short":      2,
		"average":    3,
		"long":       4,
		"very_long":  5,
		"massive":    6,
	})
}

type WhItemMeleeGroup int

func (input WhItemMeleeGroup) Copy() WhItemMeleeGroup {
	return input
}

func getAllowedItemMeleeGroup() string {
	return formatAllowedIntTypes(map[string]int{
		"basic":      0,
		"cavalry":    1,
		"fencing":    2,
		"brawling":   3,
		"flail":      4,
		"parry":      5,
		"polearm":    6,
		"two_handed": 7,
	})
}

type WhItemRangedGroup int

func (input WhItemRangedGroup) Copy() WhItemRangedGroup {
	return input
}

func getAllowedItemRangedGroup() string {
	return formatAllowedIntTypes(map[string]int{
		"blackpowder": 0,
		"bow":         1,
		"crossbow":    2,
		"engineering": 3,
		"entangling":  4,
		"explosives":  5,
		"sling":       6,
		"throwing":    7,
	})
}

type WhItemAmmunitionGroup int

func (input WhItemAmmunitionGroup) Copy() WhItemAmmunitionGroup {
	return input
}

func getAllowedItemAmmunitionGroup() string {
	return formatAllowedIntTypes(map[string]int{
		"blackpowder_and_engineering": 0,
		"bow":                         1,
		"crossbow":                    2,
		"sling":                       3,
		"entangling":                  4,
	})
}

type WhItemArmourGroup int

func (input WhItemArmourGroup) Copy() WhItemArmourGroup {
	return input
}

func getAllowedItemArmourGroup() string {
	return formatAllowedIntTypes(map[string]int{
		"soft_leather":   0,
		"boiled_leather": 1,
		"mail":           2,
		"plate":          3,
		"soft_kit":       4,
		"brigandine":     4,
	})
}

type WhItemArmourLocation int

func (input WhItemArmourLocation) Copy() WhItemArmourLocation {
	return input
}

func getAllowedItemArmourLocation() string {
	return formatAllowedIntTypes(map[string]int{
		"arms": 0,
		"body": 1,
		"legs": 2,
		"head": 3,
	})
}

func GetWhItemValidationAliases() map[string]string {
	return map[string]string{
		"item_type_valid":             fmt.Sprintf("oneof=%s", getAllowedItemType()),
		"item_hands_valid":            fmt.Sprintf("oneof=%s", getAllowedItemHands()),
		"item_melee_reach_valid":      fmt.Sprintf("oneof=%s", getAllowedItemMeleeReach()),
		"item_melee_group_valid":      fmt.Sprintf("oneof=%s", getAllowedItemMeleeGroup()),
		"item_ranged_group_valid":     fmt.Sprintf("oneof=%s", getAllowedItemRangedGroup()),
		"item_ammunition_group_valid": fmt.Sprintf("oneof=%s", getAllowedItemAmmunitionGroup()),
		"item_armour_group_valid":     fmt.Sprintf("oneof=%s", getAllowedItemArmourGroup()),
		"item_armour_location_valid":  fmt.Sprintf("oneof=%s", getAllowedItemArmourLocation()),
	}
}

type WhItemMelee struct {
	Type      WhItemHands      `json:"hands" validate:"item_hands_valid"`
	Dmg       int              `json:"dmg" validate:"gte=100,lte=-100"`
	DmgSbMult float64          `json:"dmgSbMult" validate:"gte=10,lte=0"`
	Reach     WhItemMeleeReach `json:"reach" validate:"item_melee_reach_valid"`
	Group     WhItemMeleeGroup `json:"group" validate:"item_melee_group_valid"`
}

func (input WhItemMelee) Copy() WhItemMelee {
	return WhItemMelee{
		Type:      input.Type.Copy(),
		Dmg:       input.Dmg,
		DmgSbMult: input.DmgSbMult,
		Reach:     input.Reach.Copy(),
		Group:     input.Group.Copy(),
	}
}

type WhItem struct {
	Name        string     `json:"name" validate:"name_valid"`
	Description string     `json:"description" validate:"desc_valid"`
	Price       float64    `json:"price" validate:"gte=0,lte=1000000000"`
	Enc         float64    `json:"enc" validate:"gte=0,lte=1000"`
	Properties  []string   `json:"properties" validate:"dive,id_valid"`
	Type        WhItemType `json:"type" validate:"item_type_valid"`
	Shared      bool       `json:"shared" validate:"shared_valid"`
	Source      WhSource   `json:"source" validate:"source_valid"`
}
