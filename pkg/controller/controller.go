package controller

import (
	"database/sql"
	"github.com/elolpuer/Blog/pkg/auth"
	db2 "github.com/elolpuer/Blog/pkg/db"
	"github.com/elolpuer/Blog/pkg/models"
	"github.com/elolpuer/Blog/pkg/post"
	tml2 "github.com/elolpuer/Blog/pkg/tml"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
)

var tml = tml2.GetTemplates()
var db *sql.DB

func init(){
	var err error
	db, err = db2.ConnectionToDB()
	if err != nil {
		log.Fatal(err)
	}
}

func IndexGet(ctx *gin.Context) {
	var user bool
	if isNew, err := auth.SessionIsNew(ctx.Request, "session-name"); isNew != true || err != nil {
		user = true
	}
	tml.ExecuteTemplate(ctx.Writer, "index.gohtml", struct {
		Title string
		H1 string
		User bool
	}{
		Title: "Index Page",
		H1: "Index Page",
		User: user,
	})
}

func PostsGet(ctx *gin.Context) {
	var userYet bool
	if isNew, err := auth.SessionIsNew(ctx.Request, "session-name"); isNew != true || err != nil {
		userYet = true
	}
	user, err := auth.GetSessionUser(ctx.Request,"session-name" )
	posts, err := post.Posts(db, user.ID)
	if err != nil {
		log.Fatal(err)
	}
	tml.ExecuteTemplate(ctx.Writer, "posts.gohtml", struct {
		Title string
		H1 string
		User bool
		Posts []*models.Post
	}{
		Title: "Posts",
		H1: "Add",
		User: userYet,
		Posts: posts,
	})
}


func SignUpGet(ctx *gin.Context) {
	if isNew, err := auth.SessionIsNew(ctx.Request, "session-name"); isNew != true || err != nil {
		ctx.Redirect(http.StatusSeeOther, "/posts")
	}
	user := false
	tml.ExecuteTemplate(ctx.Writer, "signup.gohtml",struct {
		Title string
		H1 string
		User bool
	}{
		Title: "Sign Up",
		H1: "Sign Up",
		User: user,
	})
}


func SignInGet(ctx *gin.Context) {
	if isNew, err := auth.SessionIsNew(ctx.Request, "session-name"); isNew != true || err != nil {
		ctx.Redirect(http.StatusSeeOther, "/posts")
	}
	user := false
	tml.ExecuteTemplate(ctx.Writer, "signin.gohtml",struct {
		Title string
		H1 string
		User bool
	}{
		Title: "Sign In",
		H1: "Sign In",
		User: user,
	})
}

func AddPost(ctx *gin.Context) {
	user, err := auth.GetSessionUser(ctx.Request,"session-name" )
	text := ctx.PostForm("body")
	err = post.Add(db, user.ID, text)
	if err != nil {
		log.Fatal(err)
	}
	ctx.Redirect(http.StatusSeeOther, "/posts")
}

func DeletePost(ctx *gin.Context) {
	session, err := auth.GetSessionStore(ctx.Request,  "session-name")
	if err != nil {
		http.Error(ctx.Writer, err.Error(), http.StatusInternalServerError)
	}
	userID := session.Values["userID"].(int)
	id := ctx.Query("id")
	err = post.DeletePost(db, id, userID)
	if err != nil {
		log.Fatal(err)
	}
	ctx.Redirect(http.StatusSeeOther, "/posts")
}

func SignUpPost(ctx *gin.Context) {
	if isNew, err := auth.SessionIsNew(ctx.Request, "session-name"); isNew != true || err != nil {
		ctx.Redirect(http.StatusSeeOther, "/")
	}
	var NewUser = new(models.User)
	NewUser.Username = ctx.PostForm("username")
	NewUser.Email = ctx.PostForm("email")
	NewUser.Password = ctx.PostForm("password")
	err := auth.SignUp(db, NewUser)
	if err != nil {
		log.Fatal(err)
	}
	ctx.Redirect(http.StatusSeeOther, "/")
}

func SignInPost(ctx *gin.Context) {
	var signUser = new(models.User)
	signUser.Email = ctx.PostForm("email")
	signUser.Password = ctx.PostForm("password")
	sessionUser, err := auth.SignIn(db, signUser)
	if err == sql.ErrNoRows || err == bcrypt.ErrMismatchedHashAndPassword{
		http.Error(ctx.Writer, "Invalid Data", 400)
		return
	}
	err = auth.CreateSessionUser(ctx.Writer, ctx.Request, sessionUser, "session-name")
	if err != nil {
		http.Error(ctx.Writer, err.Error(), http.StatusInternalServerError)
		return
	}
	ctx.Redirect(http.StatusSeeOther, "/posts")
}

func LogOutPost(ctx *gin.Context) {
	err := auth.Logout(ctx.Writer, ctx.Request, "session-name")
	if err != nil {
		http.Error(ctx.Writer, err.Error(), http.StatusInternalServerError)
		return
	}
	ctx.Redirect(http.StatusSeeOther, "/")
}