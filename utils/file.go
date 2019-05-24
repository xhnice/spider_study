package utils

import (
	"os"
	"path/filepath"
	"sync_study/logger"
)

// FileExists -
func FileExists(file string) bool {
	_, err := os.Stat(file)
	return err == nil || os.IsExist(err)
}

// OpenFile -
func OpenFile(file string) (*os.File, error) {
	path := filepath.Dir(file)
	if err := os.MkdirAll(path, 0666); err != nil {
	}

	return os.OpenFile(file, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
}

// WriterFile 写入内容到文件中
func WriterFile(filename string, data []byte, isdelete bool) {
	if isdelete {
		// 文件已存在 删除文件
		if FileExists(filename) {
			os.Remove(filename)
		}
	}

	// 写入到文件
	file, err := OpenFile(filename)
	if err != nil {
		logger.Errorf("cannot create file: %s error: %s", filename, err)
		return
	}
	defer file.Close()
	file.Write(data)
}
