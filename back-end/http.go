package main

import (
	"fmt"
	"io"
	"strings"

	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

func httpIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	f, err := os.Open("html/index.html")
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}
	defer f.Close()
	io.Copy(w, f)
}

func httpEnvironment(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	fmt.Fprintln(w, infoTest())
	fmt.Fprintf(w, "Mysql version: %s", QueryMysqlVersion(DB))
}

func httpLogin(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	username := r.PostFormValue("username")
	password := r.PostFormValue("password")
	log.Printf("[httpLogin] username:%s password:%s", username, password)
	ret, msg := DbLogin(DB, username, password)
	w.Write([]byte(CommonResponse2String(commonResponse{ret, msg})))
}

func httpModifyPassword(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	username := r.PostFormValue("username")
	password := r.PostFormValue("password")
	newPassword := r.PostFormValue("newPassword")
	log.Printf("[httpModifyPassword] username:%s password:%s newPassword:%s", username, password, newPassword)
	ret, msg := DbModifyPassword(DB, username, password, newPassword)
	w.Write([]byte(CommonResponse2String(commonResponse{ret, msg})))
}

func httpRegister(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	username := r.PostFormValue("username")
	password := r.PostFormValue("password")
	email := r.PostFormValue("email")
	log.Printf("[httpRegister] username:%s password:%s email:%s", username, password, email)
	ret, msg := DbRegister(DB, username, password, email)
	w.Write([]byte(CommonResponse2String(commonResponse{ret, msg})))
}

func httpRefresh(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	username := r.PostFormValue("username")
	ak := r.PostFormValue("AccessKey")
	log.Printf("refresh Trys: %s %s", username, ak)
	ret, msg := DbRefresh(DB, username, ak)
	w.Write([]byte(CommonResponse2String(commonResponse{ret, msg})))
}

func httpFeedback(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	username := r.PostFormValue("username")
	ak := r.PostFormValue("AccessKey")
	text := r.PostFormValue("text")
	log.Printf("feedback Trys: %s %s %s", username, ak, text)
	if ret, msg := httpVerify(username, ak); ret != 0 {
		log.Println(msg)
		w.Write([]byte(CommonResponse2String(commonResponse{ret, msg})))
		return
	}
	ret, msg := DbFeedback(DB, username, text)
	w.Write([]byte(CommonResponse2String(commonResponse{ret, msg})))
}

func httpGetFeedback(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	username := r.PostFormValue("username")
	ak := r.PostFormValue("AccessKey")
	size, _ := strconv.Atoi(r.PostFormValue("size"))
	page, _ := strconv.Atoi(r.PostFormValue("page"))
	if size <= 0 {
		size = 5
	}
	if page <= 0 {
		page = 1
	}
	log.Printf("getFeedback Trys: %s %s %d %d", username, ak, size, page)
	if ret, msg := httpVerify(username, ak); ret != 0 {
		log.Println(msg)
		w.Write([]byte(CommonResponse2String(commonResponse{ret, msg})))
		return
	}
	ret, msg := DbGetFeedback(DB, size, page)
	w.Write([]byte(CommonResponse2String(commonResponse{ret, msg})))
}

func httpRecovery(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	username := r.PostFormValue("username")
	ak := r.PostFormValue("AccessKey")
	filter := r.PostFormValue("filter")
	log.Printf("recovery Trys: %s %s %s", username, ak, filter)
	if ret, msg := httpVerify(username, ak); ret != 0 {
		log.Println(msg)
		w.Write([]byte(CommonResponse2String(commonResponse{ret, msg})))
		return
	}
	ret, msg := DbGetFileList(DB, username, filter, true)
	w.Write([]byte(CommonResponse2String(commonResponse{ret, msg})))
}

func httpFileList(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	username := r.PostFormValue("username")
	ak := r.PostFormValue("AccessKey")
	filter := r.PostFormValue("filter")
	log.Printf("filelist Trys: %s %s %s", username, ak, filter)
	if ret, msg := httpVerify(username, ak); ret != 0 {
		log.Println(msg)
		w.Write([]byte(CommonResponse2String(commonResponse{ret, msg})))
		return
	}
	ret, msg := DbGetFileList(DB, username, filter, false)
	w.Write([]byte(CommonResponse2String(commonResponse{ret, msg})))
}

