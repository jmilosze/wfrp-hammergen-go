package domain

import (
	"strings"
)

type UserDb struct {
	Id             string
	Username       *string
	PasswordHash   []byte
	Admin          *bool
	SharedAccounts []string
}

func (u *UserDb) ToUser() *User {
	if u == nil {
		return nil
	}
	return &User{Id: u.Id, Admin: u.Admin, Username: u.Username, SharedAccounts: u.SharedAccounts}
}

func (u *UserDb) Copy() *UserDb {
	if u == nil {
		return nil
	}
	userCopy := *u
	*userCopy.Username = strings.Clone(*u.Username)
	userCopy.PasswordHash = make([]byte, len(u.PasswordHash))
	copy(userCopy.PasswordHash, u.PasswordHash)
	userCopy.SharedAccounts = make([]string, len(u.SharedAccounts))
	for i, s := range u.SharedAccounts {
		userCopy.SharedAccounts[i] = strings.Clone(s)
	}

	return &userCopy
}

type UserDbService interface {
	Create(user *UserDb) *DbError
	Retrieve(fieldName string, fieldValue string) (*UserDb, *DbError)
	Update(user *UserDb) (*UserDb, *DbError)
	Delete(id string) *DbError
	List() ([]*UserDb, *DbError)
	NewUserDb() *UserDb
}
