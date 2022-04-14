package domain

import (
	"fmt"
	"strings"
)

const (
	UserNotFoundError = iota
	UserAlreadyExistsError
	UserInternalError
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserDb struct {
	Id           string
	Username     string
	PasswordHash []byte
}

func (u *UserDb) Copy() *UserDb {
	userCopy := *u
	userCopy.Username = strings.Clone(u.Username)
	userCopy.PasswordHash = make([]byte, len(u.PasswordHash))
	copy(userCopy.PasswordHash, u.PasswordHash)
	return &userCopy
}

type UserService interface {
	GetById(id string) (*UserDb, *UserError)
	GetByName(username string) (*UserDb, *UserError)
	Authenticate(user UserDb, password string) bool
	Create(new *User) (*UserDb, *UserError)
	Update(id string, new *User) (*UserDb, *UserError)
	Delete(id string) *UserError
}

type UserError struct {
	Type int
	Err  error
}

func (e *UserError) Unwrap() error {
	return e.Err
}

func (e *UserError) Error() string {
	return fmt.Sprintf("user error, %s", e.Err)
}
