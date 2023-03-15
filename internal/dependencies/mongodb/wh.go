package mongodb

import (
	"context"
	"errors"
	"github.com/jmilosze/wfrp-hammergen-go/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
		Type: domain.DbNotImplementedError,
		Err:  errors.New("not implemented"),
	}
}

func (s *WhDbService) Create(ctx context.Context, t domain.WhType, w *domain.Wh) (*domain.Wh, *domain.DbError) {
	whBsonM, err := whToBsonM(w)
	if err != nil {
		return nil, &domain.DbError{
			Type: domain.DbWriteToDbError,
			Err:  err,
		}
	}

	_, err = s.Collections[t].InsertOne(ctx, whBsonM)
	if err != nil {
		return nil, &domain.DbError{
			Type: domain.DbWriteToDbError,
			Err:  err,
		}
	}

	return w, nil
}

func whToBsonM(w *domain.Wh) (bson.M, error) {
	wBson, err := bson.Marshal(w)
	if err != nil {
		return nil, err
	}

	var whMap bson.M
	err = bson.Unmarshal(wBson, &whMap)
	if err != nil {
		return nil, err
	}

	id, err := primitive.ObjectIDFromHex(w.Id)
	if err != nil {
		return nil, err
	}

	delete(whMap, "canedit")
	delete(whMap, "id")
	whMap["_id"] = id

	return whMap, err
}

func (s *WhDbService) Update(ctx context.Context, t domain.WhType, w *domain.Wh, userId string) (*domain.Wh, *domain.DbError) {
	return nil, &domain.DbError{
		Type: domain.DbNotImplementedError,
		Err:  errors.New("not implemented"),
	}
}

func (s *WhDbService) Delete(ctx context.Context, t domain.WhType, whId string, userId string) *domain.DbError {
	return &domain.DbError{
		Type: domain.DbNotImplementedError,
		Err:  errors.New("not implemented"),
	}
}

func (s *WhDbService) RetrieveAll(ctx context.Context, t domain.WhType, users []string, sharedUsers []string) ([]*domain.Wh, *domain.DbError) {
	return nil, &domain.DbError{
		Type: domain.DbNotImplementedError,
		Err:  errors.New("not implemented"),
	}
}
