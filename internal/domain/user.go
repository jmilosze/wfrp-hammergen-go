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
	Id             string
	Username       string
	Admin          bool
	SharedAccounts []string
}

type UserService interface {
	Get(id string) (*User, *UserError)
	Create(sharedAccounts []string, username string, password string) (*User, *UserError)
	Update(id string, sharedAccounts []string) (*User, *UserError)
	UpdateCredentials(id string, currentPasswd string, username string, password string) (*User, *UserError)
	UpdateClaims(id string, admin *bool) (*User, *UserError)
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
