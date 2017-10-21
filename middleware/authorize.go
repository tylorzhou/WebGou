package middleware

import (
	"net/http"

	. "github.com/WebGou/baaplogger"
	"github.com/WebGou/handlers"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

var (
	rm handlers.Rememberme
)

func init() {
	rm = handlers.Rememberme{MaxAge: handlers.DefrmDur}
}

// AuthorizeRequest is used to authorize a request for a certain end-point group.
func AuthorizeRequest() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)

		_, _, err := rm.Check(session)

		if err != nil {
			if err.Error() != "session expired" {
				Log.Error("AuthorizeRequest faild: %s", err.Error())
			}

			c.HTML(http.StatusUnauthorized, "error.tmpl", gin.H{"message": "Please login."})
			c.Abort()
		}
		c.Next()
	}
}
