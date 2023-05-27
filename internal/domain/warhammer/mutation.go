package warhammer

import (
	"fmt"
	"strings"
)

type WhMutation struct {
	Name        string         `json:"name" validate:"name_valid"`
	Description string         `json:"description" validate:"desc_valid"`
	Type        WhMutationType `json:"type" validate:"mutation_type_valid"`
	Modifiers   WhModifiers    `json:"modifiers"`
	Shared      bool           `json:"shared" validate:"shared_valid"`
	Source      WhSourceMap    `json:"source" validate:"source_valid"`
}

func (m WhMutation) IsShared() bool {
	return m.Shared
}

func (m WhMutation) Copy() WhObject {
	return WhMutation{
		Name:        strings.Clone(m.Name),
		Description: strings.Clone(m.Description),
		Type:        m.Type.Copy(),
		Modifiers:   m.Modifiers.Copy(),
		Shared:      m.Shared,
		Source:      m.Source.Copy(),
	}
}

type WhMutationType int

const (
	WhMutationTypePhysical = 0
	WhMutationTypeMental   = 1
)

func mutationTypeValues() string {
	return formatIntegerValues([]WhMutationType{WhMutationTypePhysical, WhMutationTypeMental})
}

func (input WhMutationType) Copy() WhMutationType {
	return input
}

func GetWhMutationValidationAliases() map[string]string {
	return map[string]string{
		"mutation_type_valid": fmt.Sprintf("oneof=%s", mutationTypeValues()),
	}
}
