package services

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/jmilosze/wfrp-hammergen-go/internal/config"
	"github.com/jmilosze/wfrp-hammergen-go/internal/domain"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/url"
	"path"
	"time"
)

type UserService struct {
	BcryptCost    int
	Validator     *validator.Validate
	UserDbService domain.UserDbService
	EmailService  domain.EmailService
	JwtService    domain.JwtService
	FrontEndUrl   *url.URL
}

func NewUserService(cfg *config.UserService, db domain.UserDbService, email domain.EmailService, jwt domain.JwtService, v *validator.Validate) *UserService {
	return &UserService{
		BcryptCost:    cfg.BcryptCost,
		UserDbService: db,
		EmailService:  email,
		JwtService:    jwt,
		Validator:     v,
		FrontEndUrl:   cfg.FrontEndUrl,
	}

}

func (s *UserService) SeedUsers(ctx context.Context, su []*config.UserSeed) {
	for _, u := range su {
		userDb := domain.NewUserDb()

		userDb.Id = u.Id
		userDb.Username = u.Username
		userDb.PasswordHash, _ = bcrypt.GenerateFromPassword([]byte(u.Password), s.BcryptCost)
		userDb.SharedAccountNames = u.SharedAccounts
		userDb.Admin = &u.Admin

		if _, err := s.UserDbService.Create(ctx, userDb); err != nil {
			log.Fatal(err)
		}
	}
}

func (s *UserService) Get(ctx context.Context, c *domain.Claims, id string) (*domain.User, *domain.UserError) {
	if c.Id == "anonymous" || !(id == c.Id || c.Admin) {
		return nil, &domain.UserError{Type: domain.UserUnauthorizedError, Err: errors.New("unauthorized")}
	}

	userDb, err := s.UserDbService.Retrieve(ctx, "id", id)
	if err != nil {
		switch err.Type {
		case domain.DbNotFoundError:
			return nil, &domain.UserError{Type: domain.UserNotFoundError, Err: err}
		default:
			return nil, &domain.UserError{Type: domain.UserInternalError, Err: err}
		}
	}

	return userDb.ToUser(), nil
}

func (s *UserService) Exists(ctx context.Context, username string) (bool, *domain.UserError) {
	_, err := s.UserDbService.Retrieve(ctx, "username", username)
	if err != nil {
		switch err.Type {
		case domain.DbNotFoundError:
			return false, nil
		default:
			return false, &domain.UserError{Type: domain.UserInternalError, Err: err}
		}
	}
	return true, nil
}

func (s *UserService) Authenticate(ctx context.Context, username string, password string) (u *domain.User, sharedAccountIds []string, ue *domain.UserError) {
	userDb, err1 := s.UserDbService.Retrieve(ctx, "username", username)
	if err1 != nil {
		switch err1.Type {
		case domain.DbNotFoundError:
			return nil, nil, &domain.UserError{Type: domain.UserNotFoundError, Err: err1}
		default:
			return nil, nil, &domain.UserError{Type: domain.UserInternalError, Err: err1}
		}
	}

	if !authenticate(userDb, password) {
		return nil, nil, &domain.UserError{Type: domain.UserIncorrectPasswordError, Err: errors.New("incorrect password")}
	}

	if _, err2 := s.UserDbService.Update(ctx, &domain.UserDb{Id: userDb.Id, LastAuthOn: time.Now()}); err2 != nil {
		switch err2.Type {
		case domain.DbNotFoundError:
			return nil, nil, &domain.UserError{Type: domain.UserNotFoundError, Err: err2}
		default:
			return nil, nil, &domain.UserError{Type: domain.UserInternalError, Err: err2}
		}
	}

	return userDb.ToUser(), userDb.SharedAccountIds, nil
}

func authenticate(user *domain.UserDb, password string) (success bool) {
	if bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(password)) == nil {
		return true
	}
	return false
}

