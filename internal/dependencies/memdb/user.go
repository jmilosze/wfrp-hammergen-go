package memdb

import (
	"errors"
	"fmt"
	"github.com/hashicorp/go-memdb"
	"github.com/jmilosze/wfrp-hammergen-go/internal/config"
	"github.com/jmilosze/wfrp-hammergen-go/internal/domain"
	"github.com/rs/xid"
	"golang.org/x/crypto/bcrypt"
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

	users := []*domain.User{
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

func (s *UserService) FindUserById(id string) (*domain.User, error) {
	return findUserBy("id", id, s)
}

func (s *UserService) FindUserByName(username string) (*domain.User, error) {
	return findUserBy("username", username, s)
}

func findUserBy(fieldName string, fieldValue string, s *UserService) (*domain.User, error) {
	txn := s.Db.Txn(false)

	userRaw, err := txn.First("user", fieldName, fieldValue)
	if err != nil {
		panic(err)
	}

	if userRaw == nil {
		return nil, errors.New("user not found")
	}
	user := userRaw.(*domain.User)

	return user.Copy(), nil
}

func (s *UserService) Authenticate(user domain.User, password string) bool {
	if bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(password)) == nil {
		return true
	}
	return false
}

func (s *UserService) CreateUser(username string, password string) (*domain.User, error) {
	if _, err := findUserBy("username", username, s); err == nil {
		return nil, errors.New("user already exists")
	}

	id := xid.New().String()
	passwordHash, _ := bcrypt.GenerateFromPassword([]byte(password), s.BcryptCost)
	user := domain.User{Id: id, Username: username, PasswordHash: passwordHash}

	txn := s.Db.Txn(true)
	defer txn.Abort()
	if err := txn.Insert("user", &user); err != nil {
		return nil, fmt.Errorf("error inserting userrname: %s, id: %s", username, id)
	}
	txn.Commit()

	return user.Copy(), nil
}
