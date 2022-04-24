package memdb

import (
	"errors"
	"github.com/hashicorp/go-memdb"
	"github.com/jmilosze/wfrp-hammergen-go/internal/config"
	"github.com/jmilosze/wfrp-hammergen-go/internal/domain"
	"github.com/rs/xid"
	"golang.org/x/crypto/bcrypt"
	"strconv"
	"strings"
)

type UserService struct {
	Db         *memdb.MemDB
	BcryptCost int
}

func NewUserService(cfg *config.MockdbUserService, users []*domain.User) *UserService {
	db, err := createNewMemdb()
	if err != nil {
		panic(err)
	}

	txn := db.Txn(true)
	for i, u := range users {
		id := strconv.Itoa(i)
		pwd, _ := bcrypt.GenerateFromPassword([]byte(u.Password), cfg.BcryptCost)
		userDb := &domain.UserDb{Id: id, Username: u.Username, PasswordHash: pwd, SharedAccounts: u.SharedAccounts, Admin: u.Admin}
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

func (s *UserService) Authenticate(user domain.UserDb, password string) bool {
	if bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(password)) == nil {
		return true
	}
	return false
}

func (s *UserService) Create(newUser *domain.User) (*domain.UserDb, *domain.UserError) {
	if _, err := getUserBy("username", newUser.Username, s); err == nil {
		return nil, &domain.UserError{Type: domain.UserAlreadyExistsError, Err: errors.New("user already exists")}
	}

	newUser.Admin = false

	newId := xid.New().String()
	var userDb = &domain.UserDb{Id: newId}
	_ = updateDbUser(userDb, newUser, s.BcryptCost)

	txn := s.Db.Txn(true)
	defer txn.Abort()
	if err := txn.Insert("user", userDb); err != nil {
		return nil, &domain.UserError{Type: domain.UserInternalError, Err: err}
	}
	txn.Commit()

	return userDb.Copy(), nil
}

func (s *UserService) SimpleUpdate(id string, newUser *domain.User) (*domain.UserDb, *domain.UserError) {

	userDb, err := s.GetById(id)
	if err != nil {
		return nil, err
	}

	newUser.Password = ""
	newUser.Username = strings.Clone(userDb.Username)
	newUser.Admin = userDb.Admin

	_ = updateDbUser(userDb, newUser, s.BcryptCost)

	txn := s.Db.Txn(true)
	defer txn.Abort()
	if err := txn.Insert("user", userDb); err != nil {
		return nil, &domain.UserError{Type: domain.UserInternalError, Err: err}
	}
	txn.Commit()

	return userDb.Copy(), nil
}

func updateDbUser(userDb *domain.UserDb, user *domain.User, bcryptCost int) *domain.UserError {
	userDb.Username = strings.Clone(user.Username)

	userDb.SharedAccounts = make([]string, len(user.SharedAccounts))
	for i, s := range user.SharedAccounts {
		userDb.SharedAccounts[i] = strings.Clone(s)
	}

	if user.Password != "" {
		userDb.PasswordHash, _ = bcrypt.GenerateFromPassword([]byte(user.Password), bcryptCost)
	}

	return nil
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
