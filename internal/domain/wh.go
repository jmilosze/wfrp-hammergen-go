package domain

import (
	"context"
	"fmt"
	"strings"
)

const (
	WhTypeMutation = 0
	WhTypeSpell    = 1

	MutationPhysical = 0
	MutationMental   = 1
)

const (
	WhInvalidArgumentsError = iota
)

type WhError struct {
	WhType  int
	ErrType int
	Err     error
}

func (e *WhError) Unwrap() error {
	return e.Err
}

func (e *WhError) Error() string {
	return fmt.Sprintf("wh error, %s", e.Err)
}

type WhType interface {
	Mutation | Spell
}

type Mutation struct {
	Name        string `json:"name" validate:"omitempty,min=0,max=200,excludesall=<>"`
	Description string `json:"description" validate:"omitempty,min=0,max=100000,excludesall=<>"`
	Type        *int   `json:"type" validate:"omitempty,oneof=0 1"`
	Shared      *bool  `json:"shared" validate:"omitempty"`
	Id          string `json:"id,omitempty"`
	OwnerId     string `json:"owner_id,omitempty"`
}

func CopyMutation(from *Mutation, to *Mutation) {
	if from == nil {
		to = nil
	}

	to = &Mutation{
		Name:        strings.Clone(from.Name),
		Description: strings.Clone(from.Description),
		Type:        *(&from.Type),
		Shared:      *(&from.Shared),
		Id:          strings.Clone(from.Id),
		OwnerId:     strings.Clone(from.OwnerId),
	}
}

type Spell struct {
	Name        string `json:"name" validate:"omitempty,min=0,max=200,excludesall=<>"`
	Description string `json:"description" validate:"omitempty,min=0,max=100000,excludesall=<>"`
	Cn          *int   `json:"cn" validate:"omitempty,min=-1,max=99"`
	Range       string `json:"range" validate:"omitempty,min=0,max=200,excludesall=<>"`
	Target      string `json:"target" validate:"omitempty,min=0,max=200,excludesall=<>"`
	Duration    string `json:"duration" validate:"omitempty,min=0,max=200,excludesall=<>"`
	Shared      *bool  `json:"shared" validate:"omitempty"`
	Id          string `json:"id,omitempty"`
	OwnerId     string `json:"owner_id,omitempty"`
}

func CopySpell(from *Spell, to *Spell) {
	if from == nil {
		to = nil
	}

	to = &Spell{
		Name:        strings.Clone(from.Name),
		Description: strings.Clone(from.Description),
		Cn:          *(&from.Cn),
		Range:       strings.Clone(from.Range),
		Target:      strings.Clone(from.Target),
		Duration:    strings.Clone(from.Duration),
		Shared:      *(&from.Shared),
		Id:          strings.Clone(from.Id),
		OwnerId:     strings.Clone(from.OwnerId),
	}
}

func WhSetOwnerId[W WhType](wh *W, ownerId string) error {
	switch v := any(wh).(type) {
	case *Mutation:
		v.OwnerId = ownerId
	case *Spell:
		v.OwnerId = ownerId
	default:
		return fmt.Errorf("could not set OwnerId on type %T", v)
	}
	return nil
}

func WhCopy[W WhType](from *W, to *W) error {
	switch v := any(from).(type) {
	case *Mutation:
		CopyMutation(v, any(to).(*Mutation))
	case *Spell:
		CopySpell(v, any(to).(*Spell))
	default:
		return fmt.Errorf("could not copy type %T", v)
	}
	return nil
}

type WHService[W WhType] interface {
	Create(ctx context.Context, whWrite *W, c *Claims) (*W, *WhError)
}

type WhDbService[W WhType] interface {
	Create(ctx context.Context, wh *W) (*W, *DbError)
}
