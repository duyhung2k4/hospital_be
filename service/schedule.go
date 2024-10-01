package service

import (
	"app/config"
	"app/dto/request"
	"app/model"
	"errors"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type scheduleService struct {
	psql        *gorm.DB
	stepService StepService
}

type ScheduleService interface {
	PullSchedule() (*model.Schedule, error)
	Transit(payload request.TransitReq) error
}

func (s *scheduleService) PullSchedule() (*model.Schedule, error) {
	var schedule model.Schedule

	tx := s.psql.Begin()

	// Locking
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("status = ?", model.S_PENDING).
		First(&schedule).
		Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Model(&schedule).
		Update("status", model.S_EXAMINING).
		Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return &schedule, nil
}

func (s *scheduleService) Transit(payload request.TransitReq) error {
	var steps []model.Step

	tx := s.psql.Begin()

	if result := tx.Model(&model.Schedule{}).
		Where("id = ? AND status = ?", payload.ScheduleId, model.S_EXAMINING).
		Updates(&model.Schedule{
			Description: payload.Description,
			Status:      model.S_TRANSITED,
		}); result.RowsAffected == 0 || result.Error != nil {
		return errors.New("schedule not found")
	}

	if len(payload.DepartmentIds) == 0 {
		return nil
	}

	for i, d_id := range payload.DepartmentIds {
		steps = append(steps, model.Step{
			DepartmentId: uint(d_id),
			Index:        i + 1,
			Status:       model.ST_PENDING,
			ScheduleId:   payload.ScheduleId,
		})
	}

	if err := tx.Model(&model.Step{}).Create(&steps).Error; err != nil {
		return err
	}

	err := tx.Commit().Error
	if err != nil {
		return err
	}

	if err := s.stepService.NextStep(payload.ScheduleId); err != nil {
		return err
	}

	return nil
}

func NewScheduleService() ScheduleService {
	return &scheduleService{
		psql:        config.GetPsql(),
		stepService: NewStepService(),
	}
}
