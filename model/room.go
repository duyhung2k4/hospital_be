package model

import "gorm.io/gorm"

type Room struct {
	gorm.Model
	Name         string `json:"name"`
	Code         string `json:"code" gorm:"unique"`
	DepartmentId uint   `json:"departmentId"`

	Department *Department `json:"department" gorm:"foreignKey: DepartmentId; constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Profile    *Profile    `json:"profile" gorm:"foreignKey:RoomId"`
	Steps      []Step      `json:"steps" gorm:"foreignKey: RoomId"`
}
