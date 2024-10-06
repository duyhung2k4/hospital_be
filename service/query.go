package service

import (
	"app/config"
	"app/dto/request"

	"gorm.io/gorm"
)

type queryService[T any] struct {
	psql *gorm.DB
}

type QueryService[T any] interface {
	First(payload request.FirstPayload, agrs ...interface{}) (*T, error)
	Find(payload request.FindPayload, agrs ...interface{}) ([]T, error)
	Create(data T) (*T, error)
	Update(data T, preload []string, omit map[string][]string, condition string, args ...interface{}) (*T, error)
	Delete(condition string, args ...interface{}) error
}

func (s *queryService[T]) First(payload request.FirstPayload, agrs ...interface{}) (*T, error) {
	var item *T
	var personOmit []string

	for key, omitChild := range payload.Omit {
		if len(omitChild) == 0 {
			personOmit = append(personOmit, key)
		}
	}

	query := s.psql.Where(payload.Condition, agrs...).Omit(personOmit...)

	for _, p := range payload.Preload {
		query.Preload(p, func(tx *gorm.DB) *gorm.DB {
			return tx.Omit(payload.Omit[p]...)
		})
	}

	err := query.First(&item).Error
	if err != nil {
		return nil, err
	}

	return item, nil
}

func (s *queryService[T]) Find(payload request.FindPayload, agrs ...interface{}) ([]T, error) {
	var list []T
	var personOmit []string

	for key, omitChild := range payload.Omit {
		if len(omitChild) == 0 {
			personOmit = append(personOmit, key)
		}
	}

	query := s.psql.Where(payload.Condition, agrs...).Omit(personOmit...)

	for _, p := range payload.Preload {
		query.Preload(p, func(tx *gorm.DB) *gorm.DB {
			return tx.Omit(payload.Omit[p]...)
		})
	}

	query = query.Order(payload.Order)

	err := query.Find(&list).Error
	if err != nil {
		return nil, err
	}

	return list, nil
}

func (s *queryService[T]) Create(data T) (*T, error) {
	newData := data
	if err := s.psql.Create(&newData).Error; err != nil {
		return nil, err
	}
	return &newData, nil
}

func (s *queryService[T]) Update(data T, preload []string, omit map[string][]string, condition string, args ...interface{}) (*T, error) {
	newData := data

	query := s.psql.Where(condition, args...).Updates(&newData)
	for _, p := range preload {
		query.Preload(p, func(tx *gorm.DB) *gorm.DB {
			return tx.Omit(omit[p]...)
		})
	}

	if err := query.First(&newData).Error; err != nil {
		return nil, err
	}

	return &newData, nil
}

func (s *queryService[T]) Delete(condition string, args ...interface{}) error {
	var del T
	if err := s.psql.Where(condition, args...).Unscoped().Delete(&del).Error; err != nil {
		return err
	}
	return nil
}

func NewQueryService[T any]() QueryService[T] {
	return &queryService[T]{
		psql: config.GetPsql(),
	}
}
