package benchmark

import (
	"os"
	"path/filepath"
)

var devNull *os.File
var tmpFile *os.File

func init() {
	f, err := os.Open(os.DevNull)
	if err != nil {
		panic(err)
	}
	devNull = f
	f, err = os.OpenFile(filepath.Join(os.TempDir(), "bench.log"), os.O_CREATE|os.O_APPEND|os.O_RDWR, os.ModePerm)
	if err != nil {
		panic(err)
	}
	tmpFile = f
}

func Write2DevNull(content string) {
	_, _ = devNull.Write([]byte(content))
}

func Write2TmpFile(content string) {
	_, _ = tmpFile.Write([]byte(content))
}
