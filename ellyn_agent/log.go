package ellyn_agent

import "fmt"

var log = &logger{}

type logger struct {
}

func (l *logger) Error(format string, args ...any) {
	fmt.Printf("[Error]"+format, args...)
}

func (l *logger) Info(format string, args ...any) {
	fmt.Printf("[Info]"+format, args...)
}
