package golangjwt

import (
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/jmilosze/wfrp-hammergen-go/internal/domain"
)

type HmacService struct {
	HmacSecret []byte
}

func NewHmacService(hmacSecret string) *HmacService {
	return &HmacService{
		HmacSecret: []byte(hmacSecret),
	}
}

func (jwtService *HmacService) GenerateToken(claims *domain.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": claims.UserId,
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
		return nil, &domain.InvalidTokenError{Err: fmt.Errorf("token did not pass validation")}
	}

	jwtClaims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, &domain.InvalidTokenError{Err: fmt.Errorf("error performin type Claims type assertion")}
	}

	var claims = domain.Claims{UserId: ""}
	claims.UserId = jwtClaims["sub"].(string)

	return &claims, nil
}
