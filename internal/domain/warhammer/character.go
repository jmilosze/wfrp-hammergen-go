package warhammer

import (
	"fmt"
	"strings"
)

type WhCharacterSpecies string

const (
	WhCharacterSpeciesHumanDefault             = "0000"
	WhCharacterSpeciesHumanReikland            = "0001"
	WhCharacterSpeciesHumanAltdorfSouthBank    = "0002"
	WhCharacterSpeciesHumanAltdorfEastend      = "0003"
	WhCharacterSpeciesHumanAltdorfHexxerbezrik = "0004"
	WhCharacterSpeciesHumanAltdorfDocklands    = "0005"
	WhCharacterSpeciesHumanMiddenheim          = "0006"
	WhCharacterSpeciesHumanMiddenland          = "0007"
	WhCharacterSpeciesHumanNordland            = "0008"
	WhCharacterSpeciesHumanSalzenmund          = "0009"
	WhCharacterSpeciesHumanTilea               = "0010"
	WhCharacterSpeciesHumanNorseBjornling      = "0011"
	WhCharacterSpeciesHumanNorseSarl           = "0012"
	WhCharacterSpeciesHumanNorseSkaeling       = "0013"
	WhCharacterSpeciesHalflingDefault          = "0100"
	WhCharacterSpeciesHalflingAshfield         = "0101"
	WhCharacterSpeciesHalflingBrambledown      = "0102"
	WhCharacterSpeciesHalflingBrandysnap       = "0103"
	WhCharacterSpeciesHalflingHayfoot          = "0104"
	WhCharacterSpeciesHalflingHollyfoot        = "0105"
	WhCharacterSpeciesHalflingHayfootHollyfoot = "0106"
	WhCharacterSpeciesHalflingLostpockets      = "0107"
	WhCharacterSpeciesHalflingLowhaven         = "0108"
	WhCharacterSpeciesHalflingRumster          = "0109"
	WhCharacterSpeciesHalflingSkelfsider       = "0110"
	WhCharacterSpeciesHalflingThorncobble      = "0111"
	WhCharacterSpeciesHalflingTumbleberry      = "0112"
	WhCharacterSpeciesDwarfDefault             = "0200"
	WhCharacterSpeciesDwarfAltdorf             = "0201"
	WhCharacterSpeciesDwarfCragforgeClan       = "0202"
	WhCharacterSpeciesDwarfGrumssonClan        = "0203"
	WhCharacterSpeciesDwarfNorse               = "0204"
	WhCharacterSpeciesHighElfDefault           = "0300"
	WhCharacterSpeciesWoodElfDefault           = "0400"
	WhCharacterSpeciesGnomeDefault             = "0500"
	WhCharacterSpeciesOgreDefault              = "0600"
)

func characterSpeciesValues() string {
	return formatStringValues([]WhCharacterSpecies{
		WhCharacterSpeciesHumanDefault,
		WhCharacterSpeciesHumanReikland,
		WhCharacterSpeciesHumanAltdorfSouthBank,
		WhCharacterSpeciesHumanAltdorfEastend,
		WhCharacterSpeciesHumanAltdorfHexxerbezrik,
		WhCharacterSpeciesHumanAltdorfDocklands,
		WhCharacterSpeciesHumanMiddenheim,
		WhCharacterSpeciesHumanMiddenland,
		WhCharacterSpeciesHumanNordland,
		WhCharacterSpeciesHumanSalzenmund,
		WhCharacterSpeciesHumanTilea,
		WhCharacterSpeciesHumanNorseBjornling,
		WhCharacterSpeciesHumanNorseSarl,
		WhCharacterSpeciesHumanNorseSkaeling,
		WhCharacterSpeciesHalflingDefault,
		WhCharacterSpeciesHalflingAshfield,
		WhCharacterSpeciesHalflingBrambledown,
		WhCharacterSpeciesHalflingBrandysnap,
		WhCharacterSpeciesHalflingHayfoot,
		WhCharacterSpeciesHalflingHollyfoot,
		WhCharacterSpeciesHalflingHayfootHollyfoot,
		WhCharacterSpeciesHalflingLostpockets,
		WhCharacterSpeciesHalflingLowhaven,
		WhCharacterSpeciesHalflingRumster,
		WhCharacterSpeciesHalflingSkelfsider,
		WhCharacterSpeciesHalflingThorncobble,
		WhCharacterSpeciesHalflingTumbleberry,
		WhCharacterSpeciesDwarfDefault,
		WhCharacterSpeciesDwarfAltdorf,
		WhCharacterSpeciesDwarfCragforgeClan,
		WhCharacterSpeciesDwarfGrumssonClan,
		WhCharacterSpeciesDwarfNorse,
		WhCharacterSpeciesHighElfDefault,
		WhCharacterSpeciesWoodElfDefault,
		WhCharacterSpeciesGnomeDefault,
		WhCharacterSpeciesOgreDefault,
	})
}
func (input WhCharacterSpecies) InitAndCopy() WhCharacterSpecies {
	if input == "" {
		return WhCharacterSpeciesHumanDefault
	}
	return WhCharacterSpecies(strings.Clone(string(input)))
}

