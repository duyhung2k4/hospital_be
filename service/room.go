package service

import (
	"app/config"
	"app/model"
	"errors"

	"gorm.io/gorm"
)

type roomService struct {
	psql *gorm.DB
}

type RoomService interface {
	RoomMinStepWaiting(departmentId uint) (*model.Room, error)
}

func (s *roomService) RoomMinStepWaiting(departmentId uint) (*model.Room, error) {
	var rooms []model.Room
	var minRoom *model.Room

	if err := s.psql.Model(&model.Room{}).
		Where("department_id = ?", departmentId).
		Preload("Steps").
		Find(&rooms).Error; err != nil {
		return nil, err
	}

	if len(rooms) == 0 {
		return nil, errors.New("room not found(2)")
	}

	if len(rooms) == 1 {
		return &rooms[0], nil
	}

	minRoom = &rooms[0]
	for i := 1; i < len(rooms); i++ {
		if len(minRoom.Steps) > len(rooms[i].Steps) {
			minRoom = &rooms[i]
		}
	}

	return minRoom, nil
}

func NewRoomService() RoomService {
	return &roomService{
		psql: config.GetPsql(),
	}
}
