package model

import "gorm.io/gorm"

type Department struct {
	gorm.Model
	Name string `json:"name"`
	Code string `json:"code"`

	Rooms  []Room  `json:"rooms" gorm:"foreignKey:DepartmentId"`
	Fields []Field `json:"fields" gorm:"foreignKey:DepartmentId"`
}
