package domain

import "context"

const (
	WhTypeMutation = 0
	WhTypeSpell    = 1

	MutationPhysical = 0
	MutationMental   = 1
)

type WhError struct {
	WhType  int
	ErrType int
	Err     error
}

type WhTypePointer interface {
	*Mutation | *Spell
	SetOwnerId(string)
}
type Mutation struct {
	Name        string `json:"name" validate:"omitempty,min=0,max=200,excludesall=<>"`
	Description string `json:"description" validate:"omitempty,min=0,max=100000,excludesall=<>"`
	Shared      *bool  `json:"shared" validate:"omitempty"`
	Type        *int   `json:"type" validate:"omitempty,oneof 0 1"`
	Id          string `json:"id,omitempty"`
	OwnerId     string `json:"owner_id,omitempty"`
}

func (m *Mutation) SetOwnerId(ownerId string) {
	m.OwnerId = ownerId
}

type Spell struct {
	Name        string `json:"name" validate:"omitempty,min=0,max=200,excludesall=<>"`
	Cn          *int   `json:"cn" validate:"omitempty,min=-1,max=99"`
	Range       string `json:"range" validate:"omitempty,min=0,max=200,excludesall=<>"`
	Target      string `json:"target" validate:"omitempty,min=0,max=200,excludesall=<>"`
	Duration    string `json:"duration" validate:"omitempty,min=0,max=200,excludesall=<>"`
	Description string `json:"description" validate:"omitempty,min=0,max=100000,excludesall=<>"`
	Id          string `json:"id,omitempty"`
	OwnerId     string `json:"owner_id,omitempty"`
}

func (s *Spell) SetOwnerId(ownerId string) {
	s.OwnerId = ownerId
}

type WHService[W WhTypePointer] interface {
	Create(ctx context.Context, whWrite W) (W, *WhError)
}
