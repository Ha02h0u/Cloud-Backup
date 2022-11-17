package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
)

var DB *sqlx.DB

func init() {
	os.MkdirAll("./log", 0777)
	file := "./log/" + time.Now().String() + ".txt"
	logFile, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0766)
	if err != nil {
		panic(err)
	}
	log.SetOutput(logFile) // 将文件设置为log输出的文件
	log.SetPrefix("[Cloud-Backup]")
	log.SetFlags(log.LstdFlags | log.Lshortfile | log.LUTC)
}

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
