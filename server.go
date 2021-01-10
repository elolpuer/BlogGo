package main

import (
	"database/sql"
	"github.com/elolpuer/Blog/cfg"
	"github.com/elolpuer/Blog/pkg/index"
	"github.com/elolpuer/Blog/pkg/models"
	"github.com/elolpuer/Blog/pkg/post"
	"github.com/elolpuer/Blog/pkg/user"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	_ "github.com/lib/pq"
	"html/template"
	"log"
	"net/http"
)


var (
	db *sql.DB
	tml *template.Template
)

var store = sessions.NewCookieStore([]byte(cfg.GetSessionKey()))

func main() {
	var err error
	tml = template.Must(template.ParseGlob("templates/*.gohtml"))
	router := gin.Default()
	db, err = sql.Open("postgres", cfg.GetPostgres())
	if err != nil {
		log.Fatal(err)
	}
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}
	router.GET("/", indexPage)
	router.GET("/posts", postsPage)
	router.POST("/add/process", addPost)
	router.POST("/delete",deletePost)
	router.GET("/signup", signup)
	router.POST("/signup/auth", createUser)

	err = router.Run(":5000")
	if err != nil {
		log.Fatal(err)
	}
}

func indexPage(ctx *gin.Context) {
	resp := index.Page()
	tml.ExecuteTemplate(ctx.Writer, "index.gohtml", resp)
}

func postsPage(ctx *gin.Context) {
	posts, err := post.Posts(db)
	if err != nil {
		log.Fatal(err)
	}
	tml.ExecuteTemplate(ctx.Writer, "posts.gohtml", posts)
}

func signup(ctx *gin.Context) {
	tml.ExecuteTemplate(ctx.Writer, "signup.gohtml","")
}

func addPost(ctx *gin.Context) {
	text := ctx.PostForm("body")
	err := post.Add(db, text)
	if err != nil {
		log.Fatal(err)
	}
	ctx.Redirect(http.StatusSeeOther, "/posts")
}

func deletePost(ctx *gin.Context) {
	id := ctx.Query("id")
	err := post.DeletePost(db, id)
	if err != nil {
		log.Fatal(err)
	}
	ctx.Redirect(http.StatusSeeOther, "/posts")
}

func createUser(ctx *gin.Context) {
	var NewUser = new(models.User)
	NewUser.Username = ctx.PostForm("username")
	NewUser.Email = ctx.PostForm("email")
	NewUser.Password = ctx.PostForm("password")
	err := user.Create(db, NewUser)
	if err != nil {
		log.Fatal(err)
	}
	ctx.Redirect(http.StatusSeeOther, "/")
}