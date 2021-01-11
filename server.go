package main

import (
	"github.com/elolpuer/Blog/pkg/controller"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	router := gin.Default()

	router.GET("/", controller.IndexGet)
	router.GET("/posts", controller.PostsGet)
	router.POST("/add/process", controller.AddPost)
	router.POST("/delete",controller.DeletePost)
	router.GET("/signup", controller.SignUpGet)
	router.GET("/signin", controller.SignInGet)
	router.POST("/signup/auth", controller.SignUpPost)
	router.POST("/signin/auth", controller.SignInPost)
	router.POST("/logout", controller.LogOutPost)

	err := router.Run(":5000")
	if err != nil {
		log.Fatal(err)
	}
}

