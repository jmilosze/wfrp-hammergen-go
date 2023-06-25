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

type WhIdNumber struct {
	Id     string `json:"id" validate:"id_valid"`
	Number int    `json:"number" validate:"gte=1,lte=1000"`
}

func (input WhIdNumber) InitAndCopy() WhIdNumber {
	return WhIdNumber{
		Id:     strings.Clone(input.Id),
		Number: input.Number,
	}
}

func copyArrayIdNumber(input []WhIdNumber) []WhIdNumber {
	output := make([]WhIdNumber, len(input))
	for i, v := range input {
		output[i] = v.InitAndCopy()
	}
	return output
}

type WhCharacter struct {
	Name              string             `json:"name" validate:"name_valid"`
	Description       string             `json:"description" validate:"desc_valid"`
	Notes             string             `json:"notes" validate:"desc_valid"`
	EquippedItems     []WhIdNumber       `json:"equippedItems" validate:"dive"`
	CarriedItems      []WhIdNumber       `json:"carriedItems" validate:"dive"`
	StoredItems       []WhIdNumber       `json:"storedItems" validate:"dive"`
	Skills            []WhIdNumber       `json:"skills" validate:"dive"`
	Talents           []WhIdNumber       `json:"talents" validate:"dive"`
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
		EquippedItems:     copyArrayIdNumber(c.EquippedItems),
		CarriedItems:      copyArrayIdNumber(c.CarriedItems),
		StoredItems:       copyArrayIdNumber(c.StoredItems),
		Skills:            copyArrayIdNumber(c.Skills),
		Talents:           copyArrayIdNumber(c.Talents),
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
	}
}

type WhItemFullNumber struct {
	Item   *WhItemFull `json:"item"`
	Number int         `json:"number"`
}

type WhSkillNumber struct {
	Item   *WhSkill `json:"skill"`
	Number int      `json:"number"`
}

type WhTalentNumber struct {
	Item   *WhTalent `json:"talent"`
	Number int       `json:"number"`
}

type WhCharacterFull struct {
	Name              string             `json:"name"`
	Description       string             `json:"description"`
	Notes             string             `json:"notes"`
	EquippedItems     []WhItemFullNumber `json:"equippedItems"`
	CarriedItems      []WhItemFullNumber `json:"carriedItems"`
	StoredItems       []WhItemFullNumber `json:"storedItems"`
	Skills            []WhSkillNumber    `json:"skills"`
	Talents           []WhTalentNumber   `json:"talents"`
	Species           WhCharacterSpecies `json:"species"`
	BaseAttributes    WhAttributes
	AttributeAdvances WhAttributes
	CareerPath        []*WhCareer   `json:"careerPath"`
	Career            *WhCareer     `json:"career"`
	Fate              int           `json:"fate"`
	Fortune           int           `json:"fortune"`
	Resilience        int           `json:"resilience"`
	Resolve           int           `json:"resolve"`
	CurrentExp        int           `json:"currentExp"`
	SpentExp          int           `json:"spentExp"`
	Status            WhStatus      `json:"status"`
	Standing          WhStanding    `json:"standing"`
	Brass             int           `json:"brass"`
	Silver            int           `json:"silver"`
	Gold              int           `json:"gold"`
	Spells            []*WhSpell    `json:"spells"`
	Sin               int           `json:"sin"`
	Corruption        int           `json:"corruption"`
	Mutations         []*WhMutation `json:"mutations"`
	Shared            bool          `json:"shared"`
}
