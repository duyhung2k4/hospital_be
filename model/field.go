package model

import "gorm.io/gorm"

type Field struct {
	gorm.Model
	Lable         string `json:"lable"`
	Name          string `json:"name" gorm:"uniqueIndex:name_department_idx"`
	Type          string `json:"type"` // int - float - text - select
	DefaultValues string `json:"defaultValues"`

	DepartmentId uint        `json:"departmentId" gorm:"uniqueIndex:name_department_idx"`
	Department   *Department `json:"department" gorm:"foreignKey:DepartmentId; constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}
