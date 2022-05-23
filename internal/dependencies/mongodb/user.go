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
	Username       string               `bson:"username"`
	PasswordHash   []byte               `bson:"passwordHash"`
	Admin          *bool                `bson:"admin"`
	SharedAccounts []primitive.ObjectID `bson:"sharedAccounts"`
	CreatedOn      time.Time            `bson:"createdOn"`
	LastAuthOn     time.Time            `bson:"lastAuthOn"`
}

func fromUserDb(u *domain.UserDb, linkedUsers []*UserMongoDb) (*UserMongoDb, error) {
	id, err := primitive.ObjectIDFromHex(u.Id)
	if err != nil {
		return nil, err
	}

	userMongoDb := UserMongoDb{
		Id:             id,
		Username:       u.Username,
		PasswordHash:   u.PasswordHash,
		Admin:          u.Admin,
		SharedAccounts: usernamesToIds(u.SharedAccounts, linkedUsers),
		CreatedOn:      u.CreatedOn,
		LastAuthOn:     u.LastAuthOn,
	}

	return &userMongoDb, nil
}

func usernamesToIds(usernames []string, users []*UserMongoDb) []primitive.ObjectID {
	userMap := map[string]primitive.ObjectID{}
	for _, u := range users {
		userMap[u.Username] = u.Id
	}

	ids := make([]primitive.ObjectID, 0)
	for _, u := range usernames {
		if id, ok := userMap[u]; ok {
			ids = append(ids, id)
		}
	}
	return ids
}

func toUserDb(u *UserMongoDb, linkedUsers []*UserMongoDb) *domain.UserDb {
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
		SharedAccounts: idsToUsernames(u.SharedAccounts, linkedUsers),
		CreatedOn:      u.CreatedOn,
		LastAuthOn:     u.LastAuthOn,
	}
	return &userMongoDb
}

func idsToUsernames(ids []primitive.ObjectID, users []*UserMongoDb) []string {
	userMap := map[primitive.ObjectID]string{}
	for _, u := range users {
		userMap[u.Id] = u.Username
	}

	usernames := make([]string, 0)
	for _, id := range ids {
		if username, ok := userMap[id]; ok {
			usernames = append(usernames, username)
		}
	}
	return usernames
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
	admin := false
	return &domain.UserDb{
		Id:             primitive.NewObjectID().String(),
		Username:       "",
		PasswordHash:   []byte{},
		Admin:          &admin,
		SharedAccounts: []string{},
		CreatedOn:      time.Now(),
		LastAuthOn:     time.Time{},
	}
}

func (s *UserDbService) Retrieve(ctx context.Context, fieldName string, fieldValue string) (*domain.UserDb, *domain.DbError) {
	if fieldName != "username" && fieldName != "id" {
		return nil, &domain.DbError{Type: domain.DbInvalidUserFieldError, Err: fmt.Errorf("invalid field name %s", fieldName)}
	}

	return newUserDb(), nil
}

func getOne(ctx context.Context, coll *mongo.Collection, fieldName string, fieldValue string) (*UserMongoDb, *domain.DbError) {
	var user UserMongoDb
	var err1 error
	if fieldName == "username" {
		err1 = coll.FindOne(ctx, bson.D{{"username", fieldValue}}).Decode(&user)
	} else {
		id, err2 := primitive.ObjectIDFromHex(fieldValue)
		if err2 != nil {
			return nil, &domain.DbError{Type: domain.DbInternalError, Err: err2}
		}
		err1 = coll.FindOne(ctx, bson.D{{"_id", id}}).Decode(&user)
	}
	if err1 != nil {
		return nil, &domain.DbError{Type: domain.DbNotFoundError, Err: errors.New("user not found")}
	}

	return &user, nil
}

func getMany(ctx context.Context, coll *mongo.Collection, fieldName string, fieldValues []string) ([]*UserMongoDb, *domain.DbError) {
	getAll := false
	if fieldValues == nil {
		getAll = true
	} else {
		if len(fieldValues) == 0 {
			return []*UserMongoDb{}, nil
		}
	}

	var query bson.D
	if !getAll {
		var queryValueList []bson.D
		for _, fieldValue := range fieldValues {
			if fieldName == "username" {
				queryValueList = append(queryValueList, bson.D{{"username", fieldValue}})
			} else {
				id, err := primitive.ObjectIDFromHex(fieldValue)
				if err != nil {
					return nil, &domain.DbError{Type: domain.DbInternalError, Err: err}
				}
				queryValueList = append(queryValueList, bson.D{{"_id", id}})
			}
		}
		query = bson.D{{"$or", queryValueList}}
	} else {
		query = bson.D{{}}
	}

	cur, err := coll.Find(ctx, query)
	if err != nil {
		return nil, &domain.DbError{Type: domain.DbInternalError, Err: err}
	}

	users := make([]*UserMongoDb, 0)
	for cur.Next(ctx) {
		var user UserMongoDb
		if cur.Decode(&user) != nil {
			return nil, &domain.DbError{Type: domain.DbInternalError, Err: err}
		}
		users = append(users, &user)
	}

	if cur.Close(ctx) != nil {
		return nil, &domain.DbError{Type: domain.DbInternalError, Err: err}
	}

	return users, nil
}

func (s *UserDbService) RetrieveAll(ctx context.Context) ([]*domain.UserDb, *domain.DbError) {
	a := make([]*domain.UserDb, 1)
	a[0] = newUserDb()
	return a, nil
}

func (s *UserDbService) Create(ctx context.Context, user *domain.UserDb) (*domain.UserDb, *domain.DbError) {
	linkedUsers := make([]*UserMongoDb, 0)
	if user.SharedAccounts != nil {
		var err1 *domain.DbError
		linkedUsers, err1 = getMany(ctx, s.Collection, "username", user.SharedAccounts)
		if err1 != nil {
			return nil, err1
		}
	}

	userMongoDb, err2 := fromUserDb(user, linkedUsers)
	if err2 != nil {
		return nil, &domain.DbError{Type: domain.DbInternalError, Err: err2}
	}

	filter := bson.D{{"_id", userMongoDb.Id}}
	opts := options.Replace().SetUpsert(true)
	if _, err := s.Collection.ReplaceOne(ctx, filter, userMongoDb, opts); err != nil {
		return nil, &domain.DbError{Type: domain.DbInternalError, Err: err}
	}

	return toUserDb(userMongoDb, linkedUsers), nil
}

func (s *UserDbService) Update(ctx context.Context, user *domain.UserDb) (*domain.UserDb, *domain.DbError) {
	return newUserDb(), nil
}

func (s *UserDbService) Delete(ctx context.Context, id string) *domain.DbError {
	return nil
}
