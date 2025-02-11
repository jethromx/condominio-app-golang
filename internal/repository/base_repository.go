package repository

import (
	"errors"

	"gorm.io/gorm"
)

type Repository[T any] interface {
	Create(entity *T) error
	Update(entity *T) error
	Delete(entity *T) error
	FindAll(page, pageSize int) ([]T, int64, error)
	FindByID(entity *T, id uint) (*T, error)
	FindByField(entity *T, field, value string) (*T, error)
	FindAllWithPreloadRel(page, pageSize int, preload string) ([]T, int64, error)
	FindByIDWithPreload(entity *T, id uint, preload string) (*T, error)
}

type BaseRepository[T any] struct {
	DB *gorm.DB
}

func NewBaseRepository[T any](db *gorm.DB) *BaseRepository[T] {
	return &BaseRepository[T]{DB: db}
}

func (r *BaseRepository[T]) Create(entity *T) error {
	return r.DB.Create(entity).Error
}

func (r *BaseRepository[T]) Update(entity *T) error {
	return r.DB.Save(entity).Error
}

func (r *BaseRepository[T]) Delete(entity *T) error {
	return r.DB.Delete(entity).Error
}

func (r *BaseRepository[T]) FindAll(page, pageSize int) ([]T, int64, error) {
	var entities []T
	var totalRecords int64

	offset := (page - 1) * pageSize
	if err := r.DB.Offset(offset).Limit(pageSize).Find(&entities).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, 0, err
		}
		return nil, 0, errors.New("error to find all entities")
	}

	if err := r.DB.Model(&entities).Count(&totalRecords).Error; err != nil {
		return nil, 0, errors.New("error to count records")
	}

	return entities, totalRecords, nil
}

func (r *BaseRepository[T]) FindAllWithPreloadRel(page, pageSize int, preload string) ([]T, int64, error) {
	var entities []T
	var totalRecords int64

	offset := (page - 1) * pageSize
	if err := r.DB.Offset(offset).Limit(pageSize).Preload(preload).Find(&entities).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, 0, err
		}
		return nil, 0, errors.New("error to find all entities")
	}

	if err := r.DB.Model(&entities).Count(&totalRecords).Error; err != nil {
		return nil, 0, errors.New("error to count records")
	}

	return entities, totalRecords, nil
}

func (r *BaseRepository[T]) FindByID(entity *T, id uint) (*T, error) {

	if err := r.DB.Find(&entity, id).First(&entity).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, err
		}
		return nil, errors.New("error to find entity by id")
	}
	return entity, nil
}

func (r *BaseRepository[T]) FindByIDWithPreload(entity *T, id uint, preload string) (*T, error) {

	if err := r.DB.Find(&entity, id).Preload(preload).First(&entity).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, err
		}
		return nil, errors.New("error to find entity by id")
	}
	return entity, nil
}

func (r *BaseRepository[T]) FindByField(entity *T, field, value string) (*T, error) {

	if err := r.DB.Where(field+" = ?", value).First(&entity).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return entity, nil
}
