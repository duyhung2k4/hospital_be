package model

import (
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Field struct {
	gorm.Model
	Label         string         `json:"label"`
	Placeholder   string         `json:"placeholder"`
	Name          string         `json:"name" gorm:"uniqueIndex:name_department_idx"`
	Size          int            `json:"size"`                             // width: 1 - 12
	Type          string         `json:"type"`                             // int - float - text - area - select
	DefaultValues pq.StringArray `json:"defaultValues" gorm:"type:text[]"` // string []{ label: string, value: string }

	DepartmentId uint        `json:"departmentId" gorm:"uniqueIndex:name_department_idx"`
	Department   *Department `json:"department" gorm:"foreignKey:DepartmentId; constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}
