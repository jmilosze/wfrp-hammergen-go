package domain

import "fmt"

type Claims struct {
	UserId string
}

type JwtService interface {
	GenerateToken(claims *Claims) (string, error)
	ParseToken(tokenString string) (*Claims, error)
}

type InvalidTokenError struct {
	Err error
}

func (e *InvalidTokenError) Unwrap() error {
	return e.Err
}

func (e *InvalidTokenError) Error() string {
	return fmt.Sprintf("invalid token, %s", e.Err)
}
