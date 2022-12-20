package services

import (
	"context"
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/jmilosze/wfrp-hammergen-go/internal/domain"
	"github.com/rs/xid"
	"golang.org/x/exp/slices"
)

type WhService struct {
	Validator   *validator.Validate
	WhDbService domain.WhDbService
}

func NewWhService(v *validator.Validate, db domain.WhDbService) *WhService {
	return &WhService{Validator: v, WhDbService: db}
}

func (s *WhService) Create(ctx context.Context, whType int, w *domain.Wh, c *domain.Claims) (*domain.Wh, *domain.WhError) {
	if c.Id == "anonymous" {
		return nil, &domain.WhError{WhType: whType, ErrType: domain.WhUnauthorizedError, Err: errors.New("unauthorized")}
	}

	if err := s.Validator.Struct(w); err != nil {
		return nil, &domain.WhError{WhType: whType, ErrType: domain.WhInvalidArgumentsError, Err: err}
	}

	if c.Admin {
		w.OwnerId = "admin"
	} else {
		w.OwnerId = c.Id
	}
	w.Id = xid.New().String()

	createdWh, err := s.WhDbService.Create(ctx, whType, w)
	if err != nil {
		return nil, &domain.WhError{WhType: whType, ErrType: domain.UserInternalError, Err: err}
	}

	createdWh.CanEdit = canEdit(createdWh.OwnerId, c.Admin, c.Id, c.SharedAccounts)
	return createdWh, nil
}

func (s *WhService) Get(ctx context.Context, whType int, whId string, c *domain.Claims) (*domain.Wh, *domain.WhError) {
	users := []string{"admin", c.Id}
	wh, err := s.WhDbService.Retrieve(ctx, whType, whId, users, c.SharedAccounts)

	if err != nil {
		switch err.Type {
		case domain.DbNotFoundError:
			return nil, &domain.WhError{ErrType: domain.WhNotFoundError, WhType: whType, Err: err}
		default:
			return nil, &domain.WhError{ErrType: domain.WhInternalError, WhType: whType, Err: err}
		}
	}

	wh.CanEdit = canEdit(wh.OwnerId, c.Admin, c.Id, c.SharedAccounts)
	return wh, nil

}

func canEdit(ownerId string, isAdmin bool, userId string, sharedAccounts []string) bool {
	if (ownerId != userId) && slices.Contains(sharedAccounts, ownerId) {
		return false
	}

	if isAdmin {
		return true
	}

	if ownerId == userId {
		return true
	}

	return false
}
