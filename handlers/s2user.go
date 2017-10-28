package handlers

import (
	"github.com/WebGou/baapDB"
	. "github.com/WebGou/baaplogger"
	"github.com/gin-gonic/contrib/sessions"
)

//GetUser get user name by session
func GetUser(s sessions.Session) (user string) {
	var rm Rememberme
	_, userkey, logintype, err := rm.Check(s)
	if err != nil {
		Log.Error("GetUser failed %s", err.Error())
	}

	switch logintype {
	case llogin:
		user, _ = baapDB.CheckLocalUser(userkey)
	case glogin:
		user, _ = baapDB.CheckGUser(userkey)
	case flogin:
		user, _ = baapDB.CheckFUser(userkey)
	default:
		Log.Error("GetUser failed with empty user")
	}

	return

}
