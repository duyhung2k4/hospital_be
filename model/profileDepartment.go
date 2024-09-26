package model

import "gorm.io/gorm"

type ProfileDepartment struct {
	gorm.Model
	ProfileId    uint `json:"profileId"`
	DepartmentId uint `json:"departmentId"`

	Profile    *Profile    `json:"Profile" gorm:"foreignKey:ProfileId; constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Department *Department `json:"Department" gorm:"foreignKey:DepartmentId; constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}
