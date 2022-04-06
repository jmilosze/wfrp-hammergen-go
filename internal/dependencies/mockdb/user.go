package mockdb

import (
	"errors"
	"github.com/jmilosze/wfrp-hammergen-go/internal/config"
	"github.com/jmilosze/wfrp-hammergen-go/internal/domain"
	"golang.org/x/crypto/bcrypt"
	"sync"
)

type userStore map[string]domain.User

type UserService struct {
	users      userStore
	BcryptCost int
}

func NewUserService(cfg *config.MockdbUserService) *UserService {
	user1password, _ := bcrypt.GenerateFromPassword([]byte("123"), cfg.BcryptCost)
	user2password, _ := bcrypt.GenerateFromPassword([]byte("456"), cfg.BcryptCost)

	user1 := domain.User{Id: "0", Username: "User1", PasswordHash: user1password}
	user2 := domain.User{Id: "1", Username: "User2", PasswordHash: user2password}

	users := userStore{user1.Id: user1, user2.Id: user2}

	return &UserService{users: users, BcryptCost: cfg.BcryptCost}
}

func (s *UserService) FindUserById(id string) (*domain.User, error) {
	var m sync.RWMutex

	m.Lock()
	user, ok := s.users[id]
	m.Unlock()

	if ok {
		return &user, nil
	} else {
		return nil, errors.New("user not found")
	}
}

func (s *UserService) FindUserByName(username string) (*domain.User, error) {
	var m sync.RWMutex
	var userFound *domain.User = nil

	m.Lock()
	for _, user := range s.users {
		if user.Username == username {
			userCopy := user
			userCopy.PasswordHash = make([]byte, len(user.PasswordHash))
			copy(userCopy.PasswordHash, user.PasswordHash)
			userFound = &userCopy
			break
		}
	}
	m.Unlock()

	if userFound != nil {
		return userFound, nil
	}
	return nil, errors.New("user not found")
}

func (s *UserService) Authenticate(user domain.User, password string) bool {
	if bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(password)) == nil {
		return true
	}
	return false
}
