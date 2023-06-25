package warhammer

import (
	"fmt"
	"strings"
)

type WhIdNumberMap map[string]int

func (input WhIdNumberMap) InitAndCopy() WhIdNumberMap {
	output := make(WhIdNumberMap, len(input))
	for key, value := range input {
		output[key] = value
	}
	return output
}

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
	for k1, v1 := range input.Multiple {
		multiple[k1] = make([]string, len(v1))
		for k2, v2 := range v1 {
			multiple[k1][k2] = strings.Clone(v2)
		}
	}

	return WhSpeciesTalents{Single: single, Multiple: multiple}
}

type WhGenerationProps struct {
	Name           string
	ClassItems     map[WhCareerClass]WhItems               `json:"classItems"`
	RandomTalents  []WhRandomTalent                        `json:"randomTalents"`
	SpeciesTalents map[WhCharacterSpecies]WhSpeciesTalents `json:"speciesTalents"`
	SpeciesSkills  map[WhCharacterSpecies][]string         `json:"speciesSkills"`
}

func (gprops WhGenerationProps) InitAndCopy() WhGenerationProps {
	classItems := make(map[WhCareerClass]WhItems, len(gprops.ClassItems))
	for k, v := range gprops.ClassItems {
		classItems[k] = v.InitAndCopy()
	}

	randomTalents := make([]WhRandomTalent, len(gprops.RandomTalents))
	for k, v := range gprops.RandomTalents {
		randomTalents[k] = v.InitAndCopy()
	}

	speciesTalents := make(map[WhCharacterSpecies]WhSpeciesTalents, len(gprops.SpeciesTalents))
	for k, v := range gprops.SpeciesTalents {
		speciesTalents[k] = v.InitAndCopy()
	}

	speciesSkills := make(map[WhCharacterSpecies][]string, len(gprops.SpeciesSkills))
	for k1, v1 := range gprops.SpeciesSkills {
		skills := make([]string, len(v1))
		for k2, v2 := range v1 {
			skills[k2] = strings.Clone(v2)
		}
		speciesSkills[k1] = skills
	}

	return WhGenerationProps{
		Name:           strings.Clone(gprops.Name),
		ClassItems:     classItems,
		RandomTalents:  randomTalents,
		SpeciesTalents: speciesTalents,
		SpeciesSkills:  speciesSkills,
	}
}

func (gprops WhGenerationProps) ToMap() (map[string]any, error) {
	gMap, err := structToMap(gprops)
	if err != nil {
		return map[string]any{}, fmt.Errorf("error while mapping wh structure %s", err)
	}
	return gMap, nil
}

func GetWhGenerationPropsValidationAliases() map[string]string {
	return map[string]string{
		"id_number_map_valid": "dive,keys,id_valid,endkeys,gte=1,lte=1000",
	}
}
