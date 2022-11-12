package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func connectMysql() *sqlx.DB {
	target := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s", userName, password, ipAddrees, port, dbName, charset)
	Db, err := sqlx.Open("mysql", target)
	if err != nil {
		log.Fatalf("MySQL connecting failed: %v", err)
	}
	Db.SetConnMaxLifetime(0)
	return Db
}

func QueryMysqlVersion(Db *sqlx.DB) string {
	rows, err := Db.Query("SELECT VERSION();")
	var version string
	if err != nil {
		log.Printf("Query faied: %v", err)
		return ""
	}
	for rows.Next() {
		err := rows.Scan(&version)
		if err != nil {
			log.Printf("Fetch data failed: %v", err)
			return ""
		}
	}
	rows.Close()
	return version
}

func DbLogin(Db *sqlx.DB, username string, password string) (int, string) {
	if !ValidString(username, 40, false) {
		return -1, "Error: Invalid Username!"
	}
	if !ValidString(password, 40, false) {
		return -1, "Error: Invalid Password!"
	}

	target := fmt.Sprintf("select password from user_account where username='%s';", username)
	rows, err := Db.Query(target)
	if err != nil {
		return -1, err.Error()
	}
	defer rows.Close()
	var pwd string
	for rows.Next() {
		if err := rows.Scan(&pwd); err != nil {
			log.Printf("Fetch data failed: %v", err)
			return -1, err.Error()
		}
		if pwd == password {
			ak := RandStr(45)
			target = fmt.Sprintf("update user_account set access_key='%s', expire_time=%d where username='%s';", ak, time.Now().Unix()+300, username)
			_, err = Db.Exec(target)
			if err != nil {
				log.Printf("Exec %s faied: %v", target, err)
				return -1, err.Error()
			} else {
				return 0, ak
			}
		} else {
			return -1, "Wrong Password!"
		}
	}
	return -1, "No such username, please register first!"
}

func DbModifyPassword(Db *sqlx.DB, username, password, newPassword string) (int, string) {
	if !ValidString(username, 40, false) {
		return -1, "Error: Invalid Username!"
	}
	if !ValidString(password, 40, false) {
		return -1, "Error: Invalid Password!"
	}
	if !ValidString(newPassword, 40, false) {
		return -1, "Error: Invalid newPassword!"
	}

	target := fmt.Sprintf("select password from user_account where username='%s';", username)
	rows, err := Db.Query(target)
	if err != nil {
		return -1, err.Error()
	}
	defer rows.Close()
	var pwd string
	for rows.Next() {
		if err := rows.Scan(&pwd); err != nil {
			log.Printf("Fetch data failed: %v", err)
			return -1, err.Error()
		}
		if pwd == password {
			target = fmt.Sprintf("update user_account set password='%s' where username='%s';", newPassword, username)
			_, err = Db.Exec(target)
			if err != nil {
				log.Printf("Exec %s faied: %v", target, err)
				return -1, err.Error()
			} else {
				return 0, "success!"
			}
		} else {
			return -1, "Wrong Password!"
		}
	}
	return -1, "No such username, please register first!"
}

func DbRefresh(Db *sqlx.DB, username string, accessKey string) (int, string) {
	if !ValidString(username, 40, false) {
		return -1, "Error: Invalid Username!"
	}
	if len(accessKey) != 44 {
		return -1, "Error: Invalid AccessKey!"
	}

	target := fmt.Sprintf("select access_key, expire_time from user_account where username='%s';", username)
	rows, err := Db.Query(target)
	if err != nil {
		log.Printf("Query %s faied: %v", target, err)
		return -1, err.Error()
	}
	defer rows.Close()
	var ak string
	var expire_time int
	for rows.Next() {
		if err := rows.Scan(&ak, &expire_time); err != nil {
			log.Printf("Fetch data failed: %v", err)
			return -1, err.Error()
		}
		if int(time.Now().Unix()) > expire_time {
			return -1, "AccessKey expired, please login again!"
		}
		if ak == accessKey {
			target = fmt.Sprintf("update user_account set expire_time=%d where username='%s';", time.Now().Unix()+300, username)
			_, err = Db.Exec(target)
			if err != nil {
				log.Printf("Exec %s faied: %v", target, err)
				return -1, err.Error()
			} else {
				return 0, "Success!"
			}
		} else {
			return -1, "Wrong AccessKey!"
		}
	}
	return -1, "No such username, please register first!"
}

