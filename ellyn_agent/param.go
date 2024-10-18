package ellyn_agent

import (
	"encoding/json"
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
		switch v := item.(type) {
		case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, uintptr, float32, float64:
			res = append(res, v)
		default:
			bytes, err := json.Marshal(item)
			if err != nil {
				res = append(res, MarshalFailed)
			} else {
				res = append(res, utils.String.Bytes2string(bytes))
			}
		}
	}
	return
}
