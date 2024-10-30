package utils

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGzipUtils(t *testing.T) {
	str := "lvyahui"
	compressed := Gzip.Compress([]byte(str))
	bytes := Gzip.UnCompress(compressed)
	t.Log(string(bytes))
	require.Equal(t, str, string(bytes))
}
