package baapDB

import (
	"database/sql"
	"encoding/json"
	"errors"
	"io/ioutil"
	"path/filepath"
	"time"

	"github.com/WebGou/baaplogger"
	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

// DBconfig which store db setting.
type DBconfig struct {
	User     string `json:"user"`
	Password string `json:"password"`
	IP       string `json:"ip"`
	Port     string `json:"port"`
	DBname   string `json:"dbname"`
}

var (
	db    *sql.DB
	dblog *baaplogger.Baaplogger
)

const mysqlDateFormat = "2006-01-02 15:04:05"

func init() {

	dir, err := filepath.Abs("log")
	if err != nil {
		panic(err)
	}

	// this for baap API log
	dblog = &baaplogger.Baaplogger{
		Level: baaplogger.LevelDebug,
		Log: &baaplogger.Logger{
			Filename:   filepath.Join(dir, "db.log"),
			MaxSize:    500, // megabytes
			MaxBackups: 6,
			MaxAge:     28, // days
		},
	}

	dbcfg, err := initDBconf()
	if err != nil {
		dblog.Critical(err.Error())
	}
	//db, err := sql.Open("mysql", "root:Altigen1234@tcp(<HOST>:<port>)/<dbname>"
	con := dbcfg.User + ":" + dbcfg.Password + "@tcp" + "(" + dbcfg.IP + ":" + dbcfg.Port + ")" +
		"/" + dbcfg.DBname
	db, err = sql.Open("mysql", con)
	if err != nil {
		dblog.Critical("open database failed: %s", err.Error())
	}
	err = db.Ping()
	if err != nil {
		dblog.Error("connect to db failed: %s", err.Error())
	}
}

func initDBconf() (*DBconfig, error) {

	file, err := ioutil.ReadFile("./DBsetting.json")
	if err != nil {
		dblog.Critical("ReadFile DBsetting.json failed: %s", err.Error())
		return nil, errors.New("initDBconf failed")
	}
	var dbconf = &DBconfig{}
	json.Unmarshal(file, &dbconf)
	return dbconf, nil
}

//AddGoogleUser add user from google
func AddGoogleUser(email, name string) error {
	addError := errors.New("AddGoogleUser failed")

	stmtIns, err := db.Prepare("INSERT INTO GoogleUser SET email=?,name=?,created=?")
	if err != nil {
		dblog.Error("db prepare failed: %s", err.Error())
	}
	defer stmtIns.Close() // Close the statement when we leave main() / the program terminates

	now := time.Now()

	result, err := stmtIns.Exec(email, name, now.Format(mysqlDateFormat))
	if err != nil {
		dblog.Error("db Exec failed: %s", err.Error())
		return addError
	}

	insertid, err := result.LastInsertId()
	if err != nil {
		dblog.Error("db LastInsertId failed: %s", err.Error())
		return addError
	}

	dblog.Debug("add google user, email: %s name: %s, insertid:%d", email, name, insertid)
	return nil
}

//AddFacebookUser add user from facebook
func AddFacebookUser(id, name string) error {
	addError := errors.New("AddFacebookUser failed")

	stmtIns, err := db.Prepare("INSERT INTO FacebookUser SET id=?,name=?,created=?") // ? = placeholder
	if err != nil {
		dblog.Error("db prepare failed: %s", err.Error())
		return addError
	}
	defer stmtIns.Close() // Close the statement when we leave main() / the program terminates

	now := time.Now()
	result, err := stmtIns.Exec(id, name, now.Format(mysqlDateFormat))
	if err != nil {
		dblog.Error("db Exec failed: %s", err.Error())
		return addError
	}

	insertid, err := result.LastInsertId()
	if err != nil {
		dblog.Error("db LastInsertId failed: %s", err.Error())
		return addError
	}

	dblog.Debug("add facebook user, id: %s name: %s, insertid: %d", id, name, insertid)
	return nil
}

//IsUserExist check if user exist
func IsUserExist(provider, key string) (bool, string) {

	var rows *sql.Rows
	var err error
	if provider == "facebook" {
		rows, err = db.Query("SELECT name from FacebookUser where id = ?", key)

	} else if provider == "google" {
		rows, err = db.Query("SELECT name FROM GoogleUser WHERE email = ?", key)

	} else {
		dblog.Error("Not exist provider")
		return false, ""
	}

	defer rows.Close()

	if err != nil {
		dblog.Error("db Query failed: %s", err.Error())
	}

	var name string
	var found bool
	for rows.Next() {
		err := rows.Scan(&name)
		if err != nil {
			dblog.Critical("Scan error: %s", err.Error())
		}
		dblog.Debug("found user for provider: %s, key: %s, name: %s", provider, key, name)
		found = true
	}

	err = rows.Err()
	if err != nil {
		dblog.Debug(err.Error())
	}

	return found, name

}

//GetGUser get all G user
func GetGUser() map[string]string {
	var rows *sql.Rows
	var err error

	rows, err = db.Query("SELECT email, name  FROM GoogleUser")
	if err != nil {
		dblog.Error("db Query failed: %s", err.Error())
	}

	defer rows.Close()

	users := make(map[string]string)
	var email string
	var name string
	for rows.Next() {
		err := rows.Scan(&email, &name)
		if err != nil {
			dblog.Critical("Scan error: %s", err.Error())
		}
		users[email] = name
	}

	err = rows.Err()
	if err != nil {
		dblog.Debug(err.Error())
	}

	return users
}

//GetFUser get all F user
func GetFUser() map[string]string {
	var rows *sql.Rows
	var err error

	rows, err = db.Query("SELECT id, name  FROM FacebookUser")
	if err != nil {
		dblog.Error("db Query failed: %s", err.Error())
	}

	defer rows.Close()

	users := make(map[string]string)
	var id string
	var name string
	for rows.Next() {
		err := rows.Scan(&id, &name)
		if err != nil {
			dblog.Critical("Scan error: %s", err.Error())
		}
		users[id] = name
	}

	err = rows.Err()
	if err != nil {
		dblog.Debug(err.Error())
	}

	return users
}

//Vtime for converting time.Time data
type Vtime []byte

//Time for converting time.Time data
func (t *Vtime) Time() (time.Time, error) {
	return time.Parse(mysqlDateFormat, string(*t))
}

//AddLocalUser add local user
func AddLocalUser(user, em, pw string) error {

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	now := time.Now()
	_, err = db.Exec("INSERT INTO LocalUser(name, email, pw, created) VALUES(?, ?, ?, ?)", user, em, hashedPassword, now.Format(mysqlDateFormat))

	if err != nil {
		dblog.Debug(err.Error())
	}

	return nil
}

//CheckLocalUser check if local user exist
func CheckLocalUser(em string) (user string, b bool) {
	err := db.QueryRow("SELECT name FROM LocalUser WHERE email=?", em).Scan(&user)
	switch {
	case err == sql.ErrNoRows:
		b = false
	case err != nil:
		dblog.Error("CheckLocalUser: %s", err.Error())
	default:
		b = true
	}
	return
}

//LoginLocalUser try to login local user
func LoginLocalUser(em, pw string) error {

	_, b := CheckLocalUser(em)
	if !b {
		return errors.New(em + " not exist, please check if you registered.")
	}

	var dpw string
	err := db.QueryRow("SELECT pw FROM LocalUser WHERE email=?", em).Scan(&dpw)
	switch {
	case err == sql.ErrNoRows:
		dblog.Error("System error, the user should be there with password")
	case err != nil:
		dblog.Error("LoginLocalUser: %s", err.Error())
	default:
		err = bcrypt.CompareHashAndPassword([]byte(dpw), []byte(pw))
		if err != nil {
			return errors.New("Wrong password for user " + em)
		}
	}
	return nil

}
