package mongodb

import (
	"context"
	"github.com/jmilosze/wfrp-hammergen-go/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

type UserMongoDb struct {
	Id               primitive.ObjectID `bson:"_id"`
	Username         *string            `bson:"username"`
	PasswordHash     []byte             `bson:"passwordHash"`
	Admin            *bool              `bson:"admin"`
	SharedAccountIds []string           `bson:"sharedAccountIds"`
	CreatedOn        time.Time          `bson:"createdOn"`
	LastAuthOn       time.Time          `bson:"lastAuthOn"`
}

func fromUserDb(u *domain.UserDb) (*UserMongoDb, error) {
	id, err := primitive.ObjectIDFromHex(u.Id)
	if err != nil {
		return nil, err
	}

	userMongoDb := UserMongoDb{
		Id:               id,
		Username:         u.Username,
		PasswordHash:     u.PasswordHash,
		Admin:            u.Admin,
		SharedAccountIds: u.SharedAccountIds,
		CreatedOn:        u.CreatedOn,
		LastAuthOn:       u.LastAuthOn,
	}

	return &userMongoDb, nil
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

func (s *UserDbService) NewUserDb() *domain.UserDb {
	newId := primitive.NewObjectID().String()
	admin := false
	username := ""
	return &domain.UserDb{
		Id:               newId,
		Username:         &username,
		PasswordHash:     []byte{},
		Admin:            &admin,
		SharedAccountIds: []string{},
		CreatedOn:        time.Now(),
		LastAuthOn:       time.Time{},
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
