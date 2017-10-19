package handlers

import (
	"crypto/rand"
	"crypto/sha1"
	"crypto/subtle"
	"encoding/base64"

	"github.com/gorilla/sessions"
)

const (
	rememberMeValidatorLen = 32
)

var (
	rememberMeStore sessions.Store
)

//LoginCookie Implementation of https://paragonie.com/blog/2015/04/secure-authentication-php-with-long-term-persistence#title.2
type LoginCookie struct {
	Selector   string
	Validator  string
	CookieName string
}

// Compute the sha-256
func (l *LoginCookie) validatorHash() string {
	hasher := sha1.New()
	hasher.Write([]byte(l.Validator))
	return base64.URLEncoding.EncodeToString(hasher.Sum(nil))
}

//GenerateValidator Fills Validator with a new value and return it's sha-256 hash.
// This hash, together with Selector need to be saved for later comparison.
func (l *LoginCookie) GenerateValidator() (string, error) {
	b := make([]byte, rememberMeValidatorLen)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	l.Validator = base64.URLEncoding.EncodeToString(b)

	return l.validatorHash(), nil
}

// Check if Validator value is valid...
func (l *LoginCookie) Check(value string) bool {
	// Prevents timing atack
	return subtle.ConstantTimeCompare([]byte(l.validatorHash()), []byte(value)) == 1
}
