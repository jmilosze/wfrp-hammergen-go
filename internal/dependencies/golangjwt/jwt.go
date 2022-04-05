package golangjwt

import (
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/jmilosze/wfrp-hammergen-go/internal/domain"
	"time"
)

type HmacService struct {
	HmacSecret []byte
	ExpiryTime time.Duration
}

func NewHmacService(hmacSecret string, expiryTime time.Duration) *HmacService {
	return &HmacService{
		HmacSecret: []byte(hmacSecret),
		ExpiryTime: expiryTime,
	}
}

func (jwtService *HmacService) GenerateToken(claims *domain.Claims) (string, error) {
	currentTime := time.Now()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":      claims.UserId,
		"exp":      currentTime.Add(jwtService.ExpiryTime).Unix(),
		"orig_iat": currentTime.Unix(),
	})
	return token.SignedString(jwtService.HmacSecret)
}

func (jwtService *HmacService) ParseToken(tokenString string) (*domain.Claims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtService.HmacSecret, nil
	})

	if err != nil {
		return nil, &domain.InvalidTokenError{Err: err}
	}

	if !token.Valid {
		return nil, &domain.InvalidTokenError{Err: fmt.Errorf("invalid token")}
	}

	jwtClaims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, &domain.InvalidTokenError{Err: fmt.Errorf("error performin type Claims type assertion")}
	}

	if validErr := jwtClaims.Valid(); validErr != nil {
		return nil, &domain.InvalidTokenError{Err: fmt.Errorf("invalid token")}
	}

	var claims = domain.Claims{UserId: ""}
	claims.UserId = jwtClaims["sub"].(string)

	return &claims, nil
}
