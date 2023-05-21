package warhammer

import "strings"

type WhMutation struct {
	Name        string      `json:"name" validate:"name_valid"`
	Description string      `json:"description" validate:"desc_valid"`
	Type        int         `json:"type" validate:"oneof=0 1"`
	Modifiers   WhModifiers `json:"modifiers"`
	Shared      bool        `json:"shared" validate:"shared_valid"`
	Source      WhSource    `json:"source" validate:"source_valid"`
}

func (m WhMutation) IsShared() bool {
	return m.Shared
}

func (m WhMutation) Copy() WhObject {
	return WhMutation{
		Name:        strings.Clone(m.Name),
		Description: strings.Clone(m.Description),
		Type:        m.Type,
		Shared:      m.Shared,
		Modifiers:   m.Modifiers.Copy(),
	}
}
