package model

import "gorm.io/gorm"

type Schedule struct {
	gorm.Model
	Code string `json:"code"`

	ClinId    *uint    `json:"clinId"`
	ProfileId uint     `json:"profileId"`
	Profile   *Profile `json:"profile" gorm:"foreignKey:ProfileId; constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Clin      *Profile `json:"clin" gorm:"foreignKey:ClinId; constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Steps     []Step   `json:"steps" gorm:"foreignKey:ScheduleId"`
}
