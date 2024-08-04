package utils

import "unsafe"

var String = &stringUtils{}

type stringUtils struct {
}

func (stringUtils) String2bytes(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(&s))
}
func (stringUtils) Bytes2string(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
