package warhammer

import (
	"fmt"
	"strings"
)

type WhItemType int

const (
	WhItemTypeMelee      = 0
	WhItemTypeRanged     = 1
	WhItemTypeAmmunition = 2
	WhItemTypeArmour     = 3
	WhItemTypeContainer  = 4
	WhItemTypeGrimoire   = 6
	WhItemTypeOther      = 5
)

func itemTypeValues() string {
	return formatIntegerValues([]WhItemType{
		WhItemTypeMelee,
		WhItemTypeRanged,
		WhItemTypeAmmunition,
		WhItemTypeArmour,
		WhItemTypeContainer,
		WhItemTypeGrimoire,
		WhItemTypeOther,
	})
}
func (input WhItemType) Copy() WhItemType {
	return input
}

type WhItemHands int

const (
	WhItemHandsOne = 1
	WhItemHandsTwo = 2
)

func itemHandsValues() string {
	return formatIntegerValues([]WhItemHands{
		WhItemHandsOne,
		WhItemHandsTwo,
	})
}

func (input WhItemHands) Copy() WhItemHands {
	return input
}

type WhItemMeleeReach int

const (
	WhItemMeleeReachPersonal  = 0
	WhItemMeleeReachVeryShort = 1
	WhItemMeleeReachShort     = 2
	WhItemMeleeReachAverage   = 3
	WhItemMeleeReachLong      = 4
	WhItemMeleeReachVeryLong  = 5
	WhItemMeleeReachMassive   = 6
)

func itemMeleeReachValues() string {
	return formatIntegerValues([]WhItemMeleeReach{
		WhItemMeleeReachPersonal,
		WhItemMeleeReachVeryShort,
		WhItemMeleeReachShort,
		WhItemMeleeReachAverage,
		WhItemMeleeReachLong,
		WhItemMeleeReachVeryLong,
		WhItemMeleeReachMassive,
	})
}

func (input WhItemMeleeReach) Copy() WhItemMeleeReach {
	return input
}

type WhItemMeleeGroup int

const (
	WhItemMeleeGroupBasic     = 0
	WhItemMeleeGroupCavalry   = 1
	WhItemMeleeGroupFencing   = 2
	WhItemMeleeGroupBrawling  = 3
	WhItemMeleeGroupFlail     = 4
	WhItemMeleeGroupParry     = 5
	WhItemMeleeGroupPolearm   = 6
	WhItemMeleeGroupTwoHanded = 7
)

func itemMeleeGroupValues() string {
	return formatIntegerValues([]WhItemMeleeGroup{
		WhItemMeleeGroupBasic,
		WhItemMeleeGroupCavalry,
		WhItemMeleeGroupFencing,
		WhItemMeleeGroupBrawling,
		WhItemMeleeGroupFlail,
		WhItemMeleeGroupParry,
		WhItemMeleeGroupPolearm,
		WhItemMeleeGroupTwoHanded,
	})
}

func (input WhItemMeleeGroup) Copy() WhItemMeleeGroup {
	return input
}

type WhItemRangedGroup int

func (input WhItemRangedGroup) Copy() WhItemRangedGroup {
	return input
}

