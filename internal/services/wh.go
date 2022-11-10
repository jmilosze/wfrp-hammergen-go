package services

import (
	"context"
	"github.com/go-playground/validator/v10"
	"github.com/jmilosze/wfrp-hammergen-go/internal/domain"
)

type WhService struct {
	Validator   *validator.Validate
	WhDbService domain.WhDbService
	WhType      int
}

func NewWhService(v *validator.Validate, db domain.WhDbService, whType int) *WhService {
	return &WhService{Validator: v, WhDbService: db, WhType: whType}
}

func (s *WhService) Create(ctx context.Context, w domain.Warhammer, c *domain.Claims) (domain.Warhammer, *domain.WhError) {
	if err := w.Validate(s.Validator); err != nil {
		return nil, &domain.WhError{WhType: s.WhType, ErrType: domain.WhInvalidArgumentsError, Err: err}
	}

	if !c.Admin {
		w.SetOwnerId("admin")
	} else {
		w.SetOwnerId(c.Id)
	}
	w.SetNewId()

	createdWh, err := s.WhDbService.Create(ctx, w)
	if err != nil {
		return nil, &domain.WhError{WhType: s.WhType, ErrType: domain.UserInternalError, Err: err}
	}

	return createdWh, nil
}

func (s *WhService) Get(ctx context.Context, whId string, c *domain.Claims) (domain.Warhammer, *domain.WhError) {
	users := []string{"admin", c.Id}
	wh, err := s.WhDbService.Retrieve(ctx, whId, users, c.SharedAccounts)

	if err != nil {
		switch err.Type {
		case domain.DbNotFoundError:
			return nil, &domain.WhError{ErrType: domain.WhNotFoundError, WhType: s.WhType, Err: err}
		default:
			return nil, &domain.WhError{ErrType: domain.WhInternalError, WhType: s.WhType, Err: err}
		}
	}

	wh.SetCanEdit(c.Admin, c.Id, c.SharedAccounts)
	return wh, nil

}
