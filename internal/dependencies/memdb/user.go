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

type UserDb struct {
	Id             string
	Username       string
	PasswordHash   []byte
	Admin          bool
	SharedAccounts []string
}

func (u *UserDb) Copy() *UserDb {
	userCopy := *u
	userCopy.Username = strings.Clone(u.Username)
	userCopy.PasswordHash = make([]byte, len(u.PasswordHash))
	copy(userCopy.PasswordHash, u.PasswordHash)
	userCopy.SharedAccounts = make([]string, len(u.SharedAccounts))
	for i, s := range u.SharedAccounts {
		userCopy.SharedAccounts[i] = strings.Clone(s)
	}

	return &userCopy
}

func NewUserService(cfg *config.MockDbUserService, users map[string]*config.UserSeed) *UserService {
	db, err := createNewMemdb()
	if err != nil {
		panic(err)
	}

	txn := db.Txn(true)
	for id, u := range users {
		var userDb = UserDb{Id: id}
		_ = updateUserDbCredentials(&userDb, u.Credentials, cfg.BcryptCost)
		_ = updateUserDb(&userDb, u.User)
		_ = updateUserDbClaims(&userDb, u.Claims)
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

func (s *UserService) GetById(id string) (*UserDb, *domain.UserError) {
	return getUserBy("id", id, s)
}

func (s *UserService) GetByName(username string) (*UserDb, *domain.UserError) {
	return getUserBy("username", username, s)
}

func getUserBy(fieldName string, fieldValue string, s *UserService) (*UserDb, *domain.UserError) {
	txn := s.Db.Txn(false)
	userRaw, err := txn.First("user", fieldName, fieldValue)
	if err != nil {
		return nil, &domain.UserError{Type: domain.UserInternalError, Err: err}
	}

	if userRaw == nil {
		return nil, &domain.UserError{Type: domain.UserNotFoundError, Err: errors.New("user not found")}
	}
	user := userRaw.(*UserDb)

	return user.Copy(), nil
}

func (s *UserService) Authenticate(user *UserDb, password string) bool {
	if bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(password)) == nil {
		return true
	}
	return false
}

func (s *UserService) Create(cred *domain.UserWriteCredentials, user *domain.UserWrite) (*UserDb, *domain.UserError) {
	if _, err := getUserBy("username", cred.Username, s); err == nil {
		return nil, &domain.UserError{Type: domain.UserAlreadyExistsError, Err: errors.New("user already exists")}
	}

	newId := xid.New().String()
	userDb := UserDb{Id: newId}

	_ = updateUserDbCredentials(&userDb, cred, s.BcryptCost)
	_ = updateUserDb(&userDb, user)
	return insertUser(s, &userDb)
}

func updateUserDbCredentials(userDb *UserDb, cred *domain.UserWriteCredentials, bcryptCost int) *domain.UserError {
	userDb.Username = strings.Clone(cred.Username)
	userDb.PasswordHash, _ = bcrypt.GenerateFromPassword([]byte(cred.Password), bcryptCost)
	return nil
}

func updateUserDb(userDb *UserDb, user *domain.UserWrite) *domain.UserError {
	userDb.SharedAccounts = make([]string, len(user.SharedAccounts))
	for i, s := range user.SharedAccounts {
		userDb.SharedAccounts[i] = strings.Clone(s)
	}
	return nil
}

func updateUserDbClaims(userDb *UserDb, claims *domain.UserWriteClaims) *domain.UserError {
	userDb.Admin = claims.Admin
	return nil
}

func insertUser(s *UserService, u *UserDb) (*UserDb, *domain.UserError) {
	txn := s.Db.Txn(true)
	defer txn.Abort()
	if err := txn.Insert("user", u); err != nil {
		return nil, &domain.UserError{Type: domain.UserInternalError, Err: err}
	}
	txn.Commit()
	return u.Copy(), nil
}

func (s *UserService) Update(id string, user *domain.UserWrite) (*UserDb, *domain.UserError) {
	userDb, err := s.GetById(id)
	if err != nil {
		return nil, err
	}

	_ = updateUserDb(userDb, user)
	return insertUser(s, userDb)
}

func (s *UserService) UpdateCredentials(id string, passwd string, cred *domain.UserWriteCredentials) (*UserDb, *domain.UserError) {
	userDb, err := s.GetById(id)
	if err != nil {
		return nil, err
	}

	if !s.Authenticate(userDb, passwd) {
		return nil, &domain.UserError{Type: domain.UserIncorrectPassword, Err: errors.New("incorrect password")}
	}

	_ = updateUserDbCredentials(userDb, cred, s.BcryptCost)
	return insertUser(s, userDb)
}

func (s *UserService) UpdateClaims(id string, claims *domain.UserWriteClaims) (*UserDb, *domain.UserError) {
	userDb, err := s.GetById(id)
	if err != nil {
		return nil, err
	}

	_ = updateUserDbClaims(userDb, claims)

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

func (s *UserService) List() ([]*UserDb, *domain.UserError) {
	txn := s.Db.Txn(false)

	it, err := txn.Get("user", "id")
	if err != nil {
		return nil, &domain.UserError{Type: domain.UserInternalError, Err: err}
	}

	var users []*UserDb
	for obj := it.Next(); obj != nil; obj = it.Next() {
		u := obj.(*UserDb)
		users = append(users, u.Copy())
	}

	return users, nil
}
