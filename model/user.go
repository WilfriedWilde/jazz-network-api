package model

type User struct {
	ID           int64  `json:"id"`
	Username     string `json:"username"`
	Email        string `json:"email"`
	PasswordHash string `json:"-"`
	AccessLevel  int    `json:"access_level"`
}

type UserInput struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserCreate struct {
	Username     string
	Email        string
	PasswordHash string
	AccessLevel  int
}
