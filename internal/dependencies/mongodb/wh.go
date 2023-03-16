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
	id, err := primitive.ObjectIDFromHex(whId)
	if err != nil {
		return nil, &domain.DbError{
			Type: domain.DbInternalError,
			Err:  err,
		}
	}

	filter := bson.M{"$and": bson.A{bson.M{"_id": id}, allAllowedOwnersQuery(userIds, sharedUserIds)}}
	var whMap bson.M

	err = s.Collections[t].FindOne(ctx, filter).Decode(&whMap)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, &domain.DbError{
				Type: domain.DbNotFoundError,
				Err:  err,
			}
		} else {
			return nil, &domain.DbError{
				Type: domain.DbInternalError,
				Err:  err,
			}
		}
	}

	wh, err := bsonMToWh(whMap, t)
	if err != nil {
		return nil, &domain.DbError{
			Type: domain.DbInternalError,
			Err:  err,
		}
	}

	return wh, nil
}

func allAllowedOwnersQuery(userIds []string, sharedUserIds []string) bson.M {
	owners := bson.A{}
	for _, v := range userIds {
		owners = append(owners, bson.M{"ownerid": v})
	}

	if sharedUserIds != nil && len(sharedUserIds) > 0 {
		sharedOwners := bson.A{}
		for _, v := range sharedUserIds {
			sharedOwners = append(sharedOwners, bson.M{"ownerid": v})
		}
		owners = append(owners, bson.M{"$and": bson.A{bson.M{"shared": true}, bson.M{"$or": sharedOwners}}})
	}
	return bson.M{"$or": owners}
}

func bsonMToWh(whMap bson.M, t domain.WhType) (*domain.Wh, error) {
	id, ok := whMap["_id"].(primitive.ObjectID)
	if !ok {
		return nil, errors.New("invalid object id")
	}

	ownerId, ok := whMap["ownerid"].(string)
	if !ok {
		return nil, errors.New("invalid owner id")
	}

	wh := domain.Wh{
		Id:      id.Hex(),
		OwnerId: ownerId,
		CanEdit: false,
	}

	switch t {
	case domain.WhTypeMutation:
		wh.Object = &domain.WhMutation{}
	case domain.WhTypeSpell:
		wh.Object = &domain.WhSpell{}
	default:
		errors.New("unknown wh type")
	}

	bsonRaw, err := bson.Marshal(whMap["object"])
	if err != nil {
		errors.New("error marshaling object")
	}

	if err = bson.Unmarshal(bsonRaw, wh.Object); err != nil {
		errors.New("error unmarshalling object")
	}

	return &wh, nil
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
		if mongo.IsDuplicateKeyError(err) {
			return nil, &domain.DbError{
				Type: domain.DbAlreadyExistsError,
				Err:  err,
			}
		}
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
