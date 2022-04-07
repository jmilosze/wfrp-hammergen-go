package domain

type User struct {
	Id           string `json:"id"`
	Username     string `json:"username"`
	PasswordHash []byte `json:"password_hash"`
}

func (source *User) Copy() *User {
	userCopy := *source
	userCopy.PasswordHash = make([]byte, len(source.PasswordHash))
	copy(userCopy.PasswordHash, source.PasswordHash)
	return &userCopy
}

type UserService interface {
	FindUserById(id string) (*User, error)
	FindUserByName(username string) (*User, error)
	Authenticate(user User, password string) bool
}
