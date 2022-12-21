package memdb

import (
	"context"
	"errors"
	"fmt"
	"github.com/hashicorp/go-memdb"
	"github.com/jmilosze/wfrp-hammergen-go/internal/domain"
	"strings"
)

type UserDbService struct {
	Db *memdb.MemDB
}

func NewUserDbService() *UserDbService {
	db, err := createNewUserMemDb()
	if err != nil {
		panic(err)
	}

	return &UserDbService{Db: db}
}

func createNewUserMemDb() (*memdb.MemDB, error) {
	schema := &memdb.DBSchema{
		Tables: map[string]*memdb.TableSchema{
			"user": {
				Name: "user",
				Indexes: map[string]*memdb.IndexSchema{
					"id": {
						Name:    "id",
						Unique:  true,
						Indexer: &memdb.StringFieldIndex{Field: "Id"},
					},
					"username": {
						Name:    "username",
						Unique:  true,
						Indexer: &memdb.StringFieldIndex{Field: "Username"},
					},
				},
			},
		},
	}
	return memdb.NewMemDB(schema)
}

func (s *UserDbService) Retrieve(ctx context.Context, fieldName string, fieldValue string) (*domain.UserDb, *domain.DbError) {
	if fieldName != "username" && fieldName != "id" {
		return nil, &domain.DbError{Type: domain.DbInvalidUserFieldError, Err: fmt.Errorf("invalid field name %s", fieldName)}
	}

	user, err1 := getOneUser(s.Db, fieldName, fieldValue)
	if err1 != nil {
		return nil, err1
	}

	linkedUsers, err2 := getManyUsers(s.Db, "id", user.SharedAccountIds)
	if err2 != nil {
		return nil, err2
	}

	user.SharedAccountNames = idsToUsernames(user.SharedAccountIds, linkedUsers)

	return user, nil
}

func getOneUser(db *memdb.MemDB, fieldName string, fieldValue string) (*domain.UserDb, *domain.DbError) {
	txn := db.Txn(false)
	userRaw, err := txn.First("user", fieldName, fieldValue)
	if err != nil {
		return nil, &domain.DbError{Type: domain.DbInternalError, Err: err}
	}

	if userRaw == nil {
		return nil, &domain.DbError{Type: domain.DbNotFoundError, Err: errors.New("user not found")}
	}
	user := userRaw.(*domain.UserDb)

	return user.Copy(), nil
}

func getManyUsers(db *memdb.MemDB, fieldName string, fieldValues []string) ([]*domain.UserDb, *domain.DbError) {
	getAll := false
	if fieldValues == nil {
		getAll = true
	} else {
		if len(fieldValues) == 0 {
			return []*domain.UserDb{}, nil
		}
	}

	txn := db.Txn(false)

	it, err := txn.Get("user", "id")
	if err != nil {
		return nil, &domain.DbError{Type: domain.DbInternalError, Err: err}
	}

	var users []*domain.UserDb
	for obj := it.Next(); obj != nil; obj = it.Next() {
		u := obj.(*domain.UserDb)
		if getAll {
			users = append(users, u.Copy())
		} else {
			if fieldName == "username" && contains(fieldValues, u.Username) {
				users = append(users, u.Copy())
			}
			if fieldName == "id" && contains(fieldValues, u.Id) {
				users = append(users, u.Copy())
			}
		}
	}
	return users, nil
}

func contains(elems []string, v string) bool {
	for _, s := range elems {
		if v == s {
			return true
		}
	}
	return false
}

func idsToUsernames(ids []string, userDbs []*domain.UserDb) []string {
	userDbMap := map[string]string{}
	for _, u := range userDbs {
		userDbMap[u.Id] = u.Username
	}

	usernames := make([]string, 0)
	for _, id := range ids {
		if username, ok := userDbMap[id]; ok {
			usernames = append(usernames, username)
		}
	}
	return usernames
}

