package model

import "gorm.io/gorm"

type Profile struct {
	gorm.Model `json:"gorm"`
	FirstName  string `json:"firstName"`
	LastName   string `json:"lastName"`
	Phone      string `json:"phone"`
	Email      string `json:"email"`
	Address    string `json:"address"`
	Gender     string `json:"gender"`
	Password   string `json:"password"`
	Role       string `json:"role"` // admin - user - clin - spec

	Faces []Face `json:"faces" gorm:"foreignKey:ProfileId"`
}