type WhIdNumberMap map[string]int

func (input WhIdNumberMap) InitAndCopy() WhIdNumberMap {
	output := make(WhIdNumberMap, len(input))
	for key, value := range input {
		output[key] = value
	}
	return output
}

type WhCharacter struct {
	Name              string             `json:"name" validate:"name_valid"`
	Description       string             `json:"description" validate:"desc_valid"`
	Notes             string             `json:"notes" validate:"desc_valid"`
	EquippedItems     WhIdNumberMap      `json:"equippedItems" validate:"id_number_map_valid"`
	CarriedItems      WhIdNumberMap      `json:"carriedItems" validate:"id_number_map_valid"`
	StoredItems       WhIdNumberMap      `json:"storedItems" validate:"id_number_map_valid"`
	Skills            WhIdNumberMap      `json:"skills" validate:"id_number_map_valid"`
	Talents           WhIdNumberMap      `json:"talents" validate:"id_number_map_valid"`
	Species           WhCharacterSpecies `json:"species" validate:"character_species_valid"`
	BaseAttributes    WhAttributes
	AttributeAdvances WhAttributes
	CareerPath        []string   `json:"careerPath" validate:"dive,id_valid"`
	Career            string     `json:"career" validate:"id_valid"`
	Fate              int        `json:"fate" validate:"gte=0,lte=1000"`
	Fortune           int        `json:"fortune" validate:"gte=0,lte=1000"`
	Resilience        int        `json:"resilience" validate:"gte=0,lte=1000"`
	Resolve           int        `json:"resolve" validate:"gte=0,lte=1000"`
	CurrentExp        int        `json:"currentExp" validate:"gte=0,lte=10000000"`
	SpentExp          int        `json:"spentExp" validate:"gte=0,lte=10000000"`
	Status            WhStatus   `json:"status" validate:"status_valid"`
	Standing          WhStanding `json:"standing" validate:"standing_valid"`
	Brass             int        `json:"brass" validate:"gte=0,lte=1000000"`
	Silver            int        `json:"silver" validate:"gte=0,lte=1000000"`
	Gold              int        `json:"gold" validate:"gte=0,lte=1000000"`
	Spells            []string   `json:"spells" validate:"dive,id_valid"`
	Sin               int        `json:"sin" validate:"gte=0,lte=1000"`
	Corruption        int        `json:"corruption" validate:"gte=0,lte=1000"`
	Mutations         []string   `json:"mutations" validate:"dive,id_valid"`
	Shared            bool       `json:"shared" validate:"shared_valid"`
}

func (c WhCharacter) IsShared() bool {
	return c.Shared
}

func (c WhCharacter) InitAndCopy() WhObject {
	return WhCharacter{
		Name:              strings.Clone(c.Name),
		Description:       strings.Clone(c.Description),
		Notes:             strings.Clone(c.Notes),
		EquippedItems:     c.EquippedItems.InitAndCopy(),
		CarriedItems:      c.CarriedItems.InitAndCopy(),
		StoredItems:       c.StoredItems.InitAndCopy(),
		Skills:            c.Skills.InitAndCopy(),
		Talents:           c.Talents.InitAndCopy(),
		Species:           c.Species.InitAndCopy(),
		BaseAttributes:    c.BaseAttributes.InitAndCopy(),
		AttributeAdvances: c.AttributeAdvances.InitAndCopy(),
		CareerPath:        copyStringArray(c.CareerPath),
		Career:            strings.Clone(c.Career),
		Fate:              c.Fate,
		Fortune:           c.Fortune,
		Resilience:        c.Resilience,
		Resolve:           c.Resolve,
		CurrentExp:        c.CurrentExp,
		SpentExp:          c.SpentExp,
		Status:            c.Status.InitAndCopy(),
		Standing:          c.Standing.InitAndCopy(),
		Brass:             c.Brass,
		Silver:            c.Silver,
		Gold:              c.Gold,
		Spells:            copyStringArray(c.Spells),
		Sin:               c.Sin,
		Corruption:        c.Corruption,
		Mutations:         copyStringArray(c.Mutations),
		Shared:            c.Shared,
	}
}

func GetWhCharacterValidationAliases() map[string]string {
	return map[string]string{
		"character_species_valid": fmt.Sprintf("oneof=%s", characterSpeciesValues()),
		"id_number_map_valid":     "dive,keys,id_valid,endkeys,gte=1,lte=1000",
	}
}