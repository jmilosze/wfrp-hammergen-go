package mongodb

import (
	"context"
	"github.com/jmilosze/wfrp-hammergen-go/internal/domain"
	"github.com/rs/xid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserDbService struct {
	Db         *DbService
	Collection *mongo.Collection
}

func NewUserDbService(db *DbService, userCollection string) *UserDbService {
	coll := db.Client.Database(db.DbName).Collection(userCollection)
	return &UserDbService{Db: db, Collection: coll}
}

func (s *UserDbService) NewUserDb() *domain.UserDb {
	newId := xid.New().String()
	admin := false
	username := ""
	return &domain.UserDb{
		Id:               newId,
		Username:         &username,
		PasswordHash:     []byte{},
		Admin:            &admin,
		SharedAccountIds: []string{},
	}
}

func (s *UserDbService) Retrieve(ctx context.Context, fieldName string, fieldValue string) (*domain.UserDb, *domain.DbError) {
	return s.NewUserDb(), nil
}

func (s *UserDbService) RetrieveMany(ctx context.Context, fieldName string, fieldValues []string) ([]*domain.UserDb, *domain.DbError) {
	a := make([]*domain.UserDb, 1)
	a[0] = s.NewUserDb()
	return a, nil
}

func (s *UserDbService) Create(ctx context.Context, user *domain.UserDb) *domain.DbError {
	doc := bson.D{
		{"_id", primitive.ObjectIDFromHex(user.Id)},
		{"username", user.Username},
		{"passwordHash", user.PasswordHash},
		{"admin", user.Admin},
	}
	_, err := s.Collection.InsertOne(ctx, doc)
	if err != nil {
		panic(err)
	}

	return nil
}

func (s *UserDbService) Update(ctx context.Context, user *domain.UserDb) (*domain.UserDb, *domain.DbError) {
	return s.NewUserDb(), nil
}

func (s *UserDbService) Delete(ctx context.Context, id string) *domain.DbError {
	return nil
}

func (s *UserDbService) List(ctx context.Context) ([]*domain.UserDb, *domain.DbError) {
	a := make([]*domain.UserDb, 1)
	a[0] = s.NewUserDb()
	return a, nil
}
