package domain

import (
	"fmt"
	"strings"
)

type Claims struct {
	Id             string
	Admin          bool
	SharedAccounts []string
}

func (c *Claims) Set(u *UserRead) *Claims {
	c.Id = strings.Clone(u.Username)
	c.Admin = u.Admin

	c.SharedAccounts = make([]string, len(u.SharedAccounts))
	for i, s := range u.SharedAccounts {
		c.SharedAccounts[i] = strings.Clone(s)
	}

	return c
}

type JwtService interface {
	GenerateToken(claims *Claims) (string, error)
	ParseToken(token string) (*Claims, error)
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
