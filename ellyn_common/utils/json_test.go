package utils

import (
	"github.com/stretchr/testify/require"
	"strconv"
	"testing"
)

type Level int

func (l Level) MarshalText() (text []byte, err error) {
	return []byte(strconv.Itoa(int(l) * 100)), nil
}

func TestCodableMap(t *testing.T) {
	m := make(map[any]any)
	m[1] = 1
	m["name"] = "ellyn"
	m[Level(1)] = "xx" // key  100
	m[1.11] = 1.11
	res := GetCodableMap(m)
	require.NotEmpty(t, res)
	require.True(t, len(m) != len(res))
	_, exist := res["100"]
	require.True(t, exist)
	t.Log(res)
}
