package utils

import (
	"io/ioutil"
	"os"
)

func ReadFile(filePath string) (string, error) {
	bytes, err := ioutil.ReadFile(filePath)
	return string(bytes), err
}

func OverwriteToFile(filePath string, content string) error {
	f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write([]byte(content))
	return err
}
