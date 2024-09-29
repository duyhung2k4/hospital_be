package model

import "gorm.io/gorm"

type Room struct {
	gorm.Model
	Name string `json:"name"`
	Code string `json:"code"`

	DepartmentId uint        `json:"departmentId"`
	Department   *Department `json:"department" gorm:"foreignKey: DepartmentId; constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Steps        []Step      `json:"steps" gorm:"foreignKey: RoomId"`
}
