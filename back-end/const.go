package main

const (
	userName  string = "uestcproject"
	password  string = "MHdW6wfbxpY4TEkF"
	ipAddrees string = "localhost"
	port      int    = 3306
	dbName    string = "uestcproject"
	charset   string = "utf8"
)

/*
	type user_account struct {
		uid         int
		username    string
		password    string
		email       string
		access_key  string
		expire_time int
	}
*/
type feedback struct {
	Idfeedback int    `json:"id"`
	Username   string `json:"username"`
	Time       int    `json:"time"`
	Text       string `json:"text"`
}

type commonResponse struct {
	Ret int    `json:"ret"`
	Msg string `json:"msg"`
}

type FileInfo struct {
	Name       string `json:"name"`    // base name of the file
	Path       string `json:"path"`    // path
	Size       uint64 `json:"size"`    // length in bytes for regular files; system-dependent for others
	ModTime    uint   `json:"modTime"` // modification time
	UploadTime uint   `json:"uploadTime"`
	Md5        string `json:"md5"`
}

type FileHistory struct {
	ModTime    uint   `json:"modTime"` // modification time
	Size       uint64 `json:"size"`
	UploadTime uint   `json:"uploadTime"`
	Deleted    bool   `json:"deleted"`
	Md5        string `json:"md5"`
}
