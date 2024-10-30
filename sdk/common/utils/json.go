package utils

import (
	"encoding"
	"encoding/json"
	"github.com/lvyahui8/ellyn/sdk/common/asserts"
	"strconv"
)

func Marshal(v any) string {
	bytes, err := json.Marshal(v)
	asserts.IsNil(err)
	return String.Bytes2string(bytes)
}

func GetCodableMap(m map[any]any) map[string]any {
	// /go/go1.18/src/encoding/json/encode.go:820
	// func newMapEncoder(t reflect.Type) encoderFunc
	// go 官方的json序列化map只支持三种类型
	// - string
	// - int(int/int8-64/uint/uint8-64/uintptr)
	// - *encoding.TextMarshaler impl

	// key 转换成string 逻辑参考
	// func (me mapEncoder) encode(e *encodeState, v reflect.Value, opts encOpts)
	// func (w *reflectWithString) resolve() error
	resMap := make(map[string]any)
	for k, v := range m {
		var strKey string
		if txtEncoder, ok := k.(encoding.TextMarshaler); ok {
			text, err := txtEncoder.MarshalText()
			asserts.IsNil(err)
			strKey = string(text)
		} else {
			switch n := k.(type) {
			case string:
				strKey = n
			case int:
				strKey = strconv.FormatInt(int64(n), 10)
			case int8:
				strKey = strconv.FormatInt(int64(n), 10)
			case int16:
				strKey = strconv.FormatInt(int64(n), 10)
			case int32:
				strKey = strconv.FormatInt(int64(n), 10)
			case int64:
				strKey = strconv.FormatInt(int64(n), 10)
			case uint:
				strKey = strconv.FormatUint(uint64(n), 10)
			case uint8:
				strKey = strconv.FormatUint(uint64(n), 10)
			case uint16:
				strKey = strconv.FormatUint(uint64(n), 10)
			case uint32:
				strKey = strconv.FormatUint(uint64(n), 10)
			case uint64:
				strKey = strconv.FormatUint(uint64(n), 10)
			case uintptr:
				strKey = strconv.FormatUint(uint64(n), 10)
			default:
				// unsupported key type
				continue
			}
		}
		resMap[strKey] = v
	}
	return resMap
}
