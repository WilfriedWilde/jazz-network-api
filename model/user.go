package model

type User struct {
	ID           int64
	Username     string
	Email        string
	PasswordHash string
	AccessLevel  int
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
