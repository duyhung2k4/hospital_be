package service

import (
	"app/config"
	"app/dto/request"
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
	SaveStep(payload request.SaveStepReq) error
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
		Order("index ASC").
		Limit(1).
		First(&step).
		Error; err != nil && err != gorm.ErrRecordNotFound {
		return err
	}

	if step.ID == 0 {
		err := tx.Model(&model.Schedule{}).
			Where("id = ?", scheduleId).
			Updates(&model.Schedule{Status: model.S_FINISHED}).Error

		if err != nil {
			return err
		}

		err = tx.Commit().Error

		if err != nil {
			return err
		}

		return nil
	}

	room, err := s.roomService.RoomMinStepWaiting(step.DepartmentId)
	if err != nil {
		return err
	}
	if room == nil {
		return errors.New("room not found(1)")
	}

	if err := tx.Model(&model.Step{}).
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

func (s *stepService) SaveStep(payload request.SaveStepReq) error {
	if err := s.psql.
		Model(&model.Step{}).
		Where("schedule_id = ? AND room_id = ?", payload.ScheduleId, payload.RoomId).
		Updates(&model.Step{
			Status: model.ST_DONE,
			Result: payload.Result,
		}).Error; err != nil {
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
