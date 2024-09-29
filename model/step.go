package model

import "gorm.io/gorm"

type Step struct {
	gorm.Model
	Index  int         `json:"index"`
	Result string      `json:"result"`
	Status STEP_STATUS `json:"status"` // pending - wating - examining - done

	ScheduleId   uint        `json:"scheduleId"`
	DepartmentId uint        `json:"departmentId"`
	RoomId       *uint       `json:"roomId"`
	SpecId       *uint       `json:"specId"`
	Schedule     *Schedule   `json:"schedule" gorm:"foreignKey:ScheduleId; constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Department   *Department `json:"department" gorm:"foreignKey:DepartmentId; constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Room         *Room       `json:"room" gorm:"foreignKey:RoomId; constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Spec         *Profile    `json:"spec" gorm:"foreignKey:SpecId; constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

type STEP_STATUS string

const (
	ST_EXAMINING STEP_STATUS = "examining"
	ST_PENDING   STEP_STATUS = "pending"
	ST_WATING    STEP_STATUS = "waiting"
	ST_DONE      STEP_STATUS = "done"
)
