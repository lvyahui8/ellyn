package asserts

import (
	"fmt"
	"runtime/debug"
	"strings"
)

type AssertError struct {
	msg      string
	fileLine string
}

func NewAssertError(msg string) (ae AssertError) {
	ae.msg = msg
	stack := string(debug.Stack())
	lines := strings.Split(stack, "\n")
	visited := false
	for i := 0; i < len(lines); i++ {
		if !strings.Contains(lines[i], "/asserts.") {
			if visited {
				ae.fileLine = strings.Join(lines[i-1:i+1], "\n")
				break
			}
		} else {
			visited = true
		}
	}
	return ae
}

func (ae AssertError) Error() string {
	return fmt.Sprintf("msg:%s,fileLine:%s", ae.msg, ae.fileLine)
}

func IsNil(a interface{}) {
	if a != nil {
		panic(NewAssertError(fmt.Sprintf("must be nil. but got %+v", a)))
	}
}

func NotNil(a interface{}) {
	if a == nil {
		panic(NewAssertError("must be not nil. but got nil"))
	}
}
