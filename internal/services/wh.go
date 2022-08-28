package services

import (
	"context"
	"github.com/go-playground/validator/v10"
	"github.com/jmilosze/wfrp-hammergen-go/internal/domain"
)

type WhService[W domain.WhTypePointer] struct {
	Validator *validator.Validate
}

func NewWhService[W domain.WhTypePointer](v *validator.Validate) *WhService[W] {
	return &WhService[W]{
		Validator: v,
	}
}

func (s *WhService[W]) Create(ctx context.Context, whWrite W) (W, *domain.WhError) {

	//if err := s.Validator.Struct(cred); err != nil {
	//	return nil, &domain.UserError{Type: domain.UserInvalidArgumentsError, Err: err}
	//}

	var newWh W
	return newWh, nil
}