func (s *UserService) Create(ctx context.Context, uwc *domain.UserWriteCredentials, uw *domain.UserWrite) (*domain.User, *domain.UserError) {
	if len(uwc.Username) == 0 || len(uwc.Password) == 0 {
		return nil, &domain.UserError{Type: domain.UserInvalidArgumentsError, Err: errors.New("missing username or password")}
	}

	if err := s.Validator.Struct(uwc); err != nil {
		return nil, &domain.UserError{Type: domain.UserInvalidArgumentsError, Err: err}
	}
	if err := s.Validator.Struct(uw); err != nil {
		return nil, &domain.UserError{Type: domain.UserInvalidArgumentsError, Err: err}
	}

	passwordHash, err1 := bcrypt.GenerateFromPassword([]byte(uwc.Password), s.BcryptCost)
	if err1 != nil {
		return nil, &domain.UserError{Type: domain.UserInternalError, Err: err1}
	}

	userDb := domain.NewUserDb()
	userDb.Username = uwc.Username
	userDb.PasswordHash = passwordHash
	userDb.SharedAccountNames = uw.SharedAccountNames

	createdUserDb, err2 := s.UserDbService.Create(ctx, userDb)
	if err2 != nil {
		switch err2.Type {
		case domain.DbAlreadyExistsError:
			return nil, &domain.UserError{Type: domain.UserAlreadyExistsError, Err: err2}
		default:
			return nil, &domain.UserError{Type: domain.UserInternalError, Err: err2}
		}
	}

	return createdUserDb.ToUser(), nil
}

func (s *UserService) Update(ctx context.Context, c *domain.Claims, id string, uw *domain.UserWrite) (*domain.User, *domain.UserError) {
	if c.Id == "anonymous" || id != c.Id {
		return nil, &domain.UserError{Type: domain.UserUnauthorizedError, Err: errors.New("unauthorized")}
	}

	if err := s.Validator.Struct(uw); err != nil {
		return nil, &domain.UserError{Type: domain.UserInvalidArgumentsError, Err: err}
	}

	sharedAccounts := make([]string, 0)
	if uw.SharedAccountNames != nil {
		sharedAccounts = uw.SharedAccountNames
	}

	userUpdate := domain.UserDb{
		Id:                 id,
		SharedAccountNames: sharedAccounts,
	}

	userDb, err := s.UserDbService.Update(ctx, &userUpdate)

	if err != nil {
		switch err.Type {
		case domain.DbNotFoundError:
			return nil, &domain.UserError{Type: domain.UserNotFoundError, Err: err}
		default:
			return nil, &domain.UserError{Type: domain.UserInternalError, Err: err}
		}
	}

	return userDb.ToUser(), nil
}

func (s *UserService) UpdateCredentials(ctx context.Context, c *domain.Claims, id string, currentPasswd string, uwc *domain.UserWriteCredentials) (*domain.User, *domain.UserError) {
	if c.Id == "anonymous" || id != c.Id {
		return nil, &domain.UserError{Type: domain.UserUnauthorizedError, Err: errors.New("unauthorized")}
	}

	if err := s.Validator.Struct(uwc); err != nil {
		return nil, &domain.UserError{Type: domain.UserInvalidArgumentsError, Err: err}
	}

	userDb, err1 := s.UserDbService.Retrieve(ctx, "id", id)
	if err1 != nil {
		switch err1.Type {
		case domain.DbNotFoundError:
			return nil, &domain.UserError{Type: domain.UserNotFoundError, Err: err1}
		default:
			return nil, &domain.UserError{Type: domain.UserInternalError, Err: err1}
		}
	}

	if !authenticate(userDb, currentPasswd) {
		return nil, &domain.UserError{Type: domain.UserIncorrectPasswordError, Err: errors.New("incorrect password")}
	}

	userDb.Username = uwc.Username
	userDb.PasswordHash, _ = bcrypt.GenerateFromPassword([]byte(uwc.Password), s.BcryptCost)

	if _, err := s.UserDbService.Update(ctx, userDb); err != nil {
		switch err.Type {
		case domain.DbNotFoundError:
			return nil, &domain.UserError{Type: domain.UserNotFoundError, Err: err}
		default:
			return nil, &domain.UserError{Type: domain.UserInternalError, Err: err}
		}
	}

	return userDb.ToUser(), nil
}

func (s *UserService) UpdateClaims(ctx context.Context, c *domain.Claims, id string, uwc *domain.UserWriteClaims) (*domain.User, *domain.UserError) {
	if !c.Admin {
		return nil, &domain.UserError{Type: domain.UserUnauthorizedError, Err: errors.New("unauthorized")}
	}

	if err := s.Validator.Struct(uwc); err != nil {
		return nil, &domain.UserError{Type: domain.UserInvalidArgumentsError, Err: err}
	}

	userUpdate := domain.UserDb{Id: id, Admin: &uwc.Admin}
	userDb, err := s.UserDbService.Update(ctx, &userUpdate)
	if err != nil {
		switch err.Type {
		case domain.DbNotFoundError:
			return nil, &domain.UserError{Type: domain.UserNotFoundError, Err: err}
		default:
			return nil, &domain.UserError{Type: domain.UserInternalError, Err: err}
		}
	}

	return userDb.ToUser(), nil
}

