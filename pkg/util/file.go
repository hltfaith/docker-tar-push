package util

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
	"path/filepath"
)

//Exists 判断所给路径文件/文件夹是否存在
func Exists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

// 所有文件绝对路径
func FilesPath(path string) ([]string, error) {
	fd, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	var filespath []string = []string{}
	if !fd.IsDir() {
		filespath = append(filespath, path)
		return filespath, nil
	}

	err = filepath.Walk(path, func(subPath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			filespath = append(filespath, subPath)
		}
		return err
	})
	return filespath, err
}

//GetFileSize 获取文件大小
func GetFileSize(file string) (int64, error) {
	f, err := os.Open(file)
	if err != nil {
		return 0, err
	}
	stat, err := f.Stat() //获取文件状态
	if err != nil {
		return 0, err
	}
	defer f.Close()
	return stat.Size(), nil
}

//Sha256Hash
func Sha256Hash(file string) (string, error) {
	f, err := os.Open(file)
	if err != nil {
		return "", err
	}
	defer f.Close()
	chunkSize := 65536
	buf := make([]byte, chunkSize)
	h := sha256.New()
	for {
		n, err := f.Read(buf)
		if err == io.EOF {
			break
		}
		chunk := buf[0:n]
		h.Write(chunk)
	}
	sum := h.Sum(nil)
	//由于是十六进制表示，因此需要转换
	hash := hex.EncodeToString(sum)
	return hash, nil
}
