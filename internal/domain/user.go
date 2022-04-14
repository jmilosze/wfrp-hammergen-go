package domain

import (
	"errors"
	"fmt"
	"strings"
)

const (
	UserNotFoundError = iota
	UserAlreadyExistsError
	UserInternalError
	UserIdCannotBeUpdatedError
)

type User struct {
	Id           string `json:"id"`
	Username     string `json:"username"`
	PasswordHash []byte `json:"password_hash"`
}

func (u *User) Copy() *User {
	userCopy := *u
	userCopy.Username = strings.Clone(u.Username)
	userCopy.PasswordHash = make([]byte, len(u.PasswordHash))
	copy(userCopy.PasswordHash, u.PasswordHash)
	return &userCopy
}

func (u *User) Update(newVal *User) *UserError {
	if newVal.Id != newVal.Id {
		return &UserError{Type: UserIdCannotBeUpdatedError, Err: errors.New("user id cannot be changed")}
	}

	if newVal.Username != "" {
		u.Username = strings.Clone(newVal.Username)
	}

	if len(newVal.PasswordHash) != 0 {
		u.PasswordHash = make([]byte, len(newVal.PasswordHash))
		copy(u.PasswordHash, newVal.PasswordHash)
	}

	return nil
}

type UserService interface {
	GetById(id string) (*User, *UserError)
	GetByName(username string) (*User, *UserError)
	Authenticate(user User, password string) bool
	Create(username string, password string) (*User, *UserError)
	Update(id string, username string, password string) (*User, *UserError)
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
