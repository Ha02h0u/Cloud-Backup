### Cloud-Backup
A cloud backup project for Comprehensive Experiment of Software Development Course held by UESTC.

### Introduction

#### Client folder
- 在线文件自动同步系统.exe: executable program running on Windows PC.
> Note: The software is written in [Easy Programming Language](http://www.dywt.com.cn/) (EPL), so it may cause fake alarms of anti-virus software. You can run this software with confidence and my primise

- 客户端.e: [EPL](http://www.dywt.com.cn/) Code, charge of Windows program part.
> Note: [EPL](http://www.dywt.com.cn/) Code is **not** plain text, so programmers have to **use special IDE** offered on its official website to program)

- 精易模块[v10.4.5].ec: [EPL](http://www.dywt.com.cn/) Module, charge of encoding/decoding & accessing web interface.

- high_dpi.ec: [EPL](http://www.dywt.com.cn/) Module, charge of fitting Windows high DPI.

- logo.gif: logo of our program.

#### back-end folder

- main.go: the entrance of back-end, charge of http listening.

- http.go: http codes, handle every http requests here and use mysql in many conditions.

- mysql.go: mysql codes, handle every encapsulated database request, and run MySQL statements.

- util.go: codes of some generic methods which may be used in many different places.

- const.go: declarations of structs and const settings.

- html/index.html: html file, for front-end developers to test interface conveniently.

- go.mod & go.sum: Go Mod files.

#### front-end folder

front-end codes, most of which written in HTML & JavaScript .

#### Contact
If you have any problem in using our software, please contact us immediately through QQ or E-mail.