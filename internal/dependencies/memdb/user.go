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

func NewUserService(cfg *config.MockdbUserService) *UserService {
	db, err := createNewMemdb()
	if err != nil {
		panic(err)
	}

	user1password, _ := bcrypt.GenerateFromPassword([]byte("123"), cfg.BcryptCost)
	user2password, _ := bcrypt.GenerateFromPassword([]byte("456"), cfg.BcryptCost)

	users := []*domain.UserDb{
		{Id: "0", Username: "User1", PasswordHash: user1password},
		{Id: "1", Username: "User2", PasswordHash: user2password},
	}

	txn := db.Txn(true)
	for _, u := range users {
		if err := txn.Insert("user", u); err != nil {
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
		return nil, &domain.UserError{Type: domain.UserAlreadyExistsError, Err: err}
	}

	id := xid.New().String()
	passwordHash, _ := bcrypt.GenerateFromPassword([]byte(newUser.Password), s.BcryptCost)
	userDb := domain.UserDb{Id: id, Username: newUser.Username, PasswordHash: passwordHash}

	txn := s.Db.Txn(true)
	defer txn.Abort()
	if err := txn.Insert("user", &newUser); err != nil {
		return nil, &domain.UserError{Type: domain.UserInternalError, Err: err}
	}
	txn.Commit()

	return userDb.Copy(), nil
}

func (s *UserService) Update(id string, newUser *domain.User) (*domain.UserDb, *domain.UserError) {

	userDb, err := s.GetById(id)
	if err != nil {
		return nil, err
	}

	_ = updateDbUser(userDb, newUser, s.BcryptCost)

	txn := s.Db.Txn(true)
	defer txn.Abort()
	if err := txn.Insert("user", &userDb); err != nil {
		return nil, &domain.UserError{Type: domain.UserInternalError, Err: err}
	}
	txn.Commit()

	return userDb.Copy(), nil
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

func updateDbUser(userDb *domain.UserDb, user *domain.User, bcryptCost int) *domain.UserError {
	if user.Username != "" {
		userDb.Username = strings.Clone(user.Username)
	}

	if user.Password != "" {
		userDb.PasswordHash, _ = bcrypt.GenerateFromPassword([]byte(user.Password), bcryptCost)
	}

	return nil
}
