package domain

import (
	"context"
	"fmt"
	"time"
)

const (
	UserNotFoundError = iota
	UserAlreadyExistsError
	UserInternalError
	UserIncorrectPasswordError
	UserInvalidArgumentsError
	UserSendEmailError
	UserUnauthorizedError
)

type UserWrite struct {
	SharedAccounts []string `json:"shared_accounts" validate:"omitempty,dive,email,required"`
}

type UserWriteCredentials struct {
	Username string `json:"username" validate:"omitempty,email"`
	Password string `json:"password" validate:"omitempty,gte=5"`
}

type UserWriteClaims struct {
	Admin bool `json:"admin" validate:"omitempty"`
}

type User struct {
	Id             string
	Username       string
	Admin          bool
	SharedAccounts []string
	CreatedOn      time.Time
	LastAuthOn     time.Time
}

type UserService interface {
	Get(ctx context.Context, c *Claims, id string) (*User, *UserError)
	Exists(ctx context.Context, username string) (bool, *UserError)
	Create(ctx context.Context, uwc *UserWriteCredentials, uw *UserWrite) (*User, *UserError)
	Update(ctx context.Context, c *Claims, id string, uw *UserWrite) (*User, *UserError)
	UpdateCredentials(ctx context.Context, c *Claims, id string, currentPasswd string, uwc *UserWriteCredentials) (*User, *UserError)
	UpdateClaims(ctx context.Context, c *Claims, id string, uwc *UserWriteClaims) (*User, *UserError)
	Delete(ctx context.Context, c *Claims, id string) *UserError
	List(ctx context.Context, c *Claims) ([]*User, *UserError)
	Authenticate(ctx context.Context, username string, password string) (*User, *UserError)
	SendResetPassword(ctx context.Context, username string) *UserError
	ResetPassword(ctx context.Context, token string, newPassword string) *UserError
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
