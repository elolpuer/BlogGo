package models

type User struct {
	ID int
	Username string `validate:"required"`
	Email string `validate:"required,email"`
	Password string `validate:"required"`
}

type SessionUser struct {
	ID int
	Username string
}
