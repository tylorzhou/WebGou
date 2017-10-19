package baapDB

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/go.uuid"
)

//Rmmeinsert insert remember me info
func Rmmeinsert(user, hash string, expiration time.Time, logintype int) (selector string, err error) {
	addError := errors.New("Rmmeinsert failed")

	stmtIns, err := db.Prepare("INSERT INTO Rememberme SET selector=?, user=?,hash=?,expiration=?,logintype=?")
	if err != nil {
		dblog.Error("db prepare failed: %s", err.Error())
		return
	}
	defer stmtIns.Close()

	uniqueID := uuid.NewV1()
	selector = uniqueID.String()

	result, err := stmtIns.Exec(selector, user, hash, expiration, logintype)
	if err != nil {
		dblog.Error("db Exec failed: %s", err.Error())
		return "", addError
	}

	insertid, err := result.LastInsertId()
	if err != nil {
		dblog.Error("db LastInsertId failed: %s", err.Error())
		return "", addError
	}

	dblog.Debug("Rmmeinsert, selector: %s, user: %s hash: %s, insertid:%d from rememberme", selector, user, user, insertid)
	return
}

//Rmmeget get remember me info
func Rmmeget(selector string) (user string, hash string, expiration time.Time, err error) {

	err = db.QueryRow("SELECT user, hash, expiration FROM Rememberme WHERE selector=?", selector).Scan(&user, &hash, &expiration)
	switch {
	case err == sql.ErrNoRows:
		dblog.Debug("no rememberme for %s", selector)
	case err != nil:
		dblog.Error("Rmmeget: %s", err.Error())
	default:
		fmt.Printf("get user info for %s from rememberme\n", selector)
	}
	return
}

//Rmmeupdate update remember me info
func Rmmeupdate(selector, user, hash string, expiration time.Time) (err error) {
	addError := errors.New("Rmmeupdate failed")

	stmtIns, err := db.Prepare("UPDATE Rememberme SET validator=?, expiration =? WHERE selector=?")
	if err != nil {
		dblog.Error("db prepare failed: %s", err.Error())
	}
	defer stmtIns.Close()

	result, err := stmtIns.Exec(hash, expiration, selector)
	if err != nil {
		dblog.Error("db Exec failed: %s", err.Error())
		return addError
	}

	insertid, err := result.LastInsertId()
	if err != nil {
		dblog.Error("db LastInsertId failed: %s", err.Error())
		return addError
	}

	dblog.Debug("Rmmeupdate, selector: %s, user: %s hash: %s, insertid:%d from rememberme", selector, user, user, insertid)
	return
}
