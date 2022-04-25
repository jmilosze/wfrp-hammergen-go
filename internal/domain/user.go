package domain

import (
	"fmt"
	"strings"
)

const (
	UserNotFoundError = iota
	UserAlreadyExistsError
	UserInternalError
	UserIncorrectPassword
)

type User struct {
	Username       string   `json:"username"`
	Password       string   `json:"password"`
	SharedAccounts []string `json:"shared_accounts"`
	Admin          bool     `json:"admin"`
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

type UserService interface {
	GetById(id string) (*UserDb, *UserError)
	GetByName(username string) (*UserDb, *UserError)
	Authenticate(user *UserDb, password string) bool
	Create(new *User) (*UserDb, *UserError)
	Update(id string, new *User) (*UserDb, *UserError)
	UpdateUsername(id string, passwd string, newUsername string) (*UserDb, *UserError)
	UpdatePassword(id string, passwd string, newPasswd string) (*UserDb, *UserError)
	Delete(id string) *UserError
	List() ([]*UserDb, *UserError)
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
