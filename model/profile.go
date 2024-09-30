package model

import "gorm.io/gorm"

type Profile struct {
	gorm.Model
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Phone     string `json:"phone"`
	Email     string `json:"email"`
	Address   string `json:"address"`
	Gender    string `json:"gender"`
	Username  string `json:"username" gorm:"unique"`
	Password  string `json:"password"`
	Role      string `json:"role"` // admin - user - clin - spec - room
	RoomId    *uint  `json:"roomId"`

	Room  *Room  `json:"room" gorm:"foreignKey:RoomId; constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Faces []Face `json:"faces" gorm:"foreignKey:ProfileId"`
}