func DbRegister(Db *sqlx.DB, username string, password string, email string) (int, string) {
	if !ValidString(username, 40, false) {
		return -1, "Error: Invalid Username!"
	}
	if !ValidString(password, 40, false) {
		return -1, "Error: Invalid Password!"
	}
	if !ValidEmailAddress(email, 40) {
		return -1, "Error: Invalid Email!"
	}

	target := fmt.Sprintf("select 1 from user_account where username='%s' limit 1;", username)
	rows, err := Db.Query(target)
	if err != nil {
		log.Printf("Query %s faied: %v", target, err)
		return -1, err.Error()
	}
	defer rows.Close()
	if rows.Next() {
		return -1, "Username exists, please use login function!"
	}
	ak := RandStr(45)
	target = fmt.Sprintf("insert into user_account(username, password, email, access_key, expire_time) values('%s','%s','%s','%s',%d);", username, password, email, ak, time.Now().Unix()+300)
	_, err = Db.Exec(target)
	if err != nil {
		log.Printf("Exec %s faied: %v", target, err)
		return -1, err.Error()
	}
	return 0, ak
}

func DbFeedback(Db *sqlx.DB, username string, text string) (int, string) {
	if len(text) > 1000 || len(text) <= 0 {
		return -1, "Error: Invalid Text!"
	}

	target := fmt.Sprintf("insert into feedback(username, text, time) values('%s','%s',%d);", username, text, time.Now().Unix())
	_, err := Db.Exec(target)
	if err != nil {
		log.Printf("Exec %s faied: %v", target, err)
		return -1, err.Error()
	}
	return 0, "Success!"
}

func DbGetFeedback(Db *sqlx.DB, size int, page int) (int, string) {
	target := fmt.Sprintf("select * from feedback order by idfeedback DESC limit %d, %d;", (page-1)*size, size)
	rows, err := Db.Query(target)
	if err != nil {
		log.Printf("Query %s faied: %v", target, err)
		return -1, err.Error()
	}
	defer rows.Close()
	var fb []feedback
	for rows.Next() {
		var tmp feedback
		if err := rows.Scan(&tmp.Idfeedback, &tmp.Username, &tmp.Text, &tmp.Time); err != nil {
			log.Printf("Fetch data failed: %v", err)
			return -1, err.Error()
		}
		fb = append(fb, tmp)
	}
	res, err := json.Marshal(fb)
	if err != nil {
		log.Printf("DbGetFeedback JSON Marshal failed: %v", err)
		return -1, err.Error()
	}
	return 0, string(res)
}

func DbGetUid(Db *sqlx.DB, username string) int {
	target := fmt.Sprintf("select uid from user_account where username='%s';", username)
	rows, err := Db.Query(target)
	if err != nil {
		return -1
	}
	defer rows.Close()
	for rows.Next() {
		var uid int
		if err := rows.Scan(&uid); err != nil {
			log.Printf("Fetch data failed: %v", err)
			return -1
		} else {
			return uid
		}
	}
	log.Printf("Fetch data failed: %v", err)
	return -1
}

func DbAddFileLog(Db *sqlx.DB, username, path, name, md5 string, size uint64, modTime uint) (int, string) {
	uid := DbGetUid(Db, username)
	if uid <= 0 {
		return -1, "Error: Can't fetch user data!"
	}
	target := fmt.Sprintf("insert into file_info(uid, path, name, size, modTime, md5, uploadTime, deleted) values(%d,'%s','%s',%v,%v,'%s',%v,FALSE);", uid, path, name, size, modTime, md5, time.Now().Unix())
	res, err := Db.Exec(target)
	if err != nil {
		log.Printf("Exec %s faied: %v", target, err)
		return -1, err.Error()
	}
	lastInsertId, _ := res.LastInsertId()
	return 0, fmt.Sprintf("%v", lastInsertId)
}

