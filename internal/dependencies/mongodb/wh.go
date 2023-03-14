package mongodb

import (
	"context"
	"errors"
	"github.com/jmilosze/wfrp-hammergen-go/internal/domain"
	"go.mongodb.org/mongo-driver/mongo"
)

type WhDbService struct {
	Db          *DbService
	Collections map[domain.WhType]*mongo.Collection
}

func NewWhDbService(db *DbService) *WhDbService {
	collections := map[domain.WhType]*mongo.Collection{
		domain.WhTypeMutation: db.Client.Database(db.DbName).Collection(domain.WhTypeMutation),
		domain.WhTypeSpell:    db.Client.Database(db.DbName).Collection(domain.WhTypeSpell),
	}

	return &WhDbService{Db: db, Collections: collections}
}

func (s *WhDbService) Retrieve(ctx context.Context, t domain.WhType, whId string, userIds []string, sharedUserIds []string) (*domain.Wh, *domain.DbError) {
	return nil, &domain.DbError{
		Type: domain.DbNotImplemented,
		Err:  errors.New("not implemented"),
	}
}

func (s *WhDbService) Create(ctx context.Context, t domain.WhType, w *domain.Wh) (*domain.Wh, *domain.DbError) {
	return nil, &domain.DbError{
		Type: domain.DbNotImplemented,
		Err:  errors.New("not implemented"),
	}
}

func (s *WhDbService) Update(ctx context.Context, t domain.WhType, w *domain.Wh, userId string) (*domain.Wh, *domain.DbError) {
	return nil, &domain.DbError{
		Type: domain.DbNotImplemented,
		Err:  errors.New("not implemented"),
	}
}

func (s *WhDbService) Delete(ctx context.Context, t domain.WhType, whId string, userId string) *domain.DbError {
	return &domain.DbError{
		Type: domain.DbNotImplemented,
		Err:  errors.New("not implemented"),
	}
}

func (s *WhDbService) RetrieveAll(ctx context.Context, t domain.WhType, users []string, sharedUsers []string) ([]*domain.Wh, *domain.DbError) {
	return nil, &domain.DbError{
		Type: domain.DbNotImplemented,
		Err:  errors.New("not implemented"),
	}
}
