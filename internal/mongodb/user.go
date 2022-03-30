package mongodb

import (
	"errors"
	"github.com/jmilosze/wfrp-hammergen-go/internal/domain"
)

type UserService struct {
	user1 domain.User
	user2 domain.User
}

func NewUserService() *UserService {
	user1 := domain.User{Id: "0", Username: "User1", PasswordHash: "123"}
	user2 := domain.User{Id: "1", Username: "User2", PasswordHash: "456"}
	return &UserService{user1: user1, user2: user2}
}

func (s *UserService) FindUserById(id string) (*domain.User, error) {
	if id == s.user1.Id {
		return &s.user1, nil
	} else if id == s.user2.Id {
		return &s.user2, nil
	} else {
		return nil, errors.New("user not found")
	}
}
