package util

import (
	"io/ioutil"
	"os"
)

// WriteFile to path whit content
func WriteFile(path string, content string) error {
	fileObj, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer fileObj.Close()
	contents := []byte(content)
	if _, err := fileObj.Write(contents); err != nil {
		return err
	}
	return nil
}

// WriteFileAppend for file to append
func WriteFileAppend(path string, content string) error {
	fileObj, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	defer fileObj.Close()
	contents := []byte(content)
	if _, err := fileObj.Write(contents); err != nil {
		return err
	}
	return nil
}

// ReadFile from path
func ReadFile(path string) (string, error) {
	fi, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer fi.Close()
	fd, err := ioutil.ReadAll(fi)
	return string(fd), nil
}

// RemoveFile from path
func RemoveFile(path string) error {
	if err := os.RemoveAll(path); err != nil {
		return err
	}
	return nil
}

// MakeDirAll : create directory from path
func MakeDirAll(path string) error {
	if err := os.MkdirAll(path, 0755); err != nil {
		return err
	}
	return nil
}

// ReadDir : list for path
func ReadDir(path string) ([]string, error) {
	if dirList, err := ioutil.ReadDir(path); err != nil {
		return nil, err
	} else {
		var fileName = make([]string, len(dirList))
		for i, v := range dirList {
			if v.IsDir() {
				fileName[i] = "dir:" + v.Name()
			} else {
				fileName[i] = "file:" + v.Name()
			}

		}
		return fileName, nil
	}
}
