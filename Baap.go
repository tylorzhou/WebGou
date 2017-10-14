package main

import (
	"io"
	"os"

	. "github.com/BaapAPI/baaplogger"
	"github.com/BaapAPI/handlers"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

func main() {
	Log.Informational("service start")
	gin.DefaultWriter = io.MultiWriter(Ginlog, os.Stdout)
	gin.DefaultErrorWriter = gin.DefaultWriter
	router := gin.Default()

	store := sessions.NewCookieStore([]byte(handlers.RandToken(64)))
	store.Options(sessions.Options{
		Path:   "/",
		MaxAge: 86400 * 7,
	})
	router.Use(gin.Logger())

	router.Use(gin.Recovery())
	router.Use(sessions.Sessions("goquestsession", store))

	initrouter(router)

	router.Run("127.0.0.1:9090")
	Log.Informational("service end")
}
