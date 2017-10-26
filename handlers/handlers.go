package handlers

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/WebGou/baapDB"
	. "github.com/WebGou/baaplogger"
	"github.com/WebGou/structs"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
	"golang.org/x/oauth2/google"
)

var (
	gooleconf    *oauth2.Config
	facebookconf *oauth2.Config
	procred      ProvideCre
	rm           = Rememberme{MaxAge: 86400 * 7}
)

const (
	glogin = iota // google login
	flogin        //facebook login
	llogin        //customer local login
)

// Credentials which stores id and key.
type Credentials struct {
	Cid     string `json:"client_id"`
	Csecret string `json:"client_secret"`
}

//ProvideCre which store different id and key for provider
type ProvideCre map[string]Credentials

// RandToken generates a random @l length token.
func RandToken(l int) string {
	b := make([]byte, l)
	rand.Read(b)
	return base64.StdEncoding.EncodeToString(b)
}

func getLoginURL(state, provider string) string {
	if provider == "facebook" {
		return facebookconf.AuthCodeURL(state)
	} else if provider == "google" {
		return gooleconf.AuthCodeURL(state)
	} else {
		return ""
	}

}

func init() {
	file, err := ioutil.ReadFile("./client_secret.json")
	if err != nil {
		Log.Error("File error: %v\n", err)
		os.Exit(1)
	}
	json.Unmarshal(file, &procred)

	gooleconf = &oauth2.Config{
		ClientID:     procred["google"].Cid,
		ClientSecret: procred["google"].Csecret,
		RedirectURL:  "http://127.0.0.1:9090/google/auth",
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email", // You have to select your own scope from here -> https://developers.google.com/identity/protocols/googlescopes#google_sign-in
		},
		Endpoint: google.Endpoint,
	}

	facebookconf = &oauth2.Config{
		ClientID:     procred["facebook"].Cid,
		ClientSecret: procred["facebook"].Csecret,
		RedirectURL:  "http://localhost:9090/facebook/auth",
		Scopes:       []string{"public_profile"},
		Endpoint:     facebook.Endpoint,
	}

}

// IndexHandler handels /.
func IndexHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "index.tmpl", gin.H{})
}

// GoogleAuthHandler handles authentication of a user and initiates a session.
func GoogleAuthHandler(c *gin.Context) {
	// Handle the exchange code to initiate a transport.
	session := sessions.Default(c)
	retrievedState := session.Get("state")
	queryState := c.Request.URL.Query().Get("state")
	if retrievedState != queryState {
		Log.Error("Invalid session state: retrieved: %s; Param: %s", retrievedState, queryState)
		c.HTML(http.StatusUnauthorized, "error.tmpl", gin.H{"message": "Invalid session state."})
		return
	}
	code := c.Request.URL.Query().Get("code")
	tok, err := gooleconf.Exchange(oauth2.NoContext, code)
	if err != nil {
		Log.Error(err.Error())
		c.HTML(http.StatusBadRequest, "error.tmpl", gin.H{"message": "Login failed. Please try again."})
		return
	}

	client := gooleconf.Client(oauth2.NoContext, tok)
	userinfo, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		Log.Error(err.Error())
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	defer userinfo.Body.Close()
	data, _ := ioutil.ReadAll(userinfo.Body)
	u := structs.GoogleUser{}
	if err = json.Unmarshal(data, &u); err != nil {
		Log.Error(err.Error())
		c.HTML(http.StatusBadRequest, "error.tmpl", gin.H{"message": "Error marshalling response. Please try agian."})
		return
	}
	session.Set("user-id", u.Email)
	err = session.Save()
	if err != nil {
		Log.Error(err.Error())
		c.HTML(http.StatusBadRequest, "error.tmpl", gin.H{"message": "Error while saving session. Please try again."})
		return
	}

	Log.Debug("google user login: " + u.Email)

	exist, _ := baapDB.IsUserExist("google", u.Email)

	if !exist {
		baapDB.AddGoogleUser(u.Email, u.Name)
	}

	session.Options(sessions.Options{
		Path:   "/",
		MaxAge: 0,
	}) // thirdpart provider do need persistent cookie

	rm.SetCookie(session, u.Email, glogin)

	c.HTML(http.StatusOK, "userls.tmpl", gin.H{"GUsers": baapDB.GetGUser(), "FUsers": baapDB.GetFUser()})
}

//Loginusers handle login users
func Loginusers(c *gin.Context) {
	c.HTML(http.StatusOK, "userls.tmpl", gin.H{"GUsers": baapDB.GetGUser(), "FUsers": baapDB.GetFUser()})
}

