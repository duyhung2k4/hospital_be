package model

import "gorm.io/gorm"

type LogCheck struct {
	gorm.Model
	ProfileId uint    `json:"profileId"`
	Accuracy  float64 `json:"accuracy"`
	Url       string  `json:"url"`

	Profile *Profile `json:"Profile" gorm:"foreignKey:ProfileId; constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}
