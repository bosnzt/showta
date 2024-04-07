package util

import (
	"os"
)

func PathExist(pathStr string) (bool, error) {
	_, err := os.Stat(pathStr)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
