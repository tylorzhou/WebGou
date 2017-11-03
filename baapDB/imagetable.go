package baapDB

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

var createImageStatements = `CREATE TABLE IF NOT EXISTS %s (
		uid INT NOT NULL AUTO_INCREMENT,
		logintype INT NOT NULL,
		id INT NOT NULL,
		imageurl VARCHAR(255) NOT NULL,
		created DATETIME NULL DEFAULT NULL,
		description TEXT NULL,
		PRIMARY KEY (uid)
	)`

//Imageinfo for image information
type Imageinfo struct {
	Uqid                  int64
	Logintype, ID         int
	Imageurl, Description string
	Created               time.Time
}

// ensureTableExists checks the table exists. If not, it creates it.
func ensureTableExists(tablename string) error {

	/*
		if _, err := db.Exec("USE BaapAPI"); err != nil {
			// MySQL error 1049 is "database does not exist"
			if mErr, ok := err.(*mysql.MySQLError); ok && mErr.Number == 1049 {
				return createTable(db, createtable)
			}
		}
	*/

	/* 	Check := fmt.Sprintf("DESCRIBE %s", tablename) */

	_, err := db.Exec("USE BaapAPI;")
	if err != nil {
		dblog.Error("mysql: could not act on database: %s", err.Error())
		return err
	}

	Create := fmt.Sprintf(createImageStatements, tablename)
	_, err = db.Exec(Create)
	if err != nil {
		dblog.Error("mysql: could not act on database: %s", err.Error())
	}
	return err
}

//InsertImage record for image upload
func InsertImage(tablename string, image Imageinfo) (insertid int64, err error) {

	err = ensureTableExists(tablename)
	if err != nil {
		return 0, err
	}

	addError := errors.New("InsertImage failed")

	statement := fmt.Sprintf("INSERT INTO %s SET logintype=?, id=?,imageurl=?,description=?, created=?", tablename)
	stmtIns, err := db.Prepare(statement)
	if err != nil {
		dblog.Error("db prepare failed: %s", err.Error())
		return
	}
	defer stmtIns.Close()

	result, err := stmtIns.Exec(image.Logintype, image.ID, image.Imageurl, image.Description, image.Created)
	if err != nil {
		dblog.Error("db Exec failed: %s", err.Error())
		return 0, addError
	}

	insertid, err = result.LastInsertId()
	if err != nil {
		dblog.Error("db LastInsertId failed: %s", err.Error())
		return 0, addError
	}

	dblog.Debug("InsertImage Successfully, tablename %s, image %v", tablename, image)
	return
}

//GetAllImages to get all images from table
func GetAllImages(tablename string) ([]Imageinfo, error) {
	statement := fmt.Sprintf("SELECT uid, logintype, id, imageurl, description, created FROM %s ORDER BY created DESC",
		tablename)
	var rows *sql.Rows
	var err error

	rows, err = db.Query(statement)
	if err != nil {
		dblog.Error("GetAllImages db Query failed: %s", err.Error())
	}

	defer rows.Close()

	imageinfo := make([]Imageinfo, 0, 50)

	for rows.Next() {
		var img Imageinfo
		err := rows.Scan(&img.Uqid, &img.Logintype, &img.ID, &img.Imageurl, &img.Description, img.Created)
		if err != nil {
			dblog.Critical("GetAllImages Scan error: %s", err.Error())
		}
		imageinfo = append(imageinfo, img)
	}

	err = rows.Err()
	if err != nil {
		dblog.Debug(err.Error())
		return nil, err
	}

	return imageinfo, nil
}

//GImgTblName get image table name
func GImgTblName(logintype, id int) string {
	return fmt.Sprintf("img_%d_%d", logintype, id)
}
