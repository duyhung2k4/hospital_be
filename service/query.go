package service

import (
	"app/config"

	"gorm.io/gorm"
)

type queryService[T any] struct {
	psql *gorm.DB
}

type QueryService[T any] interface {
	First(preload []string, omit map[string][]string, condition string, agrs ...interface{}) (*T, error)
	Find(preload []string, omit map[string][]string, condition string, agrs ...interface{}) ([]T, error)
	Create(data T) (*T, error)
	Update(data T, preload []string, omit map[string][]string, condition string, args ...interface{}) (*T, error)
	Delete(condition string, args ...interface{}) error
}

func (s *queryService[T]) First(preload []string, omit map[string][]string, condition string, agrs ...interface{}) (*T, error) {
	var item *T

	query := s.psql.Where(condition, agrs...)

	for _, p := range preload {
		query.Preload(p, func(tx *gorm.DB) *gorm.DB {
			return tx.Omit(omit[p]...)
		})
	}

	err := query.First(&item).Error
	if err != nil {
		return nil, err
	}

	return item, nil
}

func (s *queryService[T]) Find(preload []string, omit map[string][]string, condition string, agrs ...interface{}) ([]T, error) {
	var list []T

	query := s.psql.Where(condition, agrs...)

	for _, p := range preload {
		query.Preload(p, func(tx *gorm.DB) *gorm.DB {
			return tx.Omit(omit[p]...)
		})
	}

	err := query.Order("id asc").Find(&list).Error
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
	if err := s.psql.Where(condition, args...).Delete(&del).Error; err != nil {
		return err
	}
	return nil
}

func NewQueryService[T any]() QueryService[T] {
	return &queryService[T]{
		psql: config.GetPsql(),
	}
}
