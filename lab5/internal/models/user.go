package models

// User - зарегистрированный пользователь.
type User struct {
	ID           int
	Email        string
	PasswordHash string // хэш пароля (сам пароль нигде не хранится)
}
