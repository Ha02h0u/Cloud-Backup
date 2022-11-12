package main

import (
	"log"
	"net/http"

	"github.com/jmoiron/sqlx"
)

var DB *sqlx.DB

func main() {
	DB = connectMysql()
	defer DB.Close()
	http.HandleFunc("/", httpIndex)
	http.HandleFunc("/environment", httpEnvironment)
	http.HandleFunc("/login", httpLogin)
	http.HandleFunc("/register", httpRegister)
	http.HandleFunc("/refresh", httpRefresh)
	http.HandleFunc("/upload", httpUpload)
	http.HandleFunc("/download", httpDownload)
	http.HandleFunc("/filehistory", httpFileHistory)
	http.HandleFunc("/delete", httpDelete)
	http.HandleFunc("/filelist", httpFileList)
	http.HandleFunc("/recovery", httpRecovery)
	http.HandleFunc("/feedback", httpFeedback)
	http.HandleFunc("/getfeedback", httpGetFeedback)
	http.HandleFunc("/modifypassword", httpModifyPassword)
	if err := http.ListenAndServe(":2022", nil); err != nil {
		log.Fatal("HTTP服务启动失败")
	}

}
