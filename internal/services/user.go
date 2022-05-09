package services

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/jmilosze/wfrp-hammergen-go/internal/config"
	"github.com/jmilosze/wfrp-hammergen-go/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	BcryptCost    int
	Validator     *validator.Validate
	UserDbService domain.UserDbService
	EmailService  domain.EmailService
	JwtService    domain.JwtService
}

func toUser(u *domain.UserDb) *domain.User {
	if u == nil {
		return nil
	}
	return &domain.User{Id: u.Id, Admin: u.Admin, Username: u.Username, SharedAccounts: u.SharedAccounts}
}

func NewUserService(cfg *config.UserServiceConfig,
	db domain.UserDbService,
	email domain.EmailService,
	jwt domain.JwtService,
	v *validator.Validate) *UserService {

	for id, u := range cfg.SeedUsers {
		userDb := db.NewUserDb()

		userDb.Id = id
		userDb.Username = &u.Username
		userDb.PasswordHash, _ = bcrypt.GenerateFromPassword([]byte(u.Password), cfg.BcryptCost)
		userDb.SharedAccounts = u.SharedAccounts
		userDb.Admin = &u.Admin

		if err := db.Create(userDb); err != nil {
			panic(err)
		}
	}

	return &UserService{BcryptCost: cfg.BcryptCost, UserDbService: db, EmailService: email, JwtService: jwt, Validator: v}
}

func (s *UserService) Get(id string) (*domain.User, *domain.UserError) {
	userDb, err := s.UserDbService.Retrieve("id", id)
	if err != nil {
		switch err.Type {
		case domain.DbNotFoundError:
			return nil, &domain.UserError{Type: domain.UserInternalError, Err: err}
		default:
			return nil, &domain.UserError{Type: domain.UserInternalError, Err: err}
		}
	}

	return toUser(userDb), nil
}

func (s *UserService) Authenticate(username string, password string) (*domain.User, *domain.UserError) {
	userDb, err := s.UserDbService.Retrieve("username", username)
	if err != nil {
		switch err.Type {
		case domain.DbNotFoundError:
			return nil, &domain.UserError{Type: domain.UserNotFoundError, Err: err}
		default:
			return nil, &domain.UserError{Type: domain.UserInternalError, Err: err}
		}
	}

	if !authenticate(userDb, password) {
		return nil, &domain.UserError{Type: domain.UserIncorrectPassword, Err: errors.New("incorrect password")}
	}

	return toUser(userDb), nil
}

func authenticate(user *domain.UserDb, password string) bool {
	if bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(password)) == nil {
		return true
	}
	return false
}

func (s *UserService) Create(cred *domain.UserWriteCredentials, user *domain.UserWrite) (*domain.User, *domain.UserError) {
	if len(cred.Username) == 0 || len(cred.Password) == 0 {
		return nil, &domain.UserError{Type: domain.UserInvalidArguments, Err: errors.New("missing username or password")}
	}

	if err := s.Validator.Struct(cred); err != nil {
		return nil, &domain.UserError{Type: domain.UserInvalidArguments, Err: err}
	}
	if err := s.Validator.Struct(user); err != nil {
		return nil, &domain.UserError{Type: domain.UserInvalidArguments, Err: err}
	}

	userDb := s.UserDbService.NewUserDb()
	userDb.Username = &cred.Username
	userDb.PasswordHash, _ = bcrypt.GenerateFromPassword([]byte(cred.Password), s.BcryptCost)
	userDb.SharedAccounts = user.SharedAccounts

	if err := s.UserDbService.Create(userDb); err != nil {
		switch err.Type {
		case domain.DbAlreadyExistsError:
			return nil, &domain.UserError{Type: domain.UserAlreadyExistsError, Err: err}
		default:
			return nil, &domain.UserError{Type: domain.UserInternalError, Err: err}
		}
	}

	return toUser(userDb), nil
}

func (s *UserService) Update(id string, user *domain.UserWrite) (*domain.User, *domain.UserError) {
	if err := s.Validator.Struct(user); err != nil {
		return nil, &domain.UserError{Type: domain.UserInvalidArguments, Err: err}
	}

	userUpdate := domain.UserDb{Id: id, SharedAccounts: user.SharedAccounts}

	userDb, err := s.UserDbService.Update(&userUpdate)

	if err != nil {
		switch err.Type {
		case domain.DbNotFoundError:
			return nil, &domain.UserError{Type: domain.UserNotFoundError, Err: err}
		default:
			return nil, &domain.UserError{Type: domain.UserInternalError, Err: err}
		}
	}

	return toUser(userDb), nil
}

