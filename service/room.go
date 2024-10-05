package service

import (
	"app/config"
	"app/dto/request"
	"app/model"
	"app/utils"
	"errors"
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type roomService struct {
	psql      *gorm.DB
	authUtils utils.AuthUtils
	smtpUtils utils.SmtpUtils
}

type RoomService interface {
	RoomMinStepWaiting(departmentId uint) (*model.Room, error)
	PullStep(roomId uint) (*model.Step, error)
	CallStep(roomId uint) (*model.Step, error)
	AddAccount(payload request.AddAccountRoomReq) (*model.Room, error)
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

func (s *roomService) PullStep(roomId uint) (*model.Step, error) {
	var step *model.Step

	tx := s.psql.Begin()

	if err := tx.Model(&model.Step{}).
		Where("room_id = ? AND status = ?", roomId, model.ST_EXAMINING).
		First(&step).
		Error; err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	if step.ID != 0 {
		return nil, errors.New("exist step examining")
	}

	if err := tx.Model(&model.Step{}).
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("room_id = ? AND status = ?", roomId, model.ST_WATING).
		First(&step).Error; err != nil {
		return nil, err
	}

	if step.ID == 0 {
		return nil, errors.New("step not found")
	}

	if err := tx.Model(&model.Step{}).
		Where("id = ?", step.ID).
		Updates(&model.Step{Status: model.ST_EXAMINING}).
		Error; err != nil {
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return step, nil
}

func (s *roomService) CallStep(roomId uint) (*model.Step, error) {
	var step *model.Step

	if err := s.psql.
		Model(&model.Step{}).
		Where("room_id = ? AND status = ?", roomId, model.ST_EXAMINING).
		First(&step).Error; err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	if step.ID == 0 {
		return nil, nil
	}

	return step, nil
}

func (s *roomService) AddAccount(payload request.AddAccountRoomReq) (*model.Room, error) {
	var room *model.Room
	var newProfile *model.Profile
	var profile *model.Profile
	var err error

	tx := s.psql.Begin()
	if err = tx.Model(&model.Room{}).Where("id = ?", payload.RoomId).First(&room).Error; err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	passwordHash, err := s.authUtils.HashPassword(payload.Password)
	if err != nil {
		return nil, err
	}

	newProfile = &model.Profile{
		Username: room.Code,
		Password: passwordHash,
		RoomId:   &room.ID,
		Role:     string(model.ROOM),
	}

	if err = tx.Model(&model.Profile{}).
		Where("room_id = ?", payload.RoomId).
		First(&profile).
		Error; err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	if profile.ID == 0 {
		err = tx.Model(&model.Profile{}).Create(&newProfile).Error
	} else {
		err = tx.Model(&model.Profile{}).Where("room_id = ?", payload.RoomId).Updates(&newProfile).Error
	}

	if err != nil {
		tx.Rollback()
		return nil, err
	}

	data := fmt.Sprintf("%s / %s", room.Code, payload.Password)
	if err = s.smtpUtils.SendEmail(data, payload.EmailAccept); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err = tx.Commit().Error; err != nil {
		return nil, err
	}

	room.Profile = newProfile

	return room, nil
}

func NewRoomService() RoomService {
	return &roomService{
		psql:      config.GetPsql(),
		authUtils: utils.NewAuthUtils(),
		smtpUtils: utils.NewSmtpUtils(),
	}
}
