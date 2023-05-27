package warhammer

import (
	"fmt"
	"strconv"
	"strings"
)

type WhSource string

const (
	WhSourceCustom                    = "0"
	WhSourceWFRP                      = "1"
	WhSourceRoughNightsAndHardDays    = "2"
	WhSourceArchivesOfTheEmpireVolI   = "3"
	WhSourceArchivesOfTheEmpireVolII  = "4"
	WhSourceArchivesOfTheEmpireVolIII = "5"
	WhSourceUpInArms                  = "6"
	WhSourceWindsOfMagic              = "7"
	WhSourceMiddenheim                = "8"
	WhSourceSalzenmund                = "9"
	WhSourceSeaOfClaws                = "10"
	WhSourceLustria                   = "11"
)

type WhSourceMap map[WhSource]string

func (input WhSourceMap) Copy() WhSourceMap {
	output := make(WhSourceMap, len(input))
	for key, value := range input {
		output[key] = value
	}
	return output
}

func sourceValues() string {
	return formatStringValues([]WhSource{
		WhSourceCustom,
		WhSourceWFRP,
		WhSourceRoughNightsAndHardDays,
		WhSourceArchivesOfTheEmpireVolI,
		WhSourceArchivesOfTheEmpireVolII,
		WhSourceArchivesOfTheEmpireVolIII,
		WhSourceUpInArms,
		WhSourceWindsOfMagic,
		WhSourceMiddenheim,
		WhSourceSalzenmund,
		WhSourceSeaOfClaws,
		WhSourceLustria,
	})
}

func GetWhCommonValidationAliases() map[string]string {
	return map[string]string{
		"name_valid":          "min=0,max=200,excludesall=<>",
		"desc_valid":          "min=0,max=100000,excludesall=<>",
		"shared_valid":        "boolean",
		"medium_string_valid": "min=0,max=200,excludesall=<>",
		"source_valid":        fmt.Sprintf("dive,keys,oneof=%s,endkeys,min=0,max=15,excludesall=<>", sourceValues()),
		"id_valid":            "hexadecimal,len=24",
	}
}

func formatAllowedIntTypesFromMap[T ~int](list map[string]T) string {
	values := make([]string, len(list))
	for _, v := range list {
		values = append(values, strconv.Itoa(int(v)))
	}
	return strings.Join(values, " ")
}

func formatIntegerValues[T ~int](list []T) string {
	values := make([]string, len(list))
	for _, v := range list {
		values = append(values, strconv.Itoa(int(v)))
	}
	return strings.Join(values, " ")
}

func formatStringValues[T ~string](list []T) string {
	values := make([]string, len(list))
	for _, v := range list {
		values = append(values, string(v))
	}
	return strings.Join(values, " ")
}

func copyStringArray(input []string) []string {
	output := make([]string, len(input))
	for i, v := range input {
		output[i] = strings.Clone(v)
	}
	return output
}
