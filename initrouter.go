package main

import (
	"github.com/BaapAPI/handlers"
	"github.com/BaapAPI/middleware"
	"github.com/gin-gonic/gin"
)

func initrouter(router *gin.Engine) {

	router.Static("/bootstrap", "./static/bootstrap-3.3.7")
	router.Static("/font-awesome", "./static/font-awesome-4.7.0")
	router.Static("/background", "./static/background")
	router.LoadHTMLGlob("templates/*")

	router.GET("/", handlers.IndexHandler)
	router.GET("/login", handlers.LoginHandler)
	router.GET("/google/auth", handlers.GoogleAuthHandler)
	router.GET("/facebook/auth", handlers.FaceBookAuthHandler)

	authorized := router.Group("/battle")
	authorized.Use(middleware.AuthorizeRequest())
	{
		authorized.GET("/field", handlers.FieldHandler)
	}

}