func (s *UserDbService) RetrieveAll(ctx context.Context) ([]*domain.UserDb, *domain.DbError) {
	users, err := getManyUsers(s.Db, "username", nil)
	if err != nil {
		return nil, err
	}

	for _, u := range users {
		u.SharedAccountNames = idsToUsernames(u.SharedAccountIds, users)
	}

	return users, nil
}

func (s *UserDbService) Create(ctx context.Context, user *domain.UserDb) (*domain.UserDb, *domain.DbError) {
	_, err1 := getOneUser(s.Db, "username", user.Username)
	if err1 == nil {
		return nil, &domain.DbError{Type: domain.DbAlreadyExistsError, Err: errors.New("user already exists")}
	}
	if err1 != nil && err1.Type != domain.DbNotFoundError {
		return nil, &domain.DbError{Type: domain.DbInternalError, Err: err1.Unwrap()}
	}

	userDbCreate := user.Copy()

	if user.SharedAccountNames != nil {
		linkedUsers := make([]*domain.UserDb, 0)
		var err2 *domain.DbError
		linkedUsers, err2 = getManyUsers(s.Db, "username", userDbCreate.SharedAccountNames)
		if err2 != nil {
			return nil, err2
		}
		userDbCreate.SharedAccountIds = usernamesToIds(userDbCreate.SharedAccountNames, linkedUsers)
	}

	txn := s.Db.Txn(true)
	defer txn.Abort()
	if err := txn.Insert("user", userDbCreate); err != nil {
		return nil, &domain.DbError{Type: domain.DbInternalError, Err: err}
	}
	txn.Commit()

	return userDbCreate.Copy(), nil
}

func usernamesToIds(usernames []string, userDbs []*domain.UserDb) []string {
	userDbMap := map[string]string{}
	for _, u := range userDbs {
		userDbMap[u.Username] = u.Id
	}

	ids := make([]string, 0)
	for _, u := range usernames {
		if id, ok := userDbMap[u]; ok {
			ids = append(ids, id)
		}
	}
	return ids
}

func (s *UserDbService) Update(ctx context.Context, user *domain.UserDb) (*domain.UserDb, *domain.DbError) {
	userDb, err1 := getOneUser(s.Db, "id", user.Id)
	if err1 != nil {
		return nil, err1
	}

	updateUserDb(userDb, user)

	if user.SharedAccountNames != nil {
		userDb.SharedAccountNames = make([]string, len(user.SharedAccountNames))
		for i, v := range user.SharedAccountNames {
			userDb.SharedAccountNames[i] = v
		}

		linkedUsers, err3 := getManyUsers(s.Db, "username", userDb.SharedAccountNames)
		if err3 != nil {
			return nil, err3
		}

		userDb.SharedAccountIds = usernamesToIds(userDb.SharedAccountNames, linkedUsers)
	}

	txn := s.Db.Txn(true)
	defer txn.Abort()
	if err4 := txn.Insert("user", userDb); err4 != nil {
		return nil, &domain.DbError{Type: domain.DbInternalError, Err: err4}
	}
	txn.Commit()

	return userDb.Copy(), nil
}

func updateUserDb(to *domain.UserDb, from *domain.UserDb) *domain.UserDb {
	if len(from.Username) != 0 {
		to.Username = strings.Clone(from.Username)
	}

	if from.PasswordHash != nil {
		to.PasswordHash = make([]byte, len(from.PasswordHash))
		for i, s := range from.PasswordHash {
			to.PasswordHash[i] = s
		}
	}

	if from.Admin != nil {
		*to.Admin = *from.Admin
	}

	if !from.LastAuthOn.IsZero() {
		to.LastAuthOn = from.LastAuthOn.UTC()
	}

	if !from.CreatedOn.IsZero() {
		to.CreatedOn = from.CreatedOn.UTC()
	}

	return to
}

func (s *UserDbService) Delete(ctx context.Context, id string) *domain.DbError {
	txn := s.Db.Txn(true)
	defer txn.Abort()
	if _, err := txn.DeleteAll("user", "id", id); err != nil {
		return &domain.DbError{Type: domain.DbInternalError, Err: err}
	}
	txn.Commit()

	return nil
}
