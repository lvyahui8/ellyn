package log

import "fmt"

// Infof 打印日志，默认输出到stdout
func Infof(format string, args ...any) {
	fmt.Printf(format+"\n", args...)
}
