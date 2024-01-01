package utils

import "os"

var OS = &osUtils{}

type osUtils struct {
}

func (osUtils) GetWorkDir() string {
	dir, _ := os.Getwd()
	return dir
}