func httpUpload(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	username := r.PostFormValue("username")
	ak := r.PostFormValue("AccessKey")
	log.Printf("upload Trys: %s %s", username, ak)
	if ret, msg := httpVerify(username, ak); ret != 0 {
		log.Println(msg)
		w.Write([]byte(CommonResponse2String(commonResponse{ret, msg})))
		return
	}
	md5 := r.PostFormValue("md5")
	md5 = strings.ToLower(md5)
	path := r.PostFormValue("path")
	path = strings.ReplaceAll(path, "\\", "/")
	modTime_str := r.PostFormValue("modTime")
	modTime, strconv_err := strconv.Atoi(modTime_str)
	if strconv_err != nil {
		log.Println(strconv_err.Error())
		w.Write([]byte(CommonResponse2String(commonResponse{-1, strconv_err.Error()})))
		return
	}
	if len(md5) != 32 || modTime < 0 || path[len(path)-1:] != "/" || path[0] != '/' {
		log.Println("Invalid params!")
		w.Write([]byte(CommonResponse2String(commonResponse{-1, "Invalid params!"})))
		return
	}
	upfile, fileHeader, err := r.FormFile("file")
	if err != nil {
		log.Println(err.Error())
		w.Write([]byte(CommonResponse2String(commonResponse{-1, "Error: upload failed, " + err.Error()})))
		return
	}
	log.Printf("path: %s\nfileName: %s\nsize: %dBytes\n", path, fileHeader.Filename, fileHeader.Size)
	filePath := "./userdata/" + username + "/" + path + fileHeader.Filename + "/" + md5 + filepath.Ext(fileHeader.Filename)
	ret, msg, fi := DbGetLatestFileInfo(DB, username, path, fileHeader.Filename)
	needNewLog := false
	if ret == -1 {
		log.Println(msg)
		w.Write([]byte(CommonResponse2String(commonResponse{ret, msg})))
		return
	} else if ret == 0 {
		fmt.Println("database md5:" + fi.Md5)
		if md5 != fi.Md5 {
			fmt.Println("md5 different, insert and create local file")
			needNewLog = true
		}
	} else if ret == -2 {
		fmt.Println("file deleted, insert and create local file")
		needNewLog = true
	} else {
		fmt.Println("No database info, insert and create local file")
		needNewLog = true
	}
	if needNewLog {
		realpath, _ := filepath.Split(filePath)
		os.MkdirAll(realpath, 0777)
		DbAddFileLog(DB, username, path, fileHeader.Filename, md5, uint64(fileHeader.Size), uint(modTime))
		newfile, err := os.Create(filePath)
		if err != nil {
			log.Println(err.Error())
			w.Write([]byte(CommonResponse2String(commonResponse{-1, "Error: Create file failed, " + err.Error()})))
			return
		} else {
			newfile.Close()
		}
	}
	fsFileInfo, _ := os.Stat(filePath)
	localFileSize := fsFileInfo.Size()
	if localFileSize == fileHeader.Size {
		tmp := fmt.Sprintf("file already exists, size %v", localFileSize)
		w.Write([]byte(CommonResponse2String(commonResponse{0, tmp})))
		return
	}
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		log.Println(err.Error())
		w.Write([]byte(CommonResponse2String(commonResponse{-1, "Error: File not exists, " + err.Error()})))
		return
	} else {
		defer file.Close()
	}
	fmt.Println("Local file size:", localFileSize)
	upfile.Seek(localFileSize, 0)
	file.Seek(localFileSize, 0)
	data := make([]byte, 1024)
	var upTotal int64 = 0
	for {
		total, err := upfile.Read(data)
		if err == io.EOF {
			fmt.Println("File uploading finished!")
			break
		}
		len, err := file.Write(data[:total])
		if err != nil {
			tmp := fmt.Sprintf("Error: upload fail. Already uploaded Bytes: %v", localFileSize)
			w.Write([]byte(CommonResponse2String(commonResponse{-1, tmp})))
			return
		}
		upTotal += int64(len)
		localFileSize += int64(len)
	}
	tmp := fmt.Sprintf("upload success! Total Bytes: %d", upTotal)
	w.Write([]byte(CommonResponse2String(commonResponse{0, tmp})))
}

