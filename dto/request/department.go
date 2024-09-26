package request

import (
	"app/constant"
	"app/model"
)

type CreateDepartmentReq struct {
	Name string `json:"name"`
	Code string `json:"code"`
}

type CreateFieldReq struct {
	Lable         string `json:"lable"`
	Name          string `json:"name"`
	Type          string `json:"type"` // int - float - text - select
	DefaultValues string `json:"defaultValues"`
	Value         string `json:"value"`

	DepartmentId uint `json:"departmentId"`
}

type QueryDepartmentReq struct {
	Data      model.Department `json:"data"`
	Condition string           `json:"condition"`
	Args      []interface{}    `json:"args"`
	Preload   []string         `json:"preload"`
	Method    constant.METHOD  `json:"method"`
}
