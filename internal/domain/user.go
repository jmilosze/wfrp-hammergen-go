package domain

import "fmt"

const (
	UserNotFoundError = iota
	UserAlreadyExistsError
	UserInternalError
)

type User struct {
	Id           string `json:"id"`
	Username     string `json:"username"`
	PasswordHash []byte `json:"password_hash"`
}

func (source *User) Copy() *User {
	userCopy := *source
	userCopy.PasswordHash = make([]byte, len(source.PasswordHash))
	copy(userCopy.PasswordHash, source.PasswordHash)
	return &userCopy
}

type UserService interface {
	GetById(id string) (*User, *UserError)
	GetByName(username string) (*User, *UserError)
	Authenticate(user User, password string) bool
	Create(username string, password string) (*User, *UserError)
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