func httpFileHistory(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	username := r.PostFormValue("username")
	ak := r.PostFormValue("AccessKey")
	filePath := r.PostFormValue("filePath")
	filePath = strings.ReplaceAll(filePath, "\\", "/")
	log.Printf("httpFileHistory Trys: %s %s %s", username, ak, filePath)
	if ret, msg := httpVerify(username, ak); ret != 0 {
		log.Println(msg)
		w.Write([]byte(CommonResponse2String(commonResponse{ret, msg})))
		return
	}
	if filePath[0] != '/' {
		log.Println("Error: illegal filePath!")
		w.Write([]byte(CommonResponse2String(commonResponse{-1, "Error: Illegal filePath!"})))
		return
	}
	fileName := filepath.Base(filePath)
	fileDir := filepath.Dir(filePath)
	if fileDir[len(fileDir)-1:] != "/" {
		fileDir = fileDir + "/"
	}
	ret, msg := DbGetFileHistory(DB, username, fileDir, fileName)
	if ret != 0 {
		log.Println(msg)
	}
	w.Write([]byte(CommonResponse2String(commonResponse{ret, msg})))
}

func httpDownload(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	username := r.PostFormValue("username")
	ak := r.PostFormValue("AccessKey")
	filePath := r.PostFormValue("filePath")
	filePath = strings.ReplaceAll(filePath, "\\", "/")
	md5 := r.PostFormValue("md5")
	md5 = strings.ToLower(md5)
	log.Printf("download Trys: %s %s %s %s", username, ak, filePath, md5)
	if ret, msg := httpVerify(username, ak); ret != 0 {
		log.Println(msg)
		w.Write([]byte(CommonResponse2String(commonResponse{ret, msg})))
		return
	}
	if filePath[0] != '/' || (len(md5) > 0 && len(md5) != 32) {
		log.Println("Error: illegal filePath!")
		w.Write([]byte(CommonResponse2String(commonResponse{-1, "Error: Illegal filePath!"})))
		return
	}
	fileName := filepath.Base(filePath)
	fileDir := filepath.Dir(filePath)
	if fileDir[len(fileDir)-1:] != "/" {
		fileDir = fileDir + "/"
	}
	if md5 == "" {
		ret, msg, fi := DbGetLatestFileInfo(DB, username, fileDir, fileName)
		if ret != 0 {
			log.Println(msg)
			w.Write([]byte(CommonResponse2String(commonResponse{ret, msg})))
			return
		}
		md5 = fi.Md5
	}
	path := "./userdata/" + username + "/" + fileDir + fileName + "/" + md5 + filepath.Ext(fileName)
	if fileBool, _ := isFileExists(path); !fileBool {
		log.Println("Error: Local file Not Exist!")
		w.Write([]byte(CommonResponse2String(commonResponse{-1, "Error: Local file Not Exist!"})))
		return
	}
	w.Header().Add("Content-Type", "application/octet-stream")
	w.Header().Add("Content-Disposition", "attachment;filename="+fileName)
	f, err := os.Open(path)
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}
	defer f.Close()
	io.Copy(w, f)
}

func httpDelete(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	username := r.PostFormValue("username")
	ak := r.PostFormValue("AccessKey")
	filePath := r.PostFormValue("filePath")
	filePath = strings.ReplaceAll(filePath, "\\", "/")
	log.Printf("download Trys: %s %s %s", username, ak, filePath)
	if ret, msg := httpVerify(username, ak); ret != 0 {
		log.Println(msg)
		w.Write([]byte(CommonResponse2String(commonResponse{ret, msg})))
		return
	}
	if filePath[0] != '/' {
		log.Println("Error: illegal filePath!")
		w.Write([]byte(CommonResponse2String(commonResponse{-1, "Error: Illegal filePath!"})))
		return
	}
	fileName := filepath.Base(filePath)
	fileDir := filepath.Dir(filePath) + "/"
	if fileDir == "//" {
		fileDir = "/"
	}
	ret, msg, fi := DbGetLatestFileInfo(DB, username, fileDir, fileName)
	if ret != 0 {
		log.Println(msg)
		w.Write([]byte(CommonResponse2String(commonResponse{ret, msg})))
		return
	}
	ret, msg = DbDeleteFile(DB, username, fileDir, fileName, fi.Md5)
	if ret != 0 {
		fmt.Println(msg)
	}
	w.Write([]byte(CommonResponse2String(commonResponse{ret, msg})))
}

func httpVerify(username, ak string) (int, string) {
	if !ValidString(username, 40, false) {
		return -1, "Error: Invalid Username!"
	}
	if len(ak) != 44 {
		return -1, "Error: Invalid AccessKey!"
	}
	return DbRefresh(DB, username, ak)
}
