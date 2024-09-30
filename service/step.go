package service

import (
	"app/config"
	"app/model"
	"errors"

	"gorm.io/gorm"
)

type stepService struct {
	psql        *gorm.DB
	roomService RoomService
}

type StepService interface {
	NextStep(scheduleId uint) error
}

func (s *stepService) NextStep(scheduleId uint) error {
	var schedule *model.Schedule
	var step model.Step

	tx := s.psql.Begin()

	if err := tx.Model(&model.Schedule{}).
		Where("id = ?", scheduleId).
		First(&schedule).
		Error; err != nil {
		return err
	}

	if schedule == nil {
		return errors.New("schedule not found")
	}

	if err := tx.Model(&model.Step{}).
		Where("status = ? AND schedule_id = ?", model.ST_PENDING, schedule.ID).
		Limit(1).
		First(&step).
		Error; err != nil {
		return err
	}
	if step.ID == 0 {
		return errors.New("step not found")
	}

	room, err := s.roomService.RoomMinStepWaiting(step.DepartmentId)
	if err != nil {
		return err
	}
	if room == nil {
		return errors.New("room not found(1)")
	}

	if err := s.psql.Model(&model.Step{}).
		Where("id = ?", step.ID).
		Updates(&model.Step{
			RoomId: &room.ID,
			Status: model.ST_WATING,
		}).
		Error; err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit().Error
	if err != nil {
		return err
	}

	return nil
}

func NewStepService() StepService {
	return &stepService{
		psql:        config.GetPsql(),
		roomService: NewRoomService(),
	}
}