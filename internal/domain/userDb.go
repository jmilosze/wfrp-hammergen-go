package domain

import (
	"context"
	"time"
)

type UserDb struct {
	Id               string
	Username         *string
	PasswordHash     []byte
	Admin            *bool
	SharedAccountIds []string
	CreatedOn        time.Time
	LastAuthOn       time.Time
}

type UserDbService interface {
	Create(ctx context.Context, user *UserDb) *DbError
	Retrieve(ctx context.Context, fieldName string, fieldValue string) (*UserDb, *DbError)
	RetrieveMany(ctx context.Context, fieldName string, fieldValues []string) ([]*UserDb, *DbError)
	Update(ctx context.Context, user *UserDb) (*UserDb, *DbError)
	Delete(ctx context.Context, id string) *DbError
	List(ctx context.Context) ([]*UserDb, *DbError)
	NewUserDb() *UserDb
}
