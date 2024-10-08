package utils

import "github.com/lvyahui8/ellyn/ellyn_common/definitions"

func Min[T definitions.Number](a, b T) T {
	if a > b {
		return b
	} else {
		return a
	}
}

func Max[T definitions.Number](a, b T) T {
	if a > b {
		return a
	} else {
		return b
	}
}
