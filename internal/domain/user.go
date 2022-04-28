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

type UserWrite struct {
	SharedAccounts []string
}

type UserWriteCredentials struct {
	Username string
	Password string
}

type UserWriteClaims struct {
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
	Create(cred *UserWriteCredentials, user *UserWrite) (*UserRead, *UserError)
	Update(id string, user *UserWrite) (*UserRead, *UserError)
	UpdateCredentials(id string, passwd string, cred *UserWriteCredentials) (*UserRead, *UserError)
	UpdateClaims(id string, claims *UserWriteClaims) (*UserRead, *UserError)
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
