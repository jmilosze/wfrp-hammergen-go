package domain

import (
	"fmt"
)

const (
	UserNotFoundError = iota
	UserAlreadyExistsError
	UserInternalError
	UserIncorrectPassword
)

type User struct {
	SharedAccounts []string
}

type UserCredentials struct {
	Username string
	Password string
}

type UserClaims struct {
	Admin bool
}

type UserRead struct {
	Id             string
	Username       string
	Admin          bool
	SharedAccounts []string
}

type UserService interface {
	GetById(id string) (*UserRead, *UserError)
	GetByName(username string) (*UserRead, *UserError)
	Authenticate(user *UserRead, password string) bool
	Create(cred *UserCredentials, user *User) (*UserRead, *UserError)
	Update(id string, user *User) (*UserRead, *UserError)
	UpdateCredentials(id string, passwd string, cred *UserCredentials) (*UserRead, *UserError)
	UpdateClaims(id string, claims *UserClaims) (*UserRead, *UserError)
	Delete(id string) *UserError
	List() ([]*UserRead, *UserError)
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
