package warhammer

import "strings"

type WhItems struct {
	Equipped WhIdNumberMap `json:"equipped"`
	Carried  WhIdNumberMap `json:"carried"`
	Stored   WhIdNumberMap `json:"stored"`
}

func (input WhItems) InitAndCopy() WhItems {
	return WhItems{
		Equipped: input.Equipped.InitAndCopy(),
		Carried:  input.Carried.InitAndCopy(),
		Stored:   input.Stored.InitAndCopy(),
	}
}

type WhRandomTalent struct {
	Id      string `json:"id"`
	MinRoll int    `json:"minRoll"`
	MaxRoll int    `json:"maxRoll"`
}

func (input WhRandomTalent) InitAndCopy() WhRandomTalent {
	return WhRandomTalent{
		Id:      strings.Clone(input.Id),
		MinRoll: input.MinRoll,
		MaxRoll: input.MaxRoll,
	}
}

type WhSpeciesTalents struct {
	Single   []string   `json:"single"`
	Multiple [][]string `json:"multiple"`
}

func (input WhSpeciesTalents) InitAndCopy() WhSpeciesTalents {
	single := make([]string, len(input.Single))
	for k, v := range input.Single {
		single[k] = strings.Clone(v)
	}

	multiple := make([][]string, len(input.Multiple))
	for outerK, outerV := range input.Multiple {
		multiple[outerK] = make([]string, len(outerV))
		for k, v := range outerV {
			single[k] = strings.Clone(v)
		}
	}

	return WhSpeciesTalents{Single: single, Multiple: multiple}
}

type WhGenerationProp struct {
	ClassItems     map[WhCareerClass]WhItems               `json:"classItems"`
	RandomTalents  []WhRandomTalent                        `json:"randomTalents"`
	SpeciesTalents map[WhCharacterSpecies]WhSpeciesTalents `json:"speciesTalents"`
	SpeciesSkills  map[WhCharacterSpecies]string           `json:"speciesSkills"`
}

func (input *WhGenerationProp) InitAndCopy() *WhGenerationProp {

	if input == nil {
		return nil
	}

	classItems := make(map[WhCareerClass]WhItems, len(input.ClassItems))
	for k, v := range input.ClassItems {
		classItems[k] = v.InitAndCopy()
	}

	randomTalents := make([]WhRandomTalent, len(input.RandomTalents))
	for k, v := range input.RandomTalents {
		randomTalents[k] = v.InitAndCopy()
	}

	speciesTalents := make(map[WhCharacterSpecies]WhSpeciesTalents, len(input.SpeciesTalents))
	for k, v := range input.SpeciesTalents {
		speciesTalents[k] = v.InitAndCopy()
	}

	speciesSkills := make(map[WhCharacterSpecies]string, len(input.SpeciesSkills))
	for k, v := range input.SpeciesSkills {
		speciesSkills[k] = strings.Clone(v)
	}

	return &WhGenerationProp{
		ClassItems:     classItems,
		RandomTalents:  randomTalents,
		SpeciesTalents: speciesTalents,
		SpeciesSkills:  speciesSkills,
	}

}
