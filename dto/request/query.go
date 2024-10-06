package request

import (
	"app/constant"
)

type QueryReq[T any] struct {
	Data      T                   `json:"data"`
	Args      []interface{}       `json:"args"`
	Condition string              `json:"condition"`
	Preload   []string            `json:"preload"`
	Omit      map[string][]string `json:"omit"`
	Method    constant.METHOD     `json:"method"`
	Order     string              `json:"order"`
}

type FindPayload struct {
	Condition string
	Preload   []string
	Omit      map[string][]string
	Order     string
}

type FirstPayload struct {
	Condition string
	Preload   []string
	Omit      map[string][]string
}
