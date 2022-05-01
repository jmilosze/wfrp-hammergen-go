package memdb

import (
	"errors"
	"github.com/hashicorp/go-memdb"
	"github.com/jmilosze/wfrp-hammergen-go/internal/config"
	"github.com/jmilosze/wfrp-hammergen-go/internal/domain"
	"github.com/rs/xid"
	"golang.org/x/crypto/bcrypt"
	"strings"
)

type UserService struct {
	Db           *memdb.MemDB
	BcryptCost   int
	EmailService domain.EmailService
	JwtService   domain.JwtService
}

type UserDb struct {
	Id             string
	Username       string
	PasswordHash   []byte
	Admin          bool
	SharedAccounts []string
}

func (u *UserDb) ToUser() *domain.User {
	if u == nil {
		return nil
	}

	return &domain.User{Id: u.Id, Admin: u.Admin, Username: u.Username, SharedAccounts: u.SharedAccounts}
}

func (u *UserDb) Copy() *UserDb {
	if u == nil {
		return nil
	}
	userCopy := *u
	userCopy.Username = strings.Clone(u.Username)
	userCopy.PasswordHash = make([]byte, len(u.PasswordHash))
	copy(userCopy.PasswordHash, u.PasswordHash)
	userCopy.SharedAccounts = make([]string, len(u.SharedAccounts))
	for i, s := range u.SharedAccounts {
		userCopy.SharedAccounts[i] = strings.Clone(s)
	}

	return &userCopy
}

func (u *UserDb) UpdateCredentials(username string, password string, bcryptCost int) *UserDb {
	if u == nil {
		return nil
	}

	if len(username) != 0 {
		u.Username = strings.Clone(username)
	}

	if len(password) != 0 {
		u.PasswordHash, _ = bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	}

	return u
}

func (u *UserDb) UpdateClaims(admin *bool) *UserDb {
	if u == nil {
		return nil
	}

	if admin != nil {
		u.Admin = *admin
	}

	return u
}

func (u *UserDb) Update(sharedAccounts []string) *UserDb {
	if u == nil {
		return nil
	}

	if sharedAccounts != nil {
		u.SharedAccounts = make([]string, len(sharedAccounts))
		for i, s := range sharedAccounts {
			u.SharedAccounts[i] = strings.Clone(s)
		}
	}

	return u
}

func NewUserService(cfg *config.MockDbUserService, users map[string]*config.UserSeed, email domain.EmailService, jwt domain.JwtService) *UserService {
	db, err := createNewMemDb()
	if err != nil {
		panic(err)
	}

	txn := db.Txn(true)
	for id, u := range users {
		var userDb = UserDb{Id: id}
		userDb.UpdateCredentials(u.Username, u.Password, cfg.BcryptCost)
		userDb.UpdateClaims(&u.Admin)
		userDb.Update(u.SharedAccounts)
		if err := txn.Insert("user", &userDb); err != nil {
			panic(err)
		}
	}

	txn.Commit()

	return &UserService{Db: db, BcryptCost: cfg.BcryptCost, EmailService: email, JwtService: jwt}
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

func (s *UserService) Get(id string) (*domain.User, *domain.UserError) {
	user, err := getFromDb("id", id, s.Db)
	return user.ToUser(), err
}

func getFromDb(fieldName string, fieldValue string, db *memdb.MemDB) (*UserDb, *domain.UserError) {
	txn := db.Txn(false)
	userRaw, err := txn.First("user", fieldName, fieldValue)
	if err != nil {
		return nil, &domain.UserError{Type: domain.UserInternalError, Err: err}
	}

	if userRaw == nil {
		return nil, &domain.UserError{Type: domain.UserNotFoundError, Err: errors.New("user not found")}
	}
	user := userRaw.(*UserDb)

	return user.Copy(), nil
}

func (s *UserService) GetAndAuth(username string, passwd string) (*domain.User, *domain.UserError) {
	userDb, err := getFromDb("username", username, s.Db)
	if err != nil {
		return nil, err
	}

	if !authenticate(userDb, passwd) {
		return nil, &domain.UserError{Type: domain.UserIncorrectPassword, Err: errors.New("incorrect password")}
	}

	return userDb.ToUser(), nil
}

func (s *UserService) Create(sharedAccounts []string, username string, password string) (*domain.User, *domain.UserError) {
	if _, err := getFromDb("username", username, s.Db); err == nil {
		return nil, &domain.UserError{Type: domain.UserAlreadyExistsError, Err: errors.New("user already exists")}
	}

	newId := xid.New().String()
	userDb := UserDb{Id: newId, Admin: false}
	userDb.Update(sharedAccounts)
	userDb.UpdateCredentials(username, password, s.BcryptCost)

	if err := insertInDb(&userDb, s.Db); err != nil {
		return nil, err
	}

	return (&userDb).ToUser(), nil
}

func insertInDb(u *UserDb, db *memdb.MemDB) *domain.UserError {

	txn := db.Txn(true)
	defer txn.Abort()
	if err := txn.Insert("user", u.Copy()); err != nil {
		return &domain.UserError{Type: domain.UserInternalError, Err: err}
	}
	txn.Commit()
	return nil
}

func (s *UserService) Update(id string, sharedAccounts []string) (*domain.User, *domain.UserError) {
	userDb, err := getFromDb("id", id, s.Db)
	if err != nil {
		return nil, err
	}

	userDb.Update(sharedAccounts)

	if e := insertInDb(userDb, s.Db); e != nil {
		return nil, err
	}

	return userDb.ToUser(), nil
}

func (s *UserService) UpdateCredentials(id string, currentPasswd string, username string, password string) (*domain.User, *domain.UserError) {
	userDb, err := getFromDb("id", id, s.Db)
	if err != nil {
		return nil, err
	}

	if !authenticate(userDb, currentPasswd) {
		return nil, &domain.UserError{Type: domain.UserIncorrectPassword, Err: errors.New("incorrect password")}
	}

	userDb.UpdateCredentials(username, password, s.BcryptCost)

	if e := insertInDb(userDb, s.Db); e != nil {
		return nil, err
	}

	return userDb.ToUser(), nil
}

func authenticate(user *UserDb, password string) bool {
	if bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(password)) == nil {
		return true
	}
	return false
}

func (s *UserService) UpdateClaims(id string, admin *bool) (*domain.User, *domain.UserError) {
	userDb, err := getFromDb("id", id, s.Db)
	if err != nil {
		return nil, err
	}

	userDb.UpdateClaims(admin)

	if e := insertInDb(userDb, s.Db); e != nil {
		return nil, err
	}

	return userDb.ToUser(), nil
}

func (s *UserService) Delete(id string) *domain.UserError {
	txn := s.Db.Txn(true)
	defer txn.Abort()
	if _, err := txn.DeleteAll("user", "id", id); err != nil {
		return &domain.UserError{Type: domain.UserInternalError, Err: err}
	}
	txn.Commit()

	return nil
}

func (s *UserService) List() ([]*domain.User, *domain.UserError) {
	txn := s.Db.Txn(false)

	it, err := txn.Get("user", "id")
	if err != nil {
		return nil, &domain.UserError{Type: domain.UserInternalError, Err: err}
	}

	var users []*domain.User
	for obj := it.Next(); obj != nil; obj = it.Next() {
		u := obj.(*UserDb)
		users = append(users, u.Copy().ToUser())
	}

	return users, nil
}
