package domain

import (
	"context"
	"github.com/rs/xid"
	"strings"
	"time"
)

type UserDb struct {
	Id                 string
	Username           string
	PasswordHash       []byte
	Admin              *bool
	SharedAccountNames []string
	SharedAccountIds   []string
	CreatedOn          time.Time
	LastAuthOn         time.Time
}

func NewUserDb() *UserDb {
	newId := xid.New().String()
	admin := false
	return &UserDb{
		Id:                 newId,
		Username:           "",
		PasswordHash:       []byte{},
		Admin:              &admin,
		SharedAccountNames: []string{},
		SharedAccountIds:   []string{},
		CreatedOn:          time.Now(),
		LastAuthOn:         time.Time{},
	}
}

func (u *UserDb) ToUser() *User {
	isAdmin := false
	if u.Admin != nil {
		isAdmin = *u.Admin
	}

	return &User{
		Id:                 u.Id,
		Admin:              isAdmin,
		Username:           u.Username,
		SharedAccountNames: u.SharedAccountNames,
		CreatedOn:          u.CreatedOn,
		LastAuthOn:         u.LastAuthOn,
	}
}

func (u *UserDb) Copy() *UserDb {
	if u == nil {
		return nil
	}
	uCopy := *u
	uCopy.Username = strings.Clone(u.Username)
	uCopy.PasswordHash = make([]byte, len(u.PasswordHash))
	copy(uCopy.PasswordHash, u.PasswordHash)
	uCopy.SharedAccountNames = make([]string, len(u.SharedAccountNames))
	for i, s := range u.SharedAccountNames {
		uCopy.SharedAccountNames[i] = strings.Clone(s)
	}
	uCopy.SharedAccountIds = make([]string, len(u.SharedAccountIds))
	for i, s := range u.SharedAccountIds {
		uCopy.SharedAccountIds[i] = strings.Clone(s)
	}

	uCopy.LastAuthOn = u.LastAuthOn.UTC()
	uCopy.CreatedOn = u.CreatedOn.UTC()

	return &uCopy
}

type UserDbService interface {
	Create(ctx context.Context, user *UserDb) (*UserDb, *DbError)
	Retrieve(ctx context.Context, fieldName string, fieldValue string) (*UserDb, *DbError)
	RetrieveAll(ctx context.Context) ([]*UserDb, *DbError)
	Update(ctx context.Context, user *UserDb) (*UserDb, *DbError)
	Delete(ctx context.Context, id string) *DbError
}
