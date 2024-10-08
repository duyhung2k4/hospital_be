package model

import (
	"fmt"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Profile struct {
	gorm.Model
	Avatar    pq.Float64Array `json:"avatar" gorm:"type:float8[]"`
	FirstName string          `json:"firstName"`
	LastName  string          `json:"lastName"`
	Phone     string          `json:"phone"`
	Email     string          `json:"email"`
	Address   string          `json:"address"`
	Gender    string          `json:"gender"`
	Username  string          `json:"username" gorm:"unique"`
	Password  string          `json:"password"`
	Role      string          `json:"role"` // admin - user - clin - spec - room
	RoomId    *uint           `json:"roomId"`
	Active    bool            `json:"active" gorm:"default:false"`

	Room  *Room  `json:"room" gorm:"foreignKey:RoomId; constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Faces []Face `json:"faces" gorm:"foreignKey:ProfileId"`
}

type ROLE string

const (
	ADMIN ROLE = "admin"
	USER  ROLE = "user"
	CLIN  ROLE = "clin"
	SPEC  ROLE = "spec"
	ROOM  ROLE = "room"
)

func (p *Profile) BeforeCreate(tx *gorm.DB) (err error) {
	var existingProfile Profile
	if err = tx.Unscoped().Where("username = ? AND deleted_at IS NULL", p.Username).First(&existingProfile).Error; err == nil {
		return fmt.Errorf("username %s đã tồn tại", p.Username)
	}
	return nil
}
