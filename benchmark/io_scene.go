package benchmark

import (
	"fmt"
	"io/ioutil"
	"net"
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
	_, err := devNull.Write([]byte(content))
	if err != nil {
		_ = fmt.Errorf("write failed. %v", err)
	}
}

func Write2TmpFile(content string) {
	_, err := tmpFile.Write([]byte(content))
	if err != nil {
		_ = fmt.Errorf("write failed. %v", err)
	}
}

func NetworkReadWrite(content string) {
	r, w := net.Pipe()
	go func() {
		_, err := w.Write([]byte(content))
		if err != nil {
			_ = fmt.Errorf("write failed. %v", err)
			return
		}
		_ = w.Close()
	}()
	data, err := ioutil.ReadAll(r)
	if err != nil {
		_ = fmt.Errorf("write failed. %v", err)
		return
	}
	_ = fmt.Sprintf("data:%s", string(data))
}
