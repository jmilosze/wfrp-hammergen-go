package domain

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/rs/xid"
	"golang.org/x/exp/slices"
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
	WhNotFoundError
	WhInternalError
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

type Warhammer interface {
	SetOwnerId(ownerId string)
	GetCommonFields() *WhCommonFields
	PopulateFromJson(jsonData []byte) error
	Validate(v *validator.Validate) error
	Copy() Warhammer
	SetNewId()
	SetCanEdit(isAdmin bool, usersId string, sharedAccounts []string)
}

type WhCommonFields struct {
	Id      string
	OwnerId string
	Shared  bool
}

type Mutation struct {
	Name        string `json:"name" validate:"omitempty,min=0,max=200,excludesall=<>"`
	Description string `json:"description" validate:"omitempty,min=0,max=100000,excludesall=<>"`
	Type        int    `json:"type" validate:"omitempty,oneof=0 1"`
	Shared      bool   `json:"shared" validate:"omitempty"`
	Id          string `json:"id"`
	OwnerId     string `json:"owner_id"`
	CanEdit     bool   `json:"can_id"`
}

func (m *Mutation) SetOwnerId(ownerId string) {
	m.OwnerId = ownerId
}

func (m *Mutation) GetCommonFields() *WhCommonFields {
	return &WhCommonFields{
		Id:      m.Id,
		OwnerId: m.OwnerId,
		Shared:  m.Shared,
	}
}

func (m *Mutation) PopulateFromJson(jsonData []byte) error {
	return json.Unmarshal(jsonData, m)
}

func (m *Mutation) Validate(v *validator.Validate) error {
	return v.Struct(m)
}

func (m *Mutation) Copy() Warhammer {
	if m == nil {
		return nil
	}

	return &Mutation{
		Name:        strings.Clone(m.Name),
		Description: strings.Clone(m.Description),
		Type:        m.Type,
		Shared:      m.Shared,
		Id:          strings.Clone(m.Id),
		OwnerId:     strings.Clone(m.OwnerId),
	}
}

func (m *Mutation) SetNewId() {
	m.Id = xid.New().String()
}

func (m *Mutation) SetCanEdit(isAdmin bool, userId string, sharedAccounts []string) {
	if (m.OwnerId != userId) && slices.Contains(sharedAccounts, m.OwnerId) {
		m.CanEdit = false
		return
	}

	if isAdmin {
		m.CanEdit = true
		return
	}

	if m.OwnerId == userId {
		m.CanEdit = true
		return
	}

	m.CanEdit = false
}

type Spell struct {
	Name        string `json:"name" validate:"omitempty,min=0,max=200,excludesall=<>"`
	Description string `json:"description" validate:"omitempty,min=0,max=100000,excludesall=<>"`
	Cn          int    `json:"cn" validate:"omitempty,min=-1,max=99"`
	Range       string `json:"range" validate:"omitempty,min=0,max=200,excludesall=<>"`
	Target      string `json:"target" validate:"omitempty,min=0,max=200,excludesall=<>"`
	Duration    string `json:"duration" validate:"omitempty,min=0,max=200,excludesall=<>"`
	Shared      bool   `json:"shared" validate:"omitempty"`
	Id          string `json:"id"`
	OwnerId     string `json:"owner_id"`
	CanEdit     bool   `json:"can_id"`
}

func (s *Spell) SetOwnerId(ownerId string) {
	s.OwnerId = ownerId
}

func (s *Spell) GetCommonFields() *WhCommonFields {
	return &WhCommonFields{
		Id:      s.Id,
		OwnerId: s.OwnerId,
		Shared:  s.Shared,
	}
}

func (s *Spell) PopulateFromJson(jsonData []byte) error {
	return json.Unmarshal(jsonData, s)
}

func (s *Spell) Validate(v *validator.Validate) error {
	return v.Struct(s)
}

func (s *Spell) Copy() Warhammer {
	if s == nil {
		return nil
	}

	return &Spell{
		Name:        strings.Clone(s.Name),
		Description: strings.Clone(s.Description),
		Cn:          s.Cn,
		Range:       strings.Clone(s.Range),
		Target:      strings.Clone(s.Target),
		Duration:    strings.Clone(s.Duration),
		Shared:      s.Shared,
		Id:          strings.Clone(s.Id),
		OwnerId:     strings.Clone(s.OwnerId),
	}
}

func (s *Spell) SetCanEdit(isAdmin bool, userId string, sharedAccounts []string) {
	if (s.OwnerId != userId) && slices.Contains(sharedAccounts, s.OwnerId) {
		s.CanEdit = false
		return
	}

	if isAdmin {
		s.CanEdit = true
		return
	}

	if s.OwnerId == userId {
		s.CanEdit = true
		return
	}

	s.CanEdit = false
}

func (s *Spell) SetNewId() {
	s.Id = xid.New().String()
}

type WhService interface {
	Create(ctx context.Context, w Warhammer, c *Claims) (Warhammer, *WhError)
	Get(ctx context.Context, whId string, c *Claims) (Warhammer, *WhError)
}