func DbDeleteFile(Db *sqlx.DB, username, path, name, md5 string) (int, string) {
	uid := DbGetUid(Db, username)
	if uid <= 0 {
		return -1, "Error: Can't fetch user data!"
	}
	if md5 == "" {
		ret, msg, fi := DbGetLatestFileInfo(Db, username, path, name)
		if ret != 0 {
			return ret, msg
		}
		md5 = fi.Md5
	}
	target := fmt.Sprintf("select deleted from file_info where uid=%d and path='%s' and name='%s' and md5='%s';", uid, path, name, md5)
	if rows, err := Db.Query(target); err != nil {
		log.Printf("Query %s faied: %v", target, err)
		return -1, err.Error()
	} else {
		defer rows.Close()
		for rows.Next() {
			var deleted bool
			rows.Scan(&deleted)
			if deleted {
				return -1, "file has already been deleted!"
			}
		}
	}
	target = fmt.Sprintf("update file_info set deleted=TRUE where uid=%d and path='%s' and name='%s' and md5='%s';", uid, path, name, md5)
	res, err := Db.Exec(target)
	if err != nil {
		log.Printf("Exec %s faied: %v", target, err)
		return -1, err.Error()
	}
	lastInsertId, _ := res.LastInsertId()
	return 0, fmt.Sprintf("%v", lastInsertId)
}

func DbGetLatestFileInfo(Db *sqlx.DB, username, path, name string) (int, string, FileInfo) {
	uid := DbGetUid(Db, username)
	if uid <= 0 {
		return -1, "Error: Can't fetch user data!", FileInfo{}
	}
	target := fmt.Sprintf("select name, path, size, modTime, md5, uploadTime, deleted from file_info where uid = %d and path = '%s' and name = '%s' order by uploadTime desc limit 1;", uid, path, name)
	rows, err := Db.Query(target)
	if err != nil {
		log.Printf("Query %s faied: %v", target, err)
		return -1, err.Error(), FileInfo{}
	}
	defer rows.Close()
	for rows.Next() {
		var tmp FileInfo
		var deleted bool
		rows.Scan(&tmp.Name, &tmp.Path, &tmp.Size, &tmp.ModTime, &tmp.Md5, &tmp.UploadTime, &deleted)
		if !deleted {
			return 0, "", tmp
		} else {
			return -2, "File has been deleted!", FileInfo{}
		}
	}
	return -3, "no such file history in database!", FileInfo{}
}

func DbGetFileHistory(Db *sqlx.DB, username, path, name string) (int, string) {
	uid := DbGetUid(Db, username)
	if uid <= 0 {
		return -1, "Error: Can't fetch user data!"
	}
	target := fmt.Sprintf("select deleted, modTime, uploadTime, size, md5 from (select * from file_info where uid = %d and path = '%s' and name = '%s' order by uploadTime desc limit 10000000000) as A group by A.md5 order by uploadTime desc;", uid, path, name)
	rows, err := Db.Query(target)
	if err != nil {
		log.Printf("Query %s faied: %v", target, err)
		return -1, err.Error()
	}
	defer rows.Close()
	var fh []FileHistory
	for rows.Next() {
		var tmp FileHistory
		rows.Scan(&tmp.Deleted, &tmp.ModTime, &tmp.UploadTime, &tmp.Size, &tmp.Md5)
		fh = append(fh, tmp)
	}
	res, err := json.Marshal(fh)
	if err != nil {
		log.Printf("DbFileHistory JSON Marshal failed: %v", err)
		return -1, err.Error()
	}
	return 0, string(res)
}

func DbGetFileList(Db *sqlx.DB, username, filter string, onlyDeleted bool) (int, string) {
	uid := DbGetUid(Db, username)
	if uid <= 0 {
		return -1, "Error: Can't fetch user data!"
	}
	target := fmt.Sprintf("select name, path, size, modTime, md5, uploadTime, deleted from (select * from file_info where uid = %d order by uploadTime desc limit 10000000000) as A group by A.path,A.name order by path,name;", uid)
	rows, err := Db.Query(target)
	if err != nil {
		log.Printf("Query %s faied: %v", target, err)
		return -1, err.Error()
	}
	defer rows.Close()
	var fi []FileInfo
	for rows.Next() {
		var tmp FileInfo
		var deleted bool
		rows.Scan(&tmp.Name, &tmp.Path, &tmp.Size, &tmp.ModTime, &tmp.Md5, &tmp.UploadTime, &deleted)
		if deleted == onlyDeleted && strings.Contains(tmp.Name, filter) {
			fi = append(fi, tmp)
		}
	}
	res, err := json.Marshal(fi)
	if err != nil {
		log.Printf("DbFileList JSON Marshal failed: %v", err)
		return -1, err.Error()
	}
	return 0, string(res)
}
