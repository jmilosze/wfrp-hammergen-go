package mockdb

import (
	"errors"
	"github.com/jmilosze/wfrp-hammergen-go/internal/config"
	"github.com/jmilosze/wfrp-hammergen-go/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	User1      domain.User
	User2      domain.User
	BcryptCost int
}

func NewUserService(cfg *config.MockdbUserService) *UserService {
	user1password, _ := bcrypt.GenerateFromPassword([]byte("123"), cfg.BcryptCost)
	user2password, _ := bcrypt.GenerateFromPassword([]byte("456"), cfg.BcryptCost)

	user1 := domain.User{Id: "0", Username: "User1", PasswordHash: user1password}
	user2 := domain.User{Id: "1", Username: "User2", PasswordHash: user2password}
	return &UserService{User1: user1, User2: user2, BcryptCost: cfg.BcryptCost}
}

func (s *UserService) FindUserById(id string) (*domain.User, error) {
	if id == s.User1.Id {
		return &s.User1, nil
	} else if id == s.User2.Id {
		return &s.User2, nil
	} else {
		return nil, errors.New("user not found")
	}
}

func (s *UserService) FindUserByName(username string) (*domain.User, error) {
	if username == s.User1.Username {
		return &s.User1, nil
	} else if username == s.User2.Username {
		return &s.User2, nil
	} else {
		return nil, errors.New("user not found")
	}
}

func (s *UserService) Authenticate(user domain.User, password string) bool {
	if bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(password)) == nil {
		return true
	}
	return false
}
