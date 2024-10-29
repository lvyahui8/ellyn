package benchmark

import "os"

var devNull *os.File

func init() {
	f, err := os.Open(os.DevNull)
	if err != nil {
		panic(err)
	}
	devNull = f
}

func Write2DevNull(content string) {
	_, _ = devNull.Write([]byte(content))
}