func (s *UserService) UpdateCredentials(id string, currentPasswd string, cred *domain.UserWriteCredentials) (*domain.User, *domain.UserError) {
	if err := s.Validator.Struct(cred); err != nil {
		return nil, &domain.UserError{Type: domain.UserInvalidArguments, Err: err}
	}

	userDb, err := s.UserDbService.Retrieve("id", id)
	if err != nil {
		switch err.Type {
		case domain.DbNotFoundError:
			return nil, &domain.UserError{Type: domain.UserNotFoundError, Err: err}
		default:
			return nil, &domain.UserError{Type: domain.UserInternalError, Err: err}
		}
	}

	if !authenticate(userDb, currentPasswd) {
		return nil, &domain.UserError{Type: domain.UserIncorrectPassword, Err: errors.New("incorrect password")}
	}

	userDb.Username = &cred.Username
	userDb.PasswordHash, _ = bcrypt.GenerateFromPassword([]byte(cred.Password), s.BcryptCost)

	if _, err := s.UserDbService.Update(userDb); err != nil {
		switch err.Type {
		case domain.DbNotFoundError:
			return nil, &domain.UserError{Type: domain.UserNotFoundError, Err: err}
		default:
			return nil, &domain.UserError{Type: domain.UserInternalError, Err: err}
		}
	}

	return toUser(userDb), nil
}

func (s *UserService) UpdateClaims(id string, claims *domain.UserWriteClaims) (*domain.User, *domain.UserError) {
	if err := s.Validator.Struct(claims); err != nil {
		return nil, &domain.UserError{Type: domain.UserInvalidArguments, Err: err}
	}

	userUpdate := domain.UserDb{Id: id, Admin: claims.Admin}
	userDb, err := s.UserDbService.Update(&userUpdate)
	if err != nil {
		switch err.Type {
		case domain.DbNotFoundError:
			return nil, &domain.UserError{Type: domain.UserNotFoundError, Err: err}
		default:
			return nil, &domain.UserError{Type: domain.UserInternalError, Err: err}
		}
	}

	return toUser(userDb), nil
}

func (s *UserService) Delete(id string) *domain.UserError {
	if err := s.UserDbService.Delete(id); err != nil {
		return &domain.UserError{Type: domain.UserInternalError, Err: err}

	} else {
		return nil
	}
}

func (s *UserService) List() ([]*domain.User, *domain.UserError) {
	usersDb, err := s.UserDbService.List()

	if err != nil {
		return nil, &domain.UserError{Type: domain.UserInternalError, Err: err}

	}

	users := make([]*domain.User, len(usersDb))
	for i, udb := range usersDb {
		users[i] = toUser(udb)
	}
	return users, nil
}

func (s *UserService) SendResetPassword(username string) *domain.UserError {
	if len(username) == 0 {
		return &domain.UserError{Type: domain.UserInvalidArguments, Err: errors.New("missing username")}
	}

	userDb, uErr := s.UserDbService.Retrieve("username", username)
	if uErr != nil {
		return &domain.UserError{Type: domain.UserInternalError, Err: uErr}
	}

	if userDb == nil {
		return &domain.UserError{Type: domain.UserNotFoundError, Err: errors.New("user not found")}
	}

	claims := domain.Claims{Id: userDb.Id, Admin: *userDb.Admin, SharedAccounts: userDb.SharedAccounts, ResetPassword: true}
	resetToken, err := s.JwtService.GenerateResetPasswordToken(&claims)

	if err != nil {
		return &domain.UserError{Type: domain.UserInternalError, Err: err}
	}

	email := domain.Email{
		ToAddress: *userDb.Username,
		Subject:   "password reset",
		Content:   resetToken,
	}

	if eErr := s.EmailService.Send(&email); eErr != nil {
		return &domain.UserError{Type: domain.UserInternalError, Err: eErr}
	}

	return nil
}

func (s *UserService) ResetPassword(token string, newPassword string) *domain.UserError {
	if len(token) == 0 || len(newPassword) == 0 {
		return &domain.UserError{Type: domain.UserInvalidArguments, Err: errors.New("missing token or username")}
	}

	claims, err1 := s.JwtService.ParseToken(token)
	if err1 != nil {
		return &domain.UserError{Type: domain.UserInvalidArguments, Err: errors.New("invalid token")}
	}

	if !claims.ResetPassword {
		return &domain.UserError{Type: domain.UserInvalidArguments, Err: errors.New("invalid token")}
	}

	var newCreds = domain.UserWriteCredentials{Username: "", Password: newPassword}
	if err2 := s.Validator.Struct(newCreds); err2 != nil {
		return &domain.UserError{Type: domain.UserInvalidArguments, Err: err2}
	}

	newHash, _ := bcrypt.GenerateFromPassword([]byte(newCreds.Password), s.BcryptCost)
	userUpdate := domain.UserDb{Id: claims.Id, PasswordHash: newHash}

	if _, err4 := s.UserDbService.Update(&userUpdate); err4 != nil {
		switch err4.Type {
		case domain.DbNotFoundError:
			return &domain.UserError{Type: domain.UserNotFoundError, Err: err4}
		default:
			return &domain.UserError{Type: domain.UserInternalError, Err: err4}
		}
	}

	return nil
}
