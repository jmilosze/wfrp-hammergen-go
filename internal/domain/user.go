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
	Admin *bool
}

type User struct {
	Id             string
	Username       string
	Admin          bool
	SharedAccounts []string
}

type UserService interface {
	Get(id string) (*User, *UserError)
	Create(cred *UserWriteCredentials, user *UserWrite) (*User, *UserError)
	Update(id string, user *UserWrite) (*User, *UserError)
	UpdateCredentials(id string, currentPasswd string, cred *UserWriteCredentials) (*User, *UserError)
	UpdateClaims(id string, claims *UserWriteClaims) (*User, *UserError)
	Delete(id string) *UserError
	List() ([]*User, *UserError)
	GetAndAuth(username string, passwd string) (*User, *UserError)
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
