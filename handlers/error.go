package handlers

import (
	"errors"
)

var (
	//ErrorBadRequest for bad requests
	ErrorBadRequest = errors.New("bad request")
	//ErrorSessionExpired indicate that session expired
	ErrorSessionExpired = errors.New("session expired")
	//ErrorBadRMSessionFormat bad remember cookie's format
	ErrorBadRMSessionFormat = errors.New("Bad remember me cookie format")
)
