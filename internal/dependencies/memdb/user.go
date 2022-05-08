package memdb

import (
	"github.com/hashicorp/go-memdb"
	"github.com/jmilosze/wfrp-hammergen-go/internal/domain"
	"github.com/rs/xid"
)

type UserDbService struct {
	Db *memdb.MemDB
}

func NewUserDbService() *UserDbService {
	db, err := createNewMemDb()
	if err != nil {
		panic(err)
	}

	return &UserDbService{Db: db}
}

func createNewMemDb() (*memdb.MemDB, error) {
	schema := &memdb.DBSchema{
		Tables: map[string]*memdb.TableSchema{
			"user": {
				Name: "user",
				Indexes: map[string]*memdb.IndexSchema{
					"id": {
						Name:    "id",
						Unique:  true,
						Indexer: &memdb.StringFieldIndex{Field: "Id"},
					},
					"username": {
						Name:    "username",
						Unique:  true,
						Indexer: &memdb.StringFieldIndex{Field: "Username"},
					},
				},
			},
		},
	}
	return memdb.NewMemDB(schema)
}

func (s *UserDbService) NewUserDb(id string) *domain.UserDb {
	var newId string
	if len(id) != 0 {
		newId = xid.New().String()
	} else {
		newId = id
	}

	return &domain.UserDb{Id: newId, Username: "", PasswordHash: []byte{}, Admin: false, SharedAccounts: []string{}}
}

func (s *UserDbService) Retrieve(fieldName string, fieldValue string) (*domain.UserDb, error) {
	txn := s.Db.Txn(false)
	userRaw, err := txn.First("user", fieldName, fieldValue)
	if err != nil {
		return nil, err
	}

	if userRaw == nil {
		return nil, nil
	}
	user := userRaw.(*domain.UserDb)

	return user.Copy(), nil
}

func (s *UserDbService) Insert(user *domain.UserDb) error {
	txn := s.Db.Txn(true)
	defer txn.Abort()
	if err := txn.Insert("user", user.Copy()); err != nil {
		return err
	}
	txn.Commit()
	return nil
}

func (s *UserDbService) Delete(id string) error {
	txn := s.Db.Txn(true)
	defer txn.Abort()
	if _, err := txn.DeleteAll("user", "id", id); err != nil {
		return &domain.UserError{Type: domain.UserInternalError, Err: err}
	}
	txn.Commit()

	return nil
}

func (s *UserDbService) List() ([]*domain.UserDb, error) {
	txn := s.Db.Txn(false)

	it, err := txn.Get("user", "id")
	if err != nil {
		return nil, err
	}

	var users []*domain.UserDb
	for obj := it.Next(); obj != nil; obj = it.Next() {
		u := obj.(*domain.UserDb)
		users = append(users, u.Copy())
	}
	return users, nil
}
