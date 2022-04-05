package domain

type User struct {
	Id           string `json:"id"`
	Username     string `json:"username"`
	PasswordHash string `json:"password_hash"`
}

type UserService interface {
	FindUserById(id string) (*User, error)
	FindUserByName(username string) (*User, error)
}
