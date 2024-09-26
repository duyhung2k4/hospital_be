package model

import "gorm.io/gorm"

type Field struct {
	gorm.Model
	Lable         string `json:"lable"`
	Name          string `json:"name"`
	Type          string `json:"type"` // int - float - text - select
	DefaultValues string `json:"default_values"`
	Value         string `json:"value"`

	DepartmentId uint        `json:"departmentId"`
	Department   *Department `json:"department" gorm:"foreignKey:DepartmentId; constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}
