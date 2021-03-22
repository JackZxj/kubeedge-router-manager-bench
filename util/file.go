package util

import (
	"io/ioutil"
	"os"
)

// FileExists checks the target file
func FileExists(path string) error {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if os.IsNotExist(err) {
		return err
	}
	return nil
}

// ReadFile returns strings of file
func ReadFile(filename string) ([]byte, error) {
	err := FileExists(filename)
	if err != nil {
		return nil, err
	}
	f, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return f, nil
}
