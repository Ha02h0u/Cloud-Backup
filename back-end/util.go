package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/mail"
	"os"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/net"
	"github.com/shirou/gopsutil/v3/mem"
)

type LSysInfo struct {
	MemAll         uint64
	MemFree        uint64
	MemUsed        uint64
	MemUsedPercent float64
	Days           int64
	Hours          int64
	Minutes        int64
	Seconds        int64
}

func init() {
	// 保证每次生成的随机数不一样
	rand.Seed(time.Now().UnixNano())
}

func sysInfo() (info LSysInfo) {
	v, _ := mem.VirtualMemory()
	info.MemAll = v.Total
	info.MemFree = v.Free
	info.MemUsed = info.MemAll - info.MemFree
	// 注：使用SwapMemory或VirtualMemory，在不同系统中使用率不一样，因此直接计算一次
	info.MemUsedPercent = float64(info.MemUsed) / float64(info.MemAll) * 100.0 // v.UsedPercent

	// 获取开机时间
	boottime, _ := host.BootTime()
	ntime := time.Now().Unix()
	btime := time.Unix(int64(boottime), 0).Unix()
	deltatime := ntime - btime

	info.Seconds = int64(deltatime)
	info.Minutes = info.Seconds / 60
	info.Seconds -= info.Minutes * 60
	info.Hours = info.Minutes / 60
	info.Minutes -= info.Hours * 60
	info.Days = info.Hours / 24
	info.Hours -= info.Days * 24
	return
}

func infoTest() (res string) {
	info := sysInfo()
	unitGb := uint64(1024 * 1024 * 1024)
	unitMb := uint64(1024 * 1024)
	c, _ := cpu.Info()
	cc, _ := cpu.Percent(time.Second, true) // 1秒
	d, _ := disk.Usage("/")
	n, _ := host.Info()
	nv, _ := net.IOCounters(true)
	physicalCnt, _ := cpu.Counts(false)
	logicalCnt, _ := cpu.Counts(true)
	for _, sub_cpu := range c {
		modelname := sub_cpu.ModelName
		cores := sub_cpu.Cores
		for j := int32(0); j < cores; j++ {
			res += fmt.Sprintf("CPU: %s(%d核) 占用率: %f%%\n", modelname, cores, cc[j])
		}
	}
	res += fmt.Sprintf("物理核心数: %d 逻辑核心数: %d\n", physicalCnt, logicalCnt)
	res += fmt.Sprintf("内存可用量: %d/%d GB (%f%%)\n", info.MemFree/unitGb, info.MemAll/unitGb, info.MemUsedPercent)
	res += fmt.Sprintf("磁盘可用量: %d/%d GB (%f%%)\n", d.Free/unitGb, d.Total/unitGb, d.UsedPercent)
	res += fmt.Sprintf("操作系统: %s(%s) %s\n", n.Platform, n.PlatformFamily, n.PlatformVersion)
	res += fmt.Sprintf("登陆用户: %s\n", n.Hostname)
	for i, value := range nv {
		res += fmt.Sprintf("网络%d(%s): 共上传%dMB/下载%vMB\n", i, value.Name, value.BytesRecv/unitMb, value.BytesSent/unitMb)
	}
	res += fmt.Sprintf("系统已经运行了: %d天%d时%d分%d秒\n", info.Days, info.Hours, info.Minutes, info.Seconds)
	return
}

func RandStr(n int) string {
	result := make([]byte, n/2)
	rand.Read(result)
	return hex.EncodeToString(result)
}

func ValidString(str string, length int, onlyNumberAndAlphabet bool) bool {
	if len(str) > length || len(str) < 5 {
		return false
	}
	for i := range str {
		if i == '\'' || i == '"' || i == '&' || i == '|' || i == '/' || i == '\\' {
			return false
		}
		if i >= 'a' && i <= 'z' {
			continue
		}
		if i >= 'A' && i <= 'Z' {
			continue
		}
		if i >= '0' && i <= '9' {
			continue
		}
		if onlyNumberAndAlphabet {
			return false
		}
	}
	return true
}

func ValidEmailAddress(email string, length int) bool {
	if len(email) > length {
		return false
	}
	_, err := mail.ParseAddress(email)
	return err == nil
}

func isFileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, err
	}
	return false, err
}

func CommonResponse2String(cRsp commonResponse) string {
	res, err := json.MarshalIndent(cRsp, "", "  ")
	if err != nil {
		log.Println(err)
	}
	return string(res)
}

func FileMD5(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()
	hash := md5.New()
	_, _ = io.Copy(hash, file)
	return hex.EncodeToString(hash.Sum(nil)), nil
}

func AesEncrypt(origData []byte, key string) []byte {
	// 转成字节数组
	k := []byte(key)
	for len(key) < 16 {
		key += "0"
	}
	// 分组秘钥
	block, _ := aes.NewCipher(k)
	// 获取秘钥块的长度
	blockSize := block.BlockSize()
	// 补全码
	origData = PKCS7Padding(origData, blockSize)
	// 加密模式
	blockMode := cipher.NewCBCEncrypter(block, k[:blockSize])
	// 创建数组
	cryted := make([]byte, len(origData))
	// 加密
	blockMode.CryptBlocks(cryted, origData)
	return cryted

}

func AesDecrypt(crytedByte []byte, key string) []byte {
	k := []byte(key)
	for len(key) < 16 {
		key += "0"
	}
	// 分组秘钥
	block, _ := aes.NewCipher(k)
	// 获取秘钥块的长度
	blockSize := block.BlockSize()
	// 加密模式
	blockMode := cipher.NewCBCDecrypter(block, k[:blockSize])
	// 创建数组
	orig := make([]byte, len(crytedByte))
	// 解密
	blockMode.CryptBlocks(orig, crytedByte)
	// 去补全码
	return PKCS7UnPadding(orig)
}

// 补码
func PKCS7Padding(ciphertext []byte, blocksize int) []byte {
	padding := blocksize - len(ciphertext)%blocksize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

// 去码
func PKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}
