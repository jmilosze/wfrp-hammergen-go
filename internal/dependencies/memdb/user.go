package memdb

import (
	"errors"
	"github.com/hashicorp/go-memdb"
	"github.com/jmilosze/wfrp-hammergen-go/internal/config"
	"github.com/jmilosze/wfrp-hammergen-go/internal/domain"
	"github.com/rs/xid"
	"golang.org/x/crypto/bcrypt"
	"strings"
)

type UserService struct {
	Db         *memdb.MemDB
	BcryptCost int
}

func NewUserService(cfg *config.MockDbUserService, users map[string]*domain.User) *UserService {
	db, err := createNewMemdb()
	if err != nil {
		panic(err)
	}

	txn := db.Txn(true)
	for id, u := range users {
		var userDb = &domain.UserDb{Id: id}
		_ = updateDbUser(userDb, u, cfg.BcryptCost, true, true)
		if err := txn.Insert("user", userDb); err != nil {
			panic(err)
		}
	}

	txn.Commit()

	return &UserService{Db: db, BcryptCost: cfg.BcryptCost}
}

func createNewMemdb() (*memdb.MemDB, error) {
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

func (s *UserService) GetById(id string) (*domain.UserDb, *domain.UserError) {
	return getUserBy("id", id, s)
}

func (s *UserService) GetByName(username string) (*domain.UserDb, *domain.UserError) {
	return getUserBy("username", username, s)
}

func getUserBy(fieldName string, fieldValue string, s *UserService) (*domain.UserDb, *domain.UserError) {
	txn := s.Db.Txn(false)
	userRaw, err := txn.First("user", fieldName, fieldValue)
	if err != nil {
		return nil, &domain.UserError{Type: domain.UserInternalError, Err: err}
	}

	if userRaw == nil {
		return nil, &domain.UserError{Type: domain.UserNotFoundError, Err: errors.New("user not found")}
	}
	user := userRaw.(*domain.UserDb)

	return user.Copy(), nil
}

func (s *UserService) Authenticate(user *domain.UserDb, password string) bool {
	if bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(password)) == nil {
		return true
	}
	return false
}

func (s *UserService) Create(newUser *domain.User) (*domain.UserDb, *domain.UserError) {
	if _, err := getUserBy("username", newUser.Username, s); err == nil {
		return nil, &domain.UserError{Type: domain.UserAlreadyExistsError, Err: errors.New("user already exists")}
	}

	newId := xid.New().String()
	var userDb = &domain.UserDb{Id: newId, Admin: false}

	_ = updateDbUser(userDb, newUser, s.BcryptCost, true, false)
	return insertUser(s, userDb)
}

func insertUser(s *UserService, u *domain.UserDb) (*domain.UserDb, *domain.UserError) {
	txn := s.Db.Txn(true)
	defer txn.Abort()
	if err := txn.Insert("user", u); err != nil {
		return nil, &domain.UserError{Type: domain.UserInternalError, Err: err}
	}
	txn.Commit()
	return u.Copy(), nil
}

func (s *UserService) Update(id string, newUser *domain.User) (*domain.UserDb, *domain.UserError) {
	userDb, err := s.GetById(id)
	if err != nil {
		return nil, err
	}

	_ = updateDbUser(userDb, newUser, s.BcryptCost, false, false)
	return insertUser(s, userDb)
}

func updateDbUser(userDb *domain.UserDb, user *domain.User, bcryptCost int, updateCredentials bool, updateAdmin bool) *domain.UserError {
	userDb.SharedAccounts = make([]string, len(user.SharedAccounts))
	for i, s := range user.SharedAccounts {
		userDb.SharedAccounts[i] = strings.Clone(s)
	}

	if updateCredentials {
		userDb.Username = strings.Clone(user.Username)
		userDb.PasswordHash, _ = bcrypt.GenerateFromPassword([]byte(user.Password), bcryptCost)
	}

	if updateAdmin {
		userDb.Admin = user.Admin
	}

	return nil
}

func (s *UserService) UpdateCredentials(id string, passwd string, newPasswd string, newUsername string) (*domain.UserDb, *domain.UserError) {
	userDb, err := s.GetById(id)
	if err != nil {
		return nil, err
	}

	if !s.Authenticate(userDb, passwd) {
		return nil, &domain.UserError{Type: domain.UserIncorrectPassword, Err: errors.New("incorrect password")}
	}

	if len(newUsername) != 0 {
		userDb.Username = strings.Clone(newUsername)
	}

	if len(newPasswd) != 0 {
		userDb.PasswordHash, _ = bcrypt.GenerateFromPassword([]byte(newUsername), s.BcryptCost)
	}

	return insertUser(s, userDb)
}

func (s *UserService) Delete(id string) *domain.UserError {
	txn := s.Db.Txn(true)
	defer txn.Abort()
	if _, err := txn.DeleteAll("user", "id", id); err != nil {
		return &domain.UserError{Type: domain.UserInternalError, Err: err}
	}
	txn.Commit()

	return nil
}

func (s *UserService) List() ([]*domain.UserDb, *domain.UserError) {
	txn := s.Db.Txn(false)

	it, err := txn.Get("user", "id")
	if err != nil {
		return nil, &domain.UserError{Type: domain.UserInternalError, Err: err}
	}

	var users []*domain.UserDb
	for obj := it.Next(); obj != nil; obj = it.Next() {
		u := obj.(*domain.UserDb)
		users = append(users, u.Copy())
	}

	return users, nil
}
