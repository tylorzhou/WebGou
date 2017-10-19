package handlers

import (
	"time"

	"github.com/WebGou/baapDB"
	"github.com/gin-gonic/contrib/sessions"
)

//Rememberme for remember me feaure
type Rememberme struct {
	MaxAge time.Duration
}

const (
	rememberme = "rememberme"
)

//Check if there is cookie
func (c *Rememberme) Check(s sessions.Session) (selector, user string, err error) {
	var now = time.Now()

	l, ok := s.Get(rememberme).(LoginCookie)
	if !ok {
		err = ErrorBadRMSessionFormat
		return
	}

	var hash string
	user, hash, expires, err := c.get(l.Selector)
	if err != nil {
		return
	}

	if expires.Before(now) {
		err = ErrorSessionExpired
		return
	}

	if !l.Check(hash) {
		err = ErrorBadRequest
	}

	selector = l.Selector

	return
}

//SetCookie set
func (c *Rememberme) SetCookie(s sessions.Session, user string, logintype int) (err error) {
	l := LoginCookie{
		CookieName: rememberme,
	}

	hash, err := l.GenerateValidator()
	if err != nil {
		return
	}

	// First save to the database
	l.Selector, err = c.insert(user, hash, time.Now().Add(c.MaxAge), logintype)
	if err != nil {
		return
	}

	// Then save to the cookie
	s.Set(rememberme, l)
	s.Save()
	return
}

//UpdateCookie update
func (c *Rememberme) UpdateCookie(s sessions.Session, selector, user string) (err error) {
	l := LoginCookie{
		Selector:   selector,
		CookieName: rememberme,
	}

	hash, err := l.GenerateValidator()
	if err != nil {
		return
	}

	// First save to the database
	err = c.update(selector, user, hash, time.Now().Add(c.MaxAge))
	if err != nil {
		return
	}

	// Then save to the cookie
	s.Set(rememberme, l)
	s.Save()
	return
}

func (c *Rememberme) get(selector string) (user string, hash string, expiration time.Time, err error) {
	return baapDB.Rmmeget(selector)
}
func (c *Rememberme) insert(user, hash string, expiration time.Time, logintype int) (selector string, err error) {
	return baapDB.Rmmeinsert(user, hash, expiration, logintype)
}
func (c *Rememberme) update(selector, user, hash string, expiration time.Time) (err error) {
	return baapDB.Rmmeupdate(selector, user, hash, expiration)
}
