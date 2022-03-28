package domain

import "context"

type User struct {
	Id           string `json:"id"`
	Username     string `json:"username"`
	PasswordHash string `json:"password_hash"`
}

type UserService interface {
	FindUserById(ctx context.Context, id string) (User, error)
}
