package user

import (
	"database/sql"
	"github.com/elolpuer/Blog/pkg/models"

	"golang.org/x/crypto/bcrypt"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func Create(db *sql.DB, user *models.User) error {
	validate = validator.New()
	err := validate.Struct(user)
	if err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return err
		}
	}
	hash, err := HashPassword(user.Password)
	if err != nil {
		return err
	}
	_, err = db.Exec("INSERT INTO users (username, email, password) VALUES ($1, $2, $3)", user.Username, user.Email, hash)
	if err != nil {
		return err
	}
	return nil
}

func createSession(user *models.User)  {

}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}