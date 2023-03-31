package app

import (
	"encoding/json"
	"fmt"

	"github.com/dumacp/go-doors/api"
)

func ParseRequest(v []byte) interface{} {
	req := new(api.DoorsStateRequest)
	err := json.Unmarshal(v, req)
	if err != nil {
		return err
	}
	if len(req.Topic) <= 0 {
		return fmt.Errorf("topin len error")
	}
	return req
}
