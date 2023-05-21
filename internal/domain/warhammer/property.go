package warhammer

import (
	"strconv"
	"strings"
)

type WhProperty struct {
	Name         string         `json:"name" validate:"name_valid"`
	Description  string         `json:"description" validate:"desc_valid"`
	Type         WhPropertyType `json:"type" validate:"property_type_valid"`
	ApplicableTo []WhItemType   `json:"applicableTo" validate:"property_applicable_to_valid"`
	Shared       bool           `json:"shared" validate:"shared_valid"`
	Source       WhSource       `json:"source" validate:"source_valid"`
}

func (p WhProperty) IsShared() bool {
	return p.Shared
}

func (p WhProperty) Copy() WhObject {
	return WhProperty{
		Name:         strings.Clone(p.Name),
		Description:  strings.Clone(p.Description),
		Type:         p.Type.Copy(),
		ApplicableTo: copyApplicableTo(p.ApplicableTo),
		Shared:       p.Shared,
		Source:       p.Source.Copy(),
	}
}

func copyApplicableTo(at []WhItemType) []WhItemType {
	r := make([]WhItemType, len(at))
	for i, v := range at {
		r[i] = v.Copy()
	}
	return r
}

type WhPropertyType int

func (input WhPropertyType) Copy() WhPropertyType {
	return input
}

func getAllowedPropertyType() string {
	var list = map[string]WhItemType{
		"quality": 0,
		"flaw":    1,
	}

	values := make([]string, 0, len(list))
	for _, v := range list {
		values = append(values, strconv.Itoa(int(v)))
	}
	return strings.Join(values, " ")
}
