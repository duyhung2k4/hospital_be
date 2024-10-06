package model

import (
	"time"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Schedule struct {
	gorm.Model

	RoomId      *uint           `json:"roomId"`
	ClinId      *uint           `json:"clinId"`
	Code        string          `json:"code" gorm:"unique"`
	Name        string          `json:"name"`
	Dob         time.Time       `json:"dob"`
	Address     string          `json:"address"`
	Gender      string          `json:"gender"`
	Phone       string          `json:"phone"`
	Description string          `json:"description"`
	Result      string          `json:"result"`
	Avatar      pq.Float64Array `json:"avatar" gorm:"type:float8[]"` // `float8[]` cho PostgreSQL array
	Status      SCHEDULE_STATUS `json:"status"`                      // pending - examining - transited - finished

	Room *Room    `json:"room" gorm:"foreignKey:RoomId; constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Clin *Profile `json:"clin" gorm:"foreignKey:ClinId; constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`

	Steps []Step `json:"steps" gorm:"foreignKey:ScheduleId"`
}

type SCHEDULE_STATUS string

const (
	S_PENDING   SCHEDULE_STATUS = "pending"
	S_EXAMINING SCHEDULE_STATUS = "examining"
	S_FINISHED  SCHEDULE_STATUS = "finished"
	S_TRANSITED SCHEDULE_STATUS = "transited"
)
