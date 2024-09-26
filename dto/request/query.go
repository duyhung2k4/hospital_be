package request

import (
	"app/constant"
)

type QueryReq[T any] struct {
	Data      T                   `json:"data"`
	Condition string              `json:"condition"`
	Args      []interface{}       `json:"args"`
	Preload   []string            `json:"preload"`
	Omit      map[string][]string `json:"omit"`
	Method    constant.METHOD     `json:"method"`
}
