package handlers

import (
	"github.com/WebGou/baapDB"
	. "github.com/WebGou/baaplogger"
	"github.com/gin-gonic/contrib/sessions"
)

//GetUser get user name by session
func GetUser(s sessions.Session) (user string, uid, logtype int) {
	var rm Rememberme
	_, userkey, logintype, err := rm.Check(s)
	if err != nil {
		Log.Error("GetUser failed %s", err.Error())
		return "", 0, -1
	}

	logtype = logintype
	switch logintype {
	case llogin:
		user, uid, _ = baapDB.CheckLocalUser(userkey)
	case glogin:
		user, uid, _ = baapDB.CheckGUser(userkey)
	case flogin:
		user, uid, _ = baapDB.CheckFUser(userkey)
	default:
		Log.Error("GetUser failed with empty user")
	}

	return

}

//GetUserkey get user key and logintype by session
func GetUserkey(s sessions.Session) (userkey string, logintype int, err error) {
	var rm Rememberme
	_, userkey, logintype, err = rm.Check(s)
	if err != nil {
		Log.Error("GetUser failed %s", err.Error())
	}

	return

}
