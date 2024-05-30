package util

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

// Decompress 解压 tar.gz 保留原始的层级结构和文件修改时间
//
// tarFile 被解压的 .tar.gz文件名
//
// dest 解压到哪个目录, 结尾的 "/" 可有可无, "" 和 "./" 和 "." 都表示解压到当前目录
func Decompress(tarFile, dest string) error {
	srcFile, err := os.Open(tarFile)
	if err != nil {
		return err
	}
	defer srcFile.Close()
	gr, err := gzip.NewReader(srcFile)
	if err != nil {
		return err
	}
	defer gr.Close()
	tr := tar.NewReader(gr)
	if dest != "" {
		_, err = makeDir(dest)
		if err != nil {
			return err
		}
	}
	type dirInfo struct {
		Name    string
		ModTime time.Time
	}
	currentDir := dirInfo{}
	for {
		header, err := tr.Next()
		if err != nil {
			if err == io.EOF {
				if currentDir.Name != "" {
					remodifyTime(currentDir.Name, currentDir.ModTime)
				}
				break
			} else {
				return err
			}
		}
		fi := header.FileInfo()
		fileName := filepath.Join(dest, header.Name)
		if !strings.HasPrefix(fileName, currentDir.Name) {
			remodifyTime(currentDir.Name, currentDir.ModTime)
		}
		if fi.IsDir() {
			foldName, err := makeDir(fileName)
			if err != nil {
				return err
			}
			currentDir = dirInfo{
				foldName,
				fi.ModTime(),
			}
			continue
		}
		file, err := createFile(fileName)
		if err != nil {
			return fmt.Errorf("can not create file %v: %v", fileName, err)
		}
		io.Copy(file, tr)
		file.Close()
		remodifyTime(fileName, header.ModTime)
	}
	return nil
}

func remodifyTime(name string, modTime time.Time) {
	if name == "" {
		return
	}
	atime := time.Now()
	os.Chtimes(name, atime, modTime)
}

func makeDir(name string) (string, error) {
	if name != "" {
		_, err := os.Stat(name)
		if err != nil {
			err = os.MkdirAll(name, 0755)
			if err != nil {
				return "", fmt.Errorf("can not make directory: %v", err)
			}
			return name, nil
		}
		return "", nil
	}
	return "", fmt.Errorf("can not make no name directory: %v", name)
}

func createFile(name string) (*os.File, error) {
	dir := path.Dir(name)
	if dir != "" {
		_, err := os.Lstat(dir)
		if err != nil {
			err := os.MkdirAll(dir, 0755)
			if err != nil {
				return nil, err
			}
		}
	}
	return os.Create(name)
}
