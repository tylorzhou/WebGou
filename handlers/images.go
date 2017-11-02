package handlers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	. "github.com/WebGou/baaplogger"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

//ImageuploadG image upload for get function
func ImageuploadG(c *gin.Context) {
	c.HTML(http.StatusOK, "uploadimg.tmpl", gin.H{})
}

//ImageuploadP image upload for post function
func ImageuploadP(c *gin.Context) {
	s := sessions.Default(c)
	_, logintype, err := GetUserkey(s)
	if err != nil {
		Log.Error(err.Error())
		return
	}

	user, uid := GetUser(s)
	usertype := UserTypePath(logintype)

	t := time.Now()
	timestamp := t.Format("2006102150405")
	form, _ := c.MultipartForm()
	images := form.File["images[]"]
	suid := strconv.Itoa(uid)
	path := filepath.Join(exePath, "images", usertype, suid, timestamp)

	err = os.MkdirAll(path, os.ModePerm)
	if err != nil {
		Log.Error(err.Error())
		return
	}
	var i int
	for i, file := range images {
		if i > 5 {
			break
		}

		name := "image" + strconv.Itoa(i)
		err := c.SaveUploadedFile(file, filepath.Join(path, name))
		if err != nil {
			Log.Error("upload image failed user: %s logintype: %d ,uid: %d,upload files: %s => %s, err %s", user, logintype,
				uid, file.Filename, err.Error())
		}
		Log.Debug("user: %s logintype: %d ,uid: %d,upload files: %s => %s", user, logintype, uid, file.Filename, name)
	}
	c.String(http.StatusOK, fmt.Sprintf("%d files uploaded!", i+1))
}

//UserTypePath get path for different login type
func UserTypePath(loginttype int) string {
	switch loginttype {
	case glogin:
		return "guser"
	case flogin:
		return "fuser"
	case llogin:
		return "luser"
	default:
		return ""
	}
}
