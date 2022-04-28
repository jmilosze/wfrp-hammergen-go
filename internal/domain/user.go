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

type UserOut struct {
	Id             string
	Username       string
	Admin          bool
	SharedAccounts []string
}

type UserService interface {
	GetById(id string) (*UserOut, *UserError)
	GetByName(username string) (*UserOut, *UserError)
	Authenticate(user *UserOut, password string) bool
	Create(cred *UserCredentials, user *User) (*UserOut, *UserError)
	Update(id string, user *User) (*UserOut, *UserError)
	UpdateCredentials(id string, passwd string, cred *UserCredentials) (*UserOut, *UserError)
	UpdateClaims(id string, claims *UserClaims) (*UserOut, *UserError)
	Delete(id string) *UserError
	List() ([]*UserOut, *UserError)
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
