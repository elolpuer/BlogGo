package auth

import (
	"database/sql"
	"github.com/elolpuer/Blog/cfg"
	"github.com/elolpuer/Blog/pkg/models"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

var validate *validator.Validate
var store = sessions.NewCookieStore([]byte(cfg.GetSessionKey()))


func SignUp(db *sql.DB, user *models.User) error {
	validate = validator.New()
	err := validate.Struct(user)
	if err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return err
		}
	}
	hash, err := hashPassword(user.Password)
	if err != nil {
		return err
	}
	_, err = db.Exec("INSERT INTO users (username, email, password) VALUES ($1, $2, $3)", user.Username, user.Email, hash)
	if err != nil {
		return err
	}
	return nil
}

func SignIn(db *sql.DB, user *models.User) (*models.SessionUser, error) {
	var signUser = new(models.User)
	row := db.QueryRow("SELECT * from users WHERE email=$1", user.Email)
	if err := row.Scan(&signUser.ID, &signUser.Username, &signUser.Email, &signUser.Password); err == sql.ErrNoRows {
		return nil, err
	}
	match := checkPasswordHash(user.Password, signUser.Password)
	if match != true {
		return nil, bcrypt.ErrMismatchedHashAndPassword
	}
	var sessionUser = new(models.SessionUser)
	sessionUser.ID = signUser.ID
	sessionUser.Username = signUser.Username
	return sessionUser, nil
}

func GetSessionStore(ctxReq *http.Request, name string) (*sessions.Session, error){
	session, err := store.Get(ctxReq, name)
	if err != nil {
		return nil,err
	}
	return session, nil
}

func CreateSessionUser(w http.ResponseWriter, ctxReq *http.Request,sUser *models.SessionUser, name string) error {
	s, err := GetSessionStore(ctxReq, name)
	if err != nil {
		return err
	}
	if s.IsNew != true {
		return sessions.MultiError{}
	}
	s.Values["userID"] = sUser.ID
	s.Values["username"] = sUser.Username
	err = s.Save(ctxReq, w)
	return nil
}

func GetSessionUser(ctxReq *http.Request, name string) (*models.SessionUser, error) {
	s, err := GetSessionStore(ctxReq, name)
	if err != nil {
		return nil, err
	}
	userID := s.Values["userID"].(int)
	username := s.Values["username"].(string)
	return &models.SessionUser{
		ID : userID,
		Username: username,
	}, nil
}

func SessionIsNew(ctxReq *http.Request, name string) (bool, error) {
	s, err := store.Get(ctxReq, name)
	if err != nil {
		return false, err
	}
	if s.IsNew != true {
		return false, nil
	}
	return true,nil
}

func Logout(w http.ResponseWriter, ctxReq *http.Request, name string) error {
	s, err := store.Get(ctxReq, name)
	if err != nil {
		return err
	}
	s.Options.MaxAge = -1
	err = s.Save(ctxReq, w)
	if err != nil {
		return err
	}
	return nil
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}