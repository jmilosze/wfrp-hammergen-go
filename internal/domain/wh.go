package domain

import (
	"context"
	"fmt"
	"strings"
)

const (
	WhTypeMutation = 1
	WhTypeSpell    = 2
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
	Id      string
	OwnerId string
	CanEdit bool
	Object  WhObject
}

func NewWh(whType int) (Wh, error) {
	var wh Wh

	switch whType {
	case WhTypeMutation:
		wh.Object = &WhMutation{}
	case WhTypeSpell:
		wh.Object = &WhSpell{}
	default:
		return wh, fmt.Errorf("invalid Wh type %d", whType)
	}

	return wh, nil
}

func (w *Wh) Copy() *Wh {
	if w == nil {
		return nil
	}

	return &Wh{
		Id:      strings.Clone(w.Id),
		OwnerId: strings.Clone(w.OwnerId),
		Object:  w.Object.Copy(),
	}
}

func (w *Wh) IsShared() bool {
	return w.Object.IsShared()
}

type WhObject interface {
	Copy() WhObject
	IsShared() bool
}

type WhModifiers struct {
	Size       int `json:"size" validate:"min=-3,max=3"`
	Movement   int `json:"movement" validate:"min=-3,max=3"`
	Attributes WHAttributes
}

type WHAttributes struct {
	WS  int `json:"WS" validate:"min=-99,max=99"`
	BS  int `json:"BS" validate:"min=-99,max=99"`
	S   int `json:"S" validate:"min=-99,max=99"`
	T   int `json:"T" validate:"min=-99,max=99"`
	I   int `json:"I" validate:"min=-99,max=99"`
	Ag  int `json:"Ag" validate:"min=-99,max=99"`
	Dex int `json:"Dex" validate:"min=-99,max=99"`
	Int int `json:"Int" validate:"min=-99,max=99"`
	WP  int `json:"WP" validate:"min=-99,max=99"`
	Fel int `json:"Fel" validate:"min=-99,max=99"`
}

type WhMutation struct {
	Name        string      `json:"name" validate:"min=0,max=200,excludesall=<>"`
	Description string      `json:"description" validate:"min=0,max=100000,excludesall=<>"`
	Type        int         `json:"type" validate:"oneof=0 1"`
	Modifiers   WhModifiers `json:"modifiers"`
	Shared      bool        `json:"shared" validate:"boolean"`
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
	}
}

type WhSpell struct {
	Name        string `json:"name" validate:"min=0,max=200,excludesall=<>"`
	Description string `json:"description" validate:"min=0,max=100000,excludesall=<>"`
	Cn          int    `json:"cn" validate:"min=-1,max=99"`
	Range       string `json:"range" validate:"min=0,max=200,excludesall=<>"`
	Target      string `json:"target" validate:"min=0,max=200,excludesall=<>"`
	Duration    string `json:"duration" validate:"min=0,max=200,excludesall=<>"`
	Shared      bool   `json:"shared" validate:"boolean"`
}

func (s WhSpell) IsShared() bool {
	return s.Shared
}

func (s WhSpell) Copy() WhObject {
	return WhSpell{
		Name:        strings.Clone(s.Name),
		Description: strings.Clone(s.Description),
		Cn:          s.Cn,
		Range:       strings.Clone(s.Range),
		Target:      strings.Clone(s.Target),
		Duration:    strings.Clone(s.Duration),
		Shared:      s.Shared,
	}
}

type WhService interface {
	Create(ctx context.Context, whType int, w *Wh, c *Claims) (*Wh, *WhError)
	Get(ctx context.Context, whType int, whId string, c *Claims) (*Wh, *WhError)
	Update(ctx context.Context, whType int, w *Wh, c *Claims) (*Wh, *WhError)
	Delete(ctx context.Context, whType int, whId string, c *Claims) *WhError
}
