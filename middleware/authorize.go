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
	rm = handlers.Rememberme{}
}

// AuthorizeRequest is used to authorize a request for a certain end-point group.
func AuthorizeRequest() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)

		selector, user, logintype, err := rm.Check(session)

		if err != nil {
			if err.Error() != "session expired" {
				Log.Error("AuthorizeRequest faild: %s", err.Error())
			}

			c.HTML(http.StatusUnauthorized, "errLogin.tmpl", gin.H{"message": "Sorry. Please login first."})
			c.Abort()
		} else {
			rm.UpdateCookie(session, selector, user, logintype)
		}

		c.Next()
	}
}
