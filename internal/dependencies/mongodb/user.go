package mongodb

import (
	"context"
	"errors"
	"fmt"
	"github.com/jmilosze/wfrp-hammergen-go/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

type UserMongoDb struct {
	Id             primitive.ObjectID   `bson:"_id"`
	Username       *string              `bson:"username"`
	PasswordHash   []byte               `bson:"passwordHash"`
	Admin          *bool                `bson:"admin"`
	SharedAccounts []primitive.ObjectID `bson:"sharedAccounts"`
	CreatedOn      *time.Time           `bson:"createdOn"`
	LastAuthOn     *time.Time           `bson:"lastAuthOn"`
}

func fromUserDb(u *domain.UserDb) (*UserMongoDb, error) {
	id, err := primitive.ObjectIDFromHex(u.Id)
	if err != nil {
		return nil, err
	}

	var sharedAccounts []primitive.ObjectID
	if u.SharedAccounts != nil {
		sharedAccounts = make([]primitive.ObjectID, len(u.SharedAccounts))
		for i, sa := range u.SharedAccounts {
			id, err = primitive.ObjectIDFromHex(sa)
			if err != nil {
				return nil, err
			}
			sharedAccounts[i] = id
		}
	} else {
		sharedAccounts = nil
	}

	userMongoDb := UserMongoDb{
		Id:             id,
		Username:       u.Username,
		PasswordHash:   u.PasswordHash,
		Admin:          u.Admin,
		SharedAccounts: sharedAccounts,
		CreatedOn:      u.CreatedOn,
		LastAuthOn:     u.LastAuthOn,
	}

	return &userMongoDb, nil
}

func toUserDb(u *UserMongoDb) *domain.UserDb {
	var sharedAccounts []string
	if u.SharedAccounts != nil {
		sharedAccounts = make([]string, len(u.SharedAccounts))
		for i, sa := range u.SharedAccounts {
			sharedAccounts[i] = sa.Hex()
		}
	} else {
		sharedAccounts = nil
	}

	userMongoDb := domain.UserDb{
		Id:             u.Id.Hex(),
		Username:       u.Username,
		PasswordHash:   u.PasswordHash,
		Admin:          u.Admin,
		SharedAccounts: sharedAccounts,
		CreatedOn:      u.CreatedOn,
		LastAuthOn:     u.LastAuthOn,
	}
	return &userMongoDb
}

type UserDbService struct {
	Db         *DbService
	Collection *mongo.Collection
}

func NewUserDbService(db *DbService, userCollection string, createIndex bool) *UserDbService {
	coll := db.Client.Database(db.DbName).Collection(userCollection)

	if createIndex {
		unique := true
		mod := mongo.IndexModel{Keys: bson.M{"username": 1}, Options: &options.IndexOptions{Unique: &unique}}
		_, err := coll.Indexes().CreateOne(context.TODO(), mod)
		if err != nil {
			log.Fatal(err)
		}
	}

	return &UserDbService{Db: db, Collection: coll}
}

func newUserDb() *domain.UserDb {
	newId := primitive.NewObjectID().String()
	admin := false
	username := ""
	now := time.Now()
	return &domain.UserDb{
		Id:             newId,
		Username:       &username,
		PasswordHash:   []byte{},
		Admin:          &admin,
		SharedAccounts: []string{},
		CreatedOn:      &now,
		LastAuthOn:     &time.Time{},
	}
}

func (s *UserDbService) Retrieve(ctx context.Context, fieldName string, fieldValue string) (*domain.UserDb, *domain.DbError) {
	if fieldName != "username" && fieldName != "id" {
		return nil, &domain.DbError{Type: domain.DbInvalidUserFieldError, Err: fmt.Errorf("invalid field name %s", fieldName)}
	}

	var userMongoDb UserMongoDb
	var err1 error
	if fieldName == "username" {
		err1 = s.Collection.FindOne(ctx, bson.D{{"username", fieldValue}}).Decode(&userMongoDb)
	} else {
		id, err2 := primitive.ObjectIDFromHex(fieldValue)
		if err2 != nil {
			return nil, &domain.DbError{Type: domain.DbInternalError, Err: err2}
		}
		err1 = s.Collection.FindOne(ctx, bson.D{{"_id", id}}).Decode(&userMongoDb)
	}
	if err1 != nil {
		return nil, &domain.DbError{Type: domain.DbNotFoundError, Err: errors.New("user not found")}
	}

	return toUserDb(&userMongoDb), nil
}

func (s *UserDbService) RetrieveAll(ctx context.Context) ([]*domain.UserDb, *domain.DbError) {
	a := make([]*domain.UserDb, 1)
	a[0] = newUserDb()
	return a, nil
}

func (s *UserDbService) Create(ctx context.Context, user *domain.UserDb) *domain.DbError {
	newUser, err := fromUserDb(user)
	if err != nil {
		return &domain.DbError{Type: domain.DbInternalError, Err: err}
	}

	filter := bson.D{{"_id", newUser.Id}}
	opts := options.Replace().SetUpsert(true)
	if _, err := s.Collection.ReplaceOne(ctx, filter, newUser, opts); err != nil {
		return &domain.DbError{Type: domain.DbInternalError, Err: err}
	}

	return nil
}

func (s *UserDbService) Update(ctx context.Context, user *domain.UserDb) (*domain.UserDb, *domain.DbError) {
	return newUserDb(), nil
}

func (s *UserDbService) Delete(ctx context.Context, id string) *domain.DbError {
	return nil
}
