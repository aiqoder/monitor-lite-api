package fileutils

import (
	"errors"
	"github.com/duke-git/lancet/v2/fileutil"
	"os"
	"path/filepath"
)

// defVal 创建文件输入的内容
func CreateFile(path string, defVal []byte) error {
	if fileutil.IsExist(path) {
		return nil
	}

	dir, _ := filepath.Split(path)
	if !fileutil.IsExist(dir) {
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return err
		}
	}

	if !fileutil.CreateFile(path) {
		return errors.New("create file failed")
	}

	if defVal != nil {
		err := fileutil.WriteBytesToFile(path, defVal)

		if err != nil {
			return err
		}
	}

	return nil
}
