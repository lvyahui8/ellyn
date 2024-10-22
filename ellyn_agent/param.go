package ellyn_agent

import (
	"context"
	"encoding/json"
	"github.com/lvyahui8/ellyn/ellyn_common/gocontext"
	"github.com/lvyahui8/ellyn/ellyn_common/utils"
)

var NotCollected = &struct {
	_ byte
}{}

var NotCollectedDisplay = utils.Marshal("[NotCollected]")

var MarshalFailed = utils.Marshal("[Marshal failed]")

func EncodeVars(vars []any) (res []any) {
	for _, item := range vars {
		if item == NotCollected {
			res = append(res, NotCollectedDisplay)
			continue
		}
		val := item
		switch v := item.(type) {
		case context.Context:
			entries := make(map[string]interface{})
			rawEntries := gocontext.GetContextKeyValues(v)
			for k, v := range rawEntries {
				if sKey, ok := k.(string); ok {
					entries[sKey] = v
				}
			}
			val = entries
		default:
		}
		bytes, err := json.Marshal(val)
		if err != nil {
			res = append(res, MarshalFailed)
		} else {
			res = append(res, utils.String.Bytes2string(bytes))
		}
	}
	return
}
