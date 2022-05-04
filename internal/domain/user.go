package domain

import (
	"fmt"
)

const (
	UserNotFoundError = iota
	UserAlreadyExistsError
	UserInternalError
	UserIncorrectPassword
	UserInvalidArguments
	UserCaptchaFailure
)

type UserWrite struct {
	SharedAccounts []string `json:"shared_accounts" validate:"omitempty,dive,required"`
}

type UserWriteCredentials struct {
	Username string `json:"username" validate:"omitempty,email"`
	Password string `json:"password" validate:"omitempty,gte=5"`
}

type UserWriteClaims struct {
	Admin *bool `json:"admin" validate:"omitempty"`
}

type User struct {
	Id             string
	Username       string
	Admin          bool
	SharedAccounts []string
}

type UserService interface {
	Get(id string) (*User, *UserError)
	Create(cred *UserWriteCredentials, user *UserWrite, captcha string) (*User, *UserError)
	Update(id string, user *UserWrite) (*User, *UserError)
	UpdateCredentials(id string, currentPasswd string, cred *UserWriteCredentials) (*User, *UserError)
	UpdateClaims(id string, claims *UserWriteClaims) (*User, *UserError)
	Delete(id string) *UserError
	List() ([]*User, *UserError)
	GetAndAuth(username string, passwd string) (*User, *UserError)
	SendResetPassword(username string, captcha string) *UserError
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