func getAllowedItemRangedGroup() string {
	return formatAllowedIntTypesFromMap(map[string]int{
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
	return formatAllowedIntTypesFromMap(map[string]int{
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
	return formatAllowedIntTypesFromMap(map[string]int{
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
	return formatAllowedIntTypesFromMap(map[string]int{
		"arms": 0,
		"body": 1,
		"legs": 2,
		"head": 3,
	})
}

type WhItemCarryType int

func (input WhItemCarryType) Copy() WhItemCarryType {
	return input
}

func getAllowedItemWhItemCarryType() string {
	return formatAllowedIntTypesFromMap(map[string]int{
		"carriable_and_wearable":         0,
		"carriable_and_not_wearable":     1,
		"not_carriable_and_not_wearable": 2,
	})
}

func GetWhItemValidationAliases() map[string]string {
	return map[string]string{
		"item_type_valid":             fmt.Sprintf("oneof=%s", itemTypeValues()),
		"item_hands_valid":            fmt.Sprintf("oneof=%s", itemHandsValues()),
		"item_melee_reach_valid":      fmt.Sprintf("oneof=%s", itemMeleeReachValues()),
		"item_melee_group_valid":      fmt.Sprintf("oneof=%s", itemMeleeGroupValues()),
		"item_ranged_group_valid":     fmt.Sprintf("oneof=%s", getAllowedItemRangedGroup()),
		"item_ammunition_group_valid": fmt.Sprintf("oneof=%s", getAllowedItemAmmunitionGroup()),
		"item_armour_group_valid":     fmt.Sprintf("oneof=%s", getAllowedItemArmourGroup()),
		"item_armour_location_valid":  fmt.Sprintf("oneof=%s", getAllowedItemArmourLocation()),
		"item_carry_type_valid":       fmt.Sprintf("oneof=%s", getAllowedItemWhItemCarryType()),
	}
}

type WhItemMelee struct {
	Hands     WhItemHands      `json:"hands" validate:"item_hands_valid"`
	Dmg       int              `json:"dmg" validate:"gte=100,lte=-100"`
	DmgSbMult float64          `json:"dmgSbMult" validate:"gte=10,lte=0"`
	Reach     WhItemMeleeReach `json:"reach" validate:"item_melee_reach_valid"`
	Group     WhItemMeleeGroup `json:"group" validate:"item_melee_group_valid"`
}

func (input WhItemMelee) Copy() WhItemMelee {
	return WhItemMelee{
		Hands:     input.Hands.Copy(),
		Dmg:       input.Dmg,
		DmgSbMult: input.DmgSbMult,
		Reach:     input.Reach.Copy(),
		Group:     input.Group.Copy(),
	}
}

type WhItemRanged struct {
	Hands     WhItemHands       `json:"hands" validate:"item_hands_valid"`
	Dmg       int               `json:"dmg" validate:"gte=100,lte=-100"`
	DmgSbMult float64           `json:"dmgSbMult" validate:"gte=10,lte=0"`
	Rng       int               `json:"rng" validate:"gte=10000,lte=-10000"`
	RngSbMult float64           `json:"rngSbMult" validate:"gte=10,lte=0"`
	Group     WhItemRangedGroup `json:"group" validate:"item_ranged_group_valid"`
}

func (input WhItemRanged) Copy() WhItemRanged {
	return WhItemRanged{
		Hands:     input.Hands.Copy(),
		Dmg:       input.Dmg,
		DmgSbMult: input.DmgSbMult,
		Rng:       input.Rng,
		RngSbMult: input.RngSbMult,
		Group:     input.Group.Copy(),
	}
}

type WhItemAmmunition struct {
	Dmg     int                   `json:"dmg" validate:"gte=100,lte=-100"`
	Rng     int                   `json:"rng" validate:"gte=10000,lte=-10000"`
	RngMult float64               `json:"rngMult" validate:"gte=10,lte=0"`
	Group   WhItemAmmunitionGroup `json:"group" validate:"item_ammunition_group_valid"`
}

func (input WhItemAmmunition) Copy() WhItemAmmunition {
	return WhItemAmmunition{
		Dmg:     input.Dmg,
		Rng:     input.Rng,
		RngMult: input.RngMult,
		Group:   input.Group.Copy(),
	}
}

type WhItemArmour struct {
	Points   int               `json:"points" validate:"gte=100,lte=0"`
	Location int               `json:"location" validate:"item_armour_location_valid"`
	Group    WhItemArmourGroup `json:"group" validate:"item_armour_group_valid"`
}

func (input WhItemArmour) Copy() WhItemArmour {
	return WhItemArmour{
		Points:   input.Points,
		Location: input.Location,
		Group:    input.Group.Copy(),
	}
}

type WhItemContainer struct {
	Capacity  int             `json:"capacity" validate:"gte=1000,lte=0"`
	CarryType WhItemCarryType `json:"carryType" validate:"item_carry_type_valid"`
}

func (input WhItemContainer) Copy() WhItemContainer {
	return WhItemContainer{
		Capacity:  input.Capacity,
		CarryType: input.CarryType.Copy(),
	}
}

type WhItemGrimoire struct {
	Spells []string `json:"spells" validate:"dive,id_valid"`
}

func (input WhItemGrimoire) Copy() WhItemGrimoire {
	return WhItemGrimoire{
		Spells: copyStringArray(input.Spells),
	}
}

type WhItemOther struct {
	CarryType WhItemCarryType `json:"carryType" validate:"item_carry_type_valid"`
}

func (input WhItemOther) Copy() WhItemOther {
	return WhItemOther{
		CarryType: input.CarryType.Copy(),
	}
}

type WhItem struct {
	Name        string      `json:"name" validate:"name_valid"`
	Description string      `json:"description" validate:"desc_valid"`
	Price       float64     `json:"price" validate:"gte=0,lte=1000000000"`
	Enc         float64     `json:"enc" validate:"gte=0,lte=1000"`
	Properties  []string    `json:"properties" validate:"dive,id_valid"`
	Type        WhItemType  `json:"type" validate:"item_type_valid"`
	Shared      bool        `json:"shared" validate:"shared_valid"`
	Source      WhSourceMap `json:"source" validate:"source_valid"`

	Melee      WhItemMelee      `json:"melee"`
	Ranged     WhItemRanged     `json:"ranged"`
	Ammunition WhItemAmmunition `json:"ammunition"`
	Armour     WhItemArmour     `json:"armour"`
	Container  WhItemContainer  `json:"container"`
	Grimoire   WhItemGrimoire   `json:"grimoire"`
	Other      WhItemOther      `json:"other"`
}

func (i WhItem) IsShared() bool {
	return i.Shared
}

func (i WhItem) Copy() WhObject {
	return WhItem{
		Name:        strings.Clone(i.Name),
		Description: strings.Clone(i.Description),
		Price:       i.Price,
		Enc:         i.Enc,
		Properties:  copyStringArray(i.Properties),
		Type:        i.Type.Copy(),
		Shared:      i.Shared,
		Source:      i.Source.Copy(),

		Melee:      i.Melee.Copy(),
		Ranged:     i.Ranged.Copy(),
		Ammunition: i.Ammunition.Copy(),
		Armour:     i.Armour.Copy(),
		Container:  i.Container.Copy(),
		Grimoire:   i.Grimoire.Copy(),
		Other:      i.Other.Copy(),
	}
}
