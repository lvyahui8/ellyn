package assert

import "fmt"

func IsNil(a interface{}) {
	if a != nil {
		panic(fmt.Sprintf("must be nil. but got %+v", a))
	}
}

func AssertNotNil(a interface{}) {
	if a == nil {
		panic("must be not nil. but got nil")
	}
}
