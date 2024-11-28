package agent

import (
	"context"
	"encoding/json"
	"github.com/lvyahui8/ellyn/sdk/common/gocontext"
	"github.com/lvyahui8/ellyn/sdk/common/utils"
)

// NotCollected 用来代表不收集的参数
var NotCollected = &struct {
	_ byte
}{}

// NotCollectedDisplay 不收集参数的展示值
var NotCollectedDisplay = utils.Marshal("[NotCollected]")

// MarshalFailed 序列化失败时的展示值
var MarshalFailed = utils.Marshal("[Marshal failed]")

// EncodeVars 编码参数列表，入参和出参都使用此方法
func EncodeVars(vars []any) *[]any {
	var res []any
	for _, item := range vars {
		if item == NotCollected {
			res = append(res, NotCollectedDisplay)
			continue
		}
		val := item
		switch v := item.(type) {
		case context.Context: // golang context 也支持序列化
			entries := gocontext.GetContextKeyValues(v)
			val = utils.GetCodableMap(entries)
		case *context.Context:
			entries := gocontext.GetContextKeyValues(*v)
			val = utils.GetCodableMap(entries)
		case map[any]any:
			val = utils.GetCodableMap(v)
		default:
		}
		bytes, err := json.Marshal(val)
		if err != nil {
			res = append(res, MarshalFailed)
		} else {
			res = append(res, utils.String.Bytes2string(bytes))
		}
	}
	return &res
}
