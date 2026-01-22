package services

import (
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type Service[T server] struct {
	db *gorm.DB
}

type server interface {
	Type() string
}

type Preloadable interface {
	Preload(db *gorm.DB) *gorm.DB
}

func NewService[T server](db *gorm.DB, model T) *Service[T] {
	fmt.Printf("%s Service initialized\n", model.Type())
	return &Service[T]{db: db}
}

func (s *Service[T]) Add(entity *T) error {
	return s.db.Create(&entity).Error // Remove the &
}

func (s *Service[T]) GetByID(id int) (*T, error) {
	var entity T
	db := s.db

	// If T implements Preloadable, apply custom preloads
	if p, ok := any(&entity).(Preloadable); ok {
		db = p.Preload(db)
	}

	if err := db.First(&entity, id).Error; err != nil {
		return nil, err
	}
	return &entity, nil
}

func (s *Service[T]) GetAll() ([]T, error) {
	var entity []T

	db := s.db

	var dummy T
	if p, ok := any(&dummy).(Preloadable); ok {
		db = p.Preload(db)
	}

	if err := db.Find(&entity).Error; err != nil {
		return nil, err
	}
	return entity, nil
}

func (s *Service[T]) Update(id int, updateEntity *T) error {
	entity, err := s.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("record with id %d not found", id)
		}
		return fmt.Errorf("ID %d retrival error: %w", id, err)
	}

	if err := s.db.Model(entity).Updates(updateEntity).Error; err != nil {
		return fmt.Errorf("failed to update movie with id %d: %w", id, err)
	}
	
	updateEntity, err = s.GetByID(id)
	if err != nil {
		return fmt.Errorf("failed to reload updated record: %w", err)
	}

	return nil
}

func (s *Service[T]) Delete(id int) error {
	entity, err := s.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("record with id %d not found", id)
		}
		return fmt.Errorf("ID %d retrival error: %w", id, err)
	}

	if err := s.db.Delete(entity).Error; err != nil {
		return fmt.Errorf("failed to delete movie with id %d: %w", id, err)
	}

	return nil
}
