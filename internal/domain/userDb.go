package domain

import (
	"context"
	"github.com/rs/xid"
	"time"
)

type UserDb struct {
	Id             string
	Username       *string
	PasswordHash   []byte
	Admin          *bool
	SharedAccounts []string
	CreatedOn      *time.Time
	LastAuthOn     *time.Time
}

func NewUserDb() *UserDb {
	newId := xid.New().String()
	admin := false
	createOn := time.Now()
	username := ""
	return &UserDb{
		Id:             newId,
		Username:       &username,
		PasswordHash:   []byte{},
		Admin:          &admin,
		SharedAccounts: []string{},
		CreatedOn:      &createOn,
		LastAuthOn:     &time.Time{},
	}
}

func (udb *UserDb) ToUser() *User {
	return &User{
		Id:             udb.Id,
		Admin:          *udb.Admin,
		Username:       *udb.Username,
		SharedAccounts: udb.SharedAccounts,
		CreatedOn:      *udb.CreatedOn,
		LastAuthOn:     *udb.LastAuthOn,
	}
}

type UserDbService interface {
	Create(ctx context.Context, user *UserDb) *DbError
	Retrieve(ctx context.Context, fieldName string, fieldValue string) (*UserDb, *DbError)
	RetrieveAll(ctx context.Context) ([]*UserDb, *DbError)
	Update(ctx context.Context, user *UserDb) (*UserDb, *DbError)
	Delete(ctx context.Context, id string) *DbError
}
