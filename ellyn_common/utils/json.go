package utils

import (
	"encoding/json"
	"github.com/lvyahui8/ellyn/ellyn_common/asserts"
)

func Marshal(v any) string {
	bytes, err := json.Marshal(v)
	asserts.IsNil(err)
	return string(bytes)
}