// FaceBookAuthHandler handles authentication of a user and initiates a session.
func FaceBookAuthHandler(c *gin.Context) {
	retrievedState := c.Request.FormValue("state")
	queryState := c.Request.URL.Query().Get("state")
	if retrievedState != queryState {
		Log.Error("Invalid session state: retrieved: %s; Param: %s", retrievedState, queryState)
		c.HTML(http.StatusUnauthorized, "error.tmpl", gin.H{"message": "Invalid session state."})
		return
	}
	code := c.Request.URL.Query().Get("code")
	tok, err := facebookconf.Exchange(oauth2.NoContext, code)
	if err != nil {
		Log.Error(err.Error())
		c.HTML(http.StatusBadRequest, "error.tmpl", gin.H{"message": "Login failed. Please try again."})
		return
	}

	client := facebookconf.Client(oauth2.NoContext, tok)
	userinfo, err := client.Get("https://graph.facebook.com/me?")
	if err != nil {
		Log.Error(err.Error())
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	defer userinfo.Body.Close()
	data, _ := ioutil.ReadAll(userinfo.Body)
	u := structs.FaceBookUser{}
	if err = json.Unmarshal(data, &u); err != nil {
		Log.Error(err.Error())
		c.HTML(http.StatusBadRequest, "error.tmpl", gin.H{"message": "Error marshalling response. Please try agian."})
		return
	}
	Log.Debug("facebook user login: " + u.Name + " id: " + u.ID)

	exist, _ := baapDB.IsUserExist("facebook", u.ID)

	if !exist {
		baapDB.AddFacebookUser(u.ID, u.Name)
	}

	session := sessions.Default(c)

	session.Options(sessions.Options{
		Path:   "/",
		MaxAge: 0,
	}) // thirdpart provider do need persistent cookie
	rm.SetCookie(session, u.ID, flogin)
	c.HTML(http.StatusOK, "userls.tmpl", gin.H{"GUsers": baapDB.GetGUser(), "FUsers": baapDB.GetFUser()})
}

// LoginHandler handles the login procedure.
func LoginHandler(c *gin.Context) {
	state := RandToken(32)
	session := sessions.Default(c)
	session.Set("state", state)
	Log.Informational("Stored session: %v\n", state)
	glink := getLoginURL(state, "google")
	flink := getLoginURL(state, "facebook")
	loginerr := session.Flashes("loginerr")[0]
	session.Save()
	c.HTML(http.StatusOK, "login.tmpl", gin.H{"glink": glink, "flink": flink, "loginerr": loginerr})
}

// FieldHandler is a rudementary handler for logged in users.
func FieldHandler(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("user-id")
	c.HTML(http.StatusOK, "field.tmpl", gin.H{"user": userID})
}

//LoginCust handle custom login
func LoginCust(c *gin.Context) {
	name := c.PostForm("InputEmail")
	pw := c.PostForm("InputPassword")
	keeplogin := c.PostForm("keeplogin")

	Log.Informational("name: %s pw:%s keeplogin:%s\n", name, pw, keeplogin)

	err := baapDB.LoginLocalUser(name, pw)
	s := sessions.Default(c)
	if err != nil {
		Log.Informational("login name %s pw %s err %s", name, pw, err.Error())

		s.AddFlash(err.Error(), "loginerr")
		s.Save()
		c.Redirect(http.StatusMovedPermanently, "/login")
	} else {
		rm.SetCookie(s, name, llogin)
		c.Redirect(http.StatusMovedPermanently, "/group/loginusers")
	}

}

//SignupG to get sign up page
func SignupG(c *gin.Context) {
	c.HTML(http.StatusOK, "signup.tmpl", gin.H{
		"name":   "",
		"em":     "",
		"pw":     "",
		"pwc":    "",
		"namerr": "",
		"emerr":  "",
		"pwerr":  "",
		"pwcerr": "",
	})
}

//SignupP to signup
func SignupP(c *gin.Context) {

	name := c.PostForm("username")
	em := c.PostForm("email")
	pw := c.PostForm("password")
	pwc := c.PostForm("password_confirm")

	var namerr, emerr, pwerr, pwcerr string
	if len(name) == 0 {
		namerr = "name cannot be null"
	}

	if em == "" {
		emerr = "email cannot be null"
	}

	if pw == "" {
		pwerr = "password cannot be null"
	}

	if pwc == "" {
		pwcerr = "confirm password cannot be null"
	}

	if pw != "" && pwc != "" && pw != pwc {
		pwcerr = "confirm password is not the same as password"
	}

	if emerr == "" && pwerr == "" && pwcerr == "" && namerr == "" {
		_, b := baapDB.CheckLocalUser(em)
		if b {
			emerr = "this email already registered"
		} else {
			err := baapDB.AddLocalUser(name, em, pw)
			if err != nil {
				Log.Error("Addlocauser failed, name: %s em: %s", name, em)
			} else {
				c.Redirect(http.StatusMovedPermanently, "/login")
				return
			}
		}
	}

	c.HTML(http.StatusOK, "signup.tmpl", gin.H{
		"name":   name,
		"em":     em,
		"pw":     pw,
		"pwc":    pwc,
		"namerr": namerr,
		"emerr":  emerr,
		"pwerr":  pwerr,
		"pwcerr": pwcerr,
	})

}
