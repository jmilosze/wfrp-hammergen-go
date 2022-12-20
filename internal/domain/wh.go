package domain

import (
	"context"
	"fmt"
	"strings"
)

const (
	WhTypeMutation = 1
	WhTypeSpell    = 2

	MutationPhysical = 0
	MutationMental   = 1
)

const (
	WhInvalidArgumentsError = iota
	WhNotFoundError
	WhInternalError
	WhUnauthorizedError
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

type Wh struct {
	Shared  bool     `json:"shared" validate:"omitempty"`
	Id      string   `json:"id"`
	OwnerId string   `json:"ownerId"`
	CanEdit bool     `json:"canEdit"`
	Object  WhObject `json:"object"`
}

func (w *Wh) Copy() *Wh {
	if w == nil {
		return nil
	}

	return &Wh{
		Shared:  w.Shared,
		Id:      strings.Clone(w.Id),
		OwnerId: strings.Clone(w.OwnerId),
		Object:  w.Object.Copy(),
	}
}

type WhObject interface {
	Copy() WhObject
}

type WhMutation struct {
	Name        string `json:"name" validate:"omitempty,min=0,max=200,excludesall=<>"`
	Description string `json:"description" validate:"omitempty,min=0,max=100000,excludesall=<>"`
	Type        int    `json:"type" validate:"omitempty,oneof=0 1"`
}

func (m WhMutation) Copy() WhObject {
	return WhMutation{
		Name:        strings.Clone(m.Name),
		Description: strings.Clone(m.Description),
		Type:        m.Type,
	}
}

type WhSpell struct {
	Name        string `json:"name" validate:"omitempty,min=0,max=200,excludesall=<>"`
	Description string `json:"description" validate:"omitempty,min=0,max=100000,excludesall=<>"`
	Cn          int    `json:"cn" validate:"omitempty,min=-1,max=99"`
	Range       string `json:"range" validate:"omitempty,min=0,max=200,excludesall=<>"`
	Target      string `json:"target" validate:"omitempty,min=0,max=200,excludesall=<>"`
	Duration    string `json:"duration" validate:"omitempty,min=0,max=200,excludesall=<>"`
}

func (s WhSpell) Copy() WhObject {
	return WhSpell{
		Name:        strings.Clone(s.Name),
		Description: strings.Clone(s.Description),
		Cn:          s.Cn,
		Range:       strings.Clone(s.Range),
		Target:      strings.Clone(s.Target),
		Duration:    strings.Clone(s.Duration),
	}
}

type WhService interface {
	Create(ctx context.Context, whType int, w *Wh, c *Claims) (*Wh, *WhError)
	Get(ctx context.Context, whType int, whId string, c *Claims) (*Wh, *WhError)
}
