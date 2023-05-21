package warhammer

import "strings"

const (
	WhTypeMutation = "mutation"
	WhTypeSpell    = "spell"
	WhTypeProperty = "property"
)

type WhType string

type WhSourceType string

type WhSource map[WhSourceType]string

func (input WhSource) Copy() WhSource {
	output := make(WhSource, len(input))
	for key, value := range input {
		output[key] = value
	}
	return output
}

func getAllowedSourceType() string {
	var list = map[string]WhSourceType{
		"custom":                       "0",
		"wfrp":                         "1",
		"rough_nights_and_hard_days":   "2",
		"archives_of_the_empire_vol_1": "3",
		"archives_of_the_empire_vol_2": "4",
		"archives_of_the_empire_vol_3": "5",
		"up_in_arms":                   "6",
		"winds_of_magic":               "7",
		"middenheim":                   "8",
		"salzenmund":                   "9",
		"sea_of_claws":                 "10",
		"lustria":                      "11",
	}

	values := make([]string, 0, len(list))
	for _, v := range list {
		values = append(values, string(v))
	}
	return strings.Join(values, " ")
}
