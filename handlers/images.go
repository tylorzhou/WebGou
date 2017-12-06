package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/WebGou/baapDB"
	. "github.com/WebGou/baaplogger"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

type imgdata struct {
	Keywords    string
	Description string
}

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

	user, uid, _ := GetUser(s)
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

	uploaded := 0
	for i, file := range images {
		if i > 5 {
			break
		}

		name := "image" + strconv.Itoa(i) + filepath.Ext(file.Filename)
		err := c.SaveUploadedFile(file, filepath.Join(path, name))
		if err != nil {
			Log.Error("upload image failed user: %s logintype: %d ,uid: %d,upload files: %s => %s, err %s", user, logintype,
				uid, file.Filename, err.Error())
		}
		uploaded = i + 1
		Log.Debug("user: %s logintype: %d ,uid: %d,upload files: %s => %s", user, logintype, uid, file.Filename, name)
	}

	if uploaded > 0 {
		tableName := baapDB.GImgTblName(logintype, uid)

		img := baapDB.Imageinfo{
			Logintype:   logintype,
			ID:          uid,
			Imageurl:    filepath.Join(usertype, suid, timestamp),
			Description: "",
			Created:     t.Format("2006-01-02 15:04:05"),
		}

		baapDB.InsertImage(tableName, img)
	}

	if uploaded > 0 {
		var idata imgdata
		idata.Keywords = c.PostForm("Keywords")
		idata.Description = c.PostForm("Description")
		b, err := json.Marshal(idata)
		if err != nil {
			Log.Error("json marshal failed for %d, %v", uid, err)
			return
		}
		ioutil.WriteFile(filepath.Join(path, "config.json"), b, 0644)
	}

	c.HTML(http.StatusOK, "UploadRes.tmpl", gin.H{"SCInfo": fmt.Sprintf("%d files uploaded!", uploaded)})
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
