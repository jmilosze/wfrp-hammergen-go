package memdb

import (
	"errors"
	"github.com/hashicorp/go-memdb"
	"github.com/jmilosze/wfrp-hammergen-go/internal/config"
	"github.com/jmilosze/wfrp-hammergen-go/internal/domain"
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
	defer txn.Abort()

	userRaw, err := txn.First("user", fieldName, fieldValue)
	if err != nil {
		panic(err)
	}

	if userRaw == nil {
		return nil, errors.New("user not found")
	}

	user := userRaw.(*domain.User)

	userCopy := *user
	userCopy.PasswordHash = make([]byte, len(user.PasswordHash))
	copy(userCopy.PasswordHash, user.PasswordHash)

	return &userCopy, nil
}

func (s *UserService) Authenticate(user domain.User, password string) bool {
	if bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(password)) == nil {
		return true
	}
	return false
}
