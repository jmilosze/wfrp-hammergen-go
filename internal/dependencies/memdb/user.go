package memdb

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/hashicorp/go-memdb"
	"github.com/jmilosze/wfrp-hammergen-go/internal/config"
	"github.com/jmilosze/wfrp-hammergen-go/internal/domain"
	"github.com/rs/xid"
	"golang.org/x/crypto/bcrypt"
	"strings"
)

type UserService struct {
	Db             *memdb.MemDB
	BcryptCost     int
	EmailService   domain.EmailService
	JwtService     domain.JwtService
	CaptchaService domain.CaptchaService
	v              *validator.Validate
}

type UserDb struct {
	Id             string
	Username       string
	PasswordHash   []byte
	Admin          bool
	SharedAccounts []string
}

func newUserDb() *UserDb {
	newId := xid.New().String()
	return &UserDb{Id: newId, Username: "", PasswordHash: []byte{}, Admin: false, SharedAccounts: []string{}}
}

func (u *UserDb) toUser() *domain.User {
	if u == nil {
		return nil
	}

	return &domain.User{Id: u.Id, Admin: u.Admin, Username: u.Username, SharedAccounts: u.SharedAccounts}
}

func (u *UserDb) copy() *UserDb {
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

func (u *UserDb) update(sharedAccounts []string) *UserDb {
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

func (u *UserDb) updateCredentials(username string, password string, bcryptCost int) *UserDb {
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

func (u *UserDb) updateClaims(admin *bool) *UserDb {
	if u == nil {
		return nil
	}

	if admin != nil {
		u.Admin = *admin
	}

	return u
}

func NewUserService(cfg *config.MemDbUserService,
	email domain.EmailService,
	jwt domain.JwtService,
	cap domain.CaptchaService,
	v *validator.Validate) *UserService {

	db, err := createNewMemDb()
	if err != nil {
		panic(err)
	}

	for id, u := range cfg.SeedUsers {
		var userDb = UserDb{Id: id}
		userDb.updateCredentials(u.Username, u.Password, cfg.BcryptCost)
		userDb.updateClaims(&u.Admin)
		userDb.update(u.SharedAccounts)
		if err := insertInDb(&userDb, db); err != nil {
			panic(err)
		}
	}

	return &UserService{Db: db, BcryptCost: cfg.BcryptCost, EmailService: email, JwtService: jwt, CaptchaService: cap, v: v}
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
	userDb, err := getFromDb("id", id, s.Db)
	if err != nil {
		return nil, &domain.UserError{Type: domain.UserInternalError, Err: err}
	}

	if userDb == nil {
		return nil, &domain.UserError{Type: domain.UserNotFoundError, Err: errors.New("user not found")}
	}

	return userDb.toUser(), nil
}

func getFromDb(fieldName string, fieldValue string, db *memdb.MemDB) (*UserDb, error) {
	txn := db.Txn(false)
	userRaw, err := txn.First("user", fieldName, fieldValue)
	if err != nil {
		return nil, err
	}

	if userRaw == nil {
		return nil, nil
	}
	user := userRaw.(*UserDb)

	return user.copy(), nil
}

func (s *UserService) GetAndAuth(username string, passwd string) (*domain.User, *domain.UserError) {
	userDb, err := getFromDb("username", username, s.Db)
	if err != nil {
		return nil, &domain.UserError{Type: domain.UserInternalError, Err: err}
	}

	if userDb == nil {
		return nil, &domain.UserError{Type: domain.UserNotFoundError, Err: errors.New("user not found")}
	}

	if !authenticate(userDb, passwd) {
		return nil, &domain.UserError{Type: domain.UserIncorrectPassword, Err: errors.New("incorrect password")}
	}

	return userDb.toUser(), nil
}

func (s *UserService) Create(cred *domain.UserWriteCredentials, user *domain.UserWrite, captcha string) (*domain.User, *domain.UserError) {
	if len(cred.Username) == 0 || len(cred.Password) == 0 {
		return nil, &domain.UserError{Type: domain.UserInvalid, Err: errors.New("missing username or password")}
	}

	if !s.CaptchaService.Verify(captcha) {
		return nil, &domain.UserError{Type: domain.UserCaptchaFailure, Err: errors.New("captcha verification failed")}
	}

	if err := s.v.Struct(cred); err != nil {
		return nil, &domain.UserError{Type: domain.UserInvalid, Err: err}
	}

	if err := s.v.Struct(user); err != nil {
		return nil, &domain.UserError{Type: domain.UserInvalid, Err: err}
	}

	userDb, err := getFromDb("username", cred.Username, s.Db)
	if err != nil {
		return nil, &domain.UserError{Type: domain.UserInternalError, Err: err}
	}

	if userDb != nil {
		return nil, &domain.UserError{Type: domain.UserAlreadyExistsError, Err: errors.New("user already exists")}
	}

	userDb = newUserDb()
	userDb.update(user.SharedAccounts)
	userDb.updateCredentials(cred.Username, cred.Password, s.BcryptCost)

	if err := insertInDb(userDb, s.Db); err != nil {
		return nil, &domain.UserError{Type: domain.UserInternalError, Err: err}
	}

	return userDb.toUser(), nil
}

func insertInDb(u *UserDb, db *memdb.MemDB) error {

	txn := db.Txn(true)
	defer txn.Abort()
	if err := txn.Insert("user", u.copy()); err != nil {
		return err
	}
	txn.Commit()
	return nil
}

func (s *UserService) Update(id string, user *domain.UserWrite) (*domain.User, *domain.UserError) {
	if err := s.v.Struct(user); err != nil {
		return nil, &domain.UserError{Type: domain.UserInvalid, Err: err}
	}

	userDb, err := getFromDb("id", id, s.Db)
	if err != nil {
		return nil, &domain.UserError{Type: domain.UserInternalError, Err: err}
	}

	if userDb == nil {
		return nil, &domain.UserError{Type: domain.UserNotFoundError, Err: errors.New("user not found")}
	}

	userDb.update(user.SharedAccounts)

	if e := insertInDb(userDb, s.Db); e != nil {
		return nil, &domain.UserError{Type: domain.UserInternalError, Err: err}
	}

	return userDb.toUser(), nil
}

func (s *UserService) UpdateCredentials(id string, currentPasswd string, cred *domain.UserWriteCredentials) (*domain.User, *domain.UserError) {
	if err := s.v.Struct(cred); err != nil {
		return nil, &domain.UserError{Type: domain.UserInvalid, Err: err}
	}

	userDb, err := getFromDb("id", id, s.Db)
	if err != nil {
		return nil, &domain.UserError{Type: domain.UserInternalError, Err: err}
	}

	if userDb == nil {
		return nil, &domain.UserError{Type: domain.UserNotFoundError, Err: errors.New("user not found")}
	}

	if !authenticate(userDb, currentPasswd) {
		return nil, &domain.UserError{Type: domain.UserIncorrectPassword, Err: errors.New("incorrect password")}
	}

	userDb.updateCredentials(cred.Username, cred.Password, s.BcryptCost)

	if e := insertInDb(userDb, s.Db); e != nil {
		return nil, &domain.UserError{Type: domain.UserInternalError, Err: err}
	}

	return userDb.toUser(), nil
}

func authenticate(user *UserDb, password string) bool {
	if bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(password)) == nil {
		return true
	}
	return false
}

func (s *UserService) UpdateClaims(id string, claims *domain.UserWriteClaims) (*domain.User, *domain.UserError) {
	if err := s.v.Struct(claims); err != nil {
		return nil, &domain.UserError{Type: domain.UserInvalid, Err: err}
	}

	userDb, err := getFromDb("id", id, s.Db)
	if err != nil {
		return nil, &domain.UserError{Type: domain.UserInternalError, Err: err}
	}

	if userDb == nil {
		return nil, &domain.UserError{Type: domain.UserNotFoundError, Err: errors.New("user not found")}
	}

	userDb.updateClaims(claims.Admin)

	if e := insertInDb(userDb, s.Db); e != nil {
		return nil, &domain.UserError{Type: domain.UserInternalError, Err: err}
	}

	return userDb.toUser(), nil
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
		users = append(users, u.copy().toUser())
	}

	return users, nil
}

func (s *UserService) SendResetPassword(username string, captcha string) *domain.UserError {
	if !s.CaptchaService.Verify(captcha) {
		return &domain.UserError{Type: domain.UserCaptchaFailure, Err: errors.New("captcha verification failed")}
	}

	userDb, uErr := getFromDb("username", username, s.Db)
	if uErr != nil {
		return &domain.UserError{Type: domain.UserInternalError, Err: uErr}
	}

	if userDb == nil {
		return &domain.UserError{Type: domain.UserNotFoundError, Err: errors.New("user not found")}
	}

	claims := domain.Claims{Id: userDb.Id, Admin: userDb.Admin, SharedAccounts: userDb.SharedAccounts}
	resetToken, err := s.JwtService.GenerateResetPasswordToken(&claims)

	if err != nil {
		return &domain.UserError{Type: domain.UserInternalError, Err: err}
	}

	email := domain.Email{ToAddress: userDb.Username, Subject: "password reset", Content: resetToken}

	if eErr := s.EmailService.Send(&email); eErr != nil {
		return &domain.UserError{Type: domain.UserInternalError, Err: eErr}
	}

	return nil
}
