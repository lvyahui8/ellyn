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
		default:
			bytes, err := json.Marshal(v)
			if err != nil {
				res = append(res, MarshalFailed)
			} else {
				res = append(res, utils.String.Bytes2string(bytes))
			}
		}
	}
	return
}
