package util

import (
	"io/ioutil"
	"os"
)

func writeFile(filepath string, b []byte) error {
	err := ioutil.WriteFile(keypath, b, 0666)
	if err != nil {
		return err
	}
	return nil
}

func existsFile(filename string) bool {
	_, err := os.Stat(filename)
	if err != nil {
		return !os.IsNotExist(err)
	}
	return err == nil
}
