package domain

import (
	"fmt"
)

const (
	UserNotFoundError = iota
	UserAlreadyExistsError
	UserInternalError
	UserIncorrectPassword
	UserInvalid
)

type UserWrite struct {
	SharedAccounts []string `validate:"omitempty"`
}

type UserWriteCredentials struct {
	Username string `validate:"omitempty,email"`
	Password string `validate:"omitempty,gte=5"`
}

type UserWriteClaims struct {
	Admin *bool `validate:"omitempty"`
}

type User struct {
	Id             string
	Username       string
	Admin          bool
	SharedAccounts []string
}

//func (c *UserWriteCredentials) Validate() error {
//	validate := validator.New()
//	return validate.Struct(c)
//}

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
