package main

import (
	"io"
	"os"
	"runtime"

	. "github.com/WebGou/baaplogger"
	"github.com/WebGou/handlers"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

var (
	cookiekey = "kOgyG9KZ5lIQJMGeZvm59ivKZxtfFP0R06q3+1F1gFaqSRIA/D4MFnURAGLcHhc0pzT90xi0Z6xfl5m0xSVWCg=="
)

func main() {
	Log.Informational("service start")
	Log.Informational("GOMAXPROCS: %d", runtime.GOMAXPROCS(-1))
	gin.DefaultWriter = io.MultiWriter(Ginlog, os.Stdout)
	gin.DefaultErrorWriter = gin.DefaultWriter
	router := gin.Default()

	store := sessions.NewCookieStore([]byte(cookiekey))
	store.Options(sessions.Options{
		Path:   "/",
		MaxAge: 86400 * 7,
	})
	router.Use(gin.Logger())

	router.Use(gin.Recovery())
	router.Use(sessions.Sessions("webgousession", store))

	initrouter(router)
	go handlers.InitSearch()

	router.Run("127.0.0.1:9090")
	Log.Informational("service end")
}
