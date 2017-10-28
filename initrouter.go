package main

import (
	"github.com/WebGou/handlers"
	"github.com/WebGou/middleware"
	"github.com/gin-gonic/gin"
)

func initrouter(router *gin.Engine) {

	router.Static("/bootstrap", "./static/bootstrap-3.3.7")
	router.Static("/font-awesome", "./static/font-awesome-4.7.0")
	router.Static("/background", "./static/background")
	router.Static("/assets", "./static/assets")
	router.LoadHTMLGlob("templates/*")

	router.GET("/", handlers.IndexHandler)
	router.GET("/login", handlers.LoginHandler)
	router.POST("/login", handlers.LoginCust)
	router.GET("/signup", handlers.SignupG)
	router.POST("/signup", handlers.SignupP)

	router.GET("/logout", handlers.Logout)

	router.GET("/dashboard", handlers.Dashboard)

	router.GET("/google/auth", handlers.GoogleAuthHandler)
	router.GET("/facebook/auth", handlers.FaceBookAuthHandler)

	authorized := router.Group("/group")
	authorized.Use(middleware.AuthorizeRequest())
	{
		authorized.GET("/field", handlers.FieldHandler)
		authorized.GET("/loginusers", handlers.Loginusers)
	}

}
