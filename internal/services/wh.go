package services

import (
	"context"
	"github.com/go-playground/validator/v10"
	"github.com/jmilosze/wfrp-hammergen-go/internal/domain"
)

type WhService[W domain.WhType] struct {
	Validator *validator.Validate
}

func NewWhService[W domain.WhType](v *validator.Validate) *WhService[W] {
	return &WhService[W]{
		Validator: v,
	}
}

func (s *WhService[W]) Create(ctx context.Context, whWrite *W, c *domain.Claims) (*W, *domain.WhError) {

	if err := s.Validator.Struct(whWrite); err != nil {
		return nil, &domain.WhError{ErrType: domain.WhInvalidArgumentsError, Err: err}
	}

	var newWh W
	return &newWh, nil
}
