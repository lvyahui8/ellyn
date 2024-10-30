package utils

import (
	"bytes"
	"compress/gzip"
	"github.com/lvyahui8/ellyn/sdk/common/asserts"
	"io/ioutil"
)

var Gzip = &gzipUtils{}

type gzipUtils struct {
}

func (gzipUtils) Compress(content []byte) []byte {
	var b bytes.Buffer
	gz := gzip.NewWriter(&b)
	_, err := gz.Write(content)
	asserts.IsNil(err)
	err = gz.Close()
	asserts.IsNil(err)
	return b.Bytes()
}

func (gzipUtils) UnCompress(compressed []byte) []byte {
	reader := bytes.NewReader(compressed)
	gz, err := gzip.NewReader(reader)
	asserts.IsNil(err)
	out, err := ioutil.ReadAll(gz)
	asserts.IsNil(err)
	return out
}