func (s *UserService) Delete(ctx context.Context, c *domain.Claims, id string) *domain.UserError {
	if c.Id == "anonymous" || id != c.Id {
		return &domain.UserError{Type: domain.UserUnauthorizedError, Err: errors.New("unauthorized")}
	}

	if err := s.UserDbService.Delete(ctx, id); err != nil {
		return &domain.UserError{Type: domain.UserInternalError, Err: err}

	} else {
		return nil
	}
}

func (s *UserService) List(ctx context.Context, c *domain.Claims) ([]*domain.User, *domain.UserError) {
	if !c.Admin {
		return nil, &domain.UserError{Type: domain.UserUnauthorizedError, Err: errors.New("unauthorized")}
	}

	usersDb, err := s.UserDbService.RetrieveAll(ctx)

	if err != nil {
		return nil, &domain.UserError{Type: domain.UserInternalError, Err: err}
	}

	users := make([]*domain.User, len(usersDb))
	for i, udb := range usersDb {
		users[i] = udb.ToUser()
	}
	return users, nil
}

func (s *UserService) SendResetPassword(ctx context.Context, username string) *domain.UserError {
	if len(username) == 0 {
		return &domain.UserError{Type: domain.UserInvalidArgumentsError, Err: errors.New("missing username")}
	}

	userDb, err1 := s.UserDbService.Retrieve(ctx, "username", username)
	if err1 != nil {
		switch err1.Type {
		case domain.DbNotFoundError:
			return &domain.UserError{Type: domain.UserNotFoundError, Err: err1}
		default:
			return &domain.UserError{Type: domain.UserInternalError, Err: err1}
		}
	}

	if userDb == nil {
		return &domain.UserError{Type: domain.UserNotFoundError, Err: errors.New("user not found")}
	}

	claims := domain.Claims{Id: userDb.Id, Admin: false, SharedAccounts: []string{}, ResetPassword: true}
	resetToken, err2 := s.JwtService.GenerateResetPasswordToken(&claims)

	if err2 != nil {
		return &domain.UserError{Type: domain.UserInternalError, Err: err2}
	}

	clickUrl := s.FrontEndUrl
	resetTokenPath := fmt.Sprintf("/resetPassword/%s", resetToken)
	clickUrl.Path = path.Join(clickUrl.Path, resetTokenPath)

	emailMessage := fmt.Sprintf("Please reset your password by <a href=%s>clicking here</a>", clickUrl.String())
	email := domain.Email{
		ToAddress: userDb.Username,
		Subject:   "Reset password",
		Content:   emailMessage,
	}

	if err3 := s.EmailService.Send(ctx, &email); err3 != nil {
		return &domain.UserError{Type: domain.UserSendEmailError, Err: err3}
	}

	return nil
}

func (s *UserService) ResetPassword(ctx context.Context, token string, newPassword string) *domain.UserError {
	if len(token) == 0 || len(newPassword) == 0 {
		return &domain.UserError{Type: domain.UserInvalidArgumentsError, Err: errors.New("missing token or username")}
	}

	claims, err1 := s.JwtService.ParseToken(token)
	if err1 != nil {
		return &domain.UserError{Type: domain.UserInvalidArgumentsError, Err: errors.New("invalid token")}
	}

	if !claims.ResetPassword {
		return &domain.UserError{Type: domain.UserInvalidArgumentsError, Err: errors.New("invalid token")}
	}

	var newCreds = domain.UserWriteCredentials{Username: "", Password: newPassword}
	if err2 := s.Validator.Struct(newCreds); err2 != nil {
		return &domain.UserError{Type: domain.UserInvalidArgumentsError, Err: err2}
	}

	newHash, _ := bcrypt.GenerateFromPassword([]byte(newCreds.Password), s.BcryptCost)
	userUpdate := domain.UserDb{Id: claims.Id, PasswordHash: newHash}

	if _, err4 := s.UserDbService.Update(ctx, &userUpdate); err4 != nil {
		switch err4.Type {
		case domain.DbNotFoundError:
			return &domain.UserError{Type: domain.UserNotFoundError, Err: err4}
		default:
			return &domain.UserError{Type: domain.UserInternalError, Err: err4}
		}
	}

	return nil
}
