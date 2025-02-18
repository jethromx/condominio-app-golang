package repository

import (
	"errors"

	"gorm.io/gorm"
)

type Repository[T any] interface {
	Create(entity *T) error
	Update(entity *T) error
	UpdatePreLoad(entity *T, preload string) error
	Delete(entity *T) error
	FindAll(page, pageSize int) ([]T, int64, error)

	FindID(entity *T, id uint) error
	FindField(entity *T, field, value string) error
	FindByIDPreload(entity *T, id uint, preload string) error

	FindAllWithPreloadRel(page, pageSize int, preload string) ([]T, int64, error)
}

var (
	ErrorEntittyNoyFound = errors.New("records not found")
	ErrorCountRecords    = errors.New("error to count records")
)

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

func (r *BaseRepository[T]) UpdatePreLoad(entity *T, preload string) error {
	return r.DB.Save(entity).Preload(preload).Error
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
		return nil, 0, ErrorEntittyNoyFound
	}

	if err := r.DB.Model(&entities).Count(&totalRecords).Error; err != nil {
		return nil, 0, ErrorCountRecords
	}

	return entities, totalRecords, nil
}

func (r *BaseRepository[T]) FindAllWithPreloadRel(page, pageSize int, preload string) ([]T, int64, error) {
	var entities []T
	var totalRecords int64

	if len(preload) > 0 {
		offset := (page - 1) * pageSize

		if err := r.DB.Offset(offset).Limit(pageSize).Preload(preload).Find(&entities).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return nil, 0, err
			}
			return nil, 0, ErrorEntittyNoyFound
		}

		if err := r.DB.Model(&entities).Count(&totalRecords).Error; err != nil {
			return nil, 0, ErrorCountRecords
		}

		return entities, totalRecords, nil
	} else {
		return r.FindAll(page, pageSize)

	}
}

func (r *BaseRepository[T]) FindID(entity *T, id uint) error {

	if err := r.DB.Find(&entity, id).First(&entity).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return err
		}
		return ErrorEntittyNoyFound
	}
	return nil
}

func (r *BaseRepository[T]) FindByIDPreload(entity *T, id uint, preload string) error {

	if len(preload) > 0 {
		if err := r.DB.Find(&entity, id).Preload(preload).First(&entity).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return err
			}
			return ErrorEntittyNoyFound
		}
		return nil
	} else {
		if err := r.DB.Find(&entity, id).First(&entity).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return err
			}
			return ErrorEntittyNoyFound
		}
		return nil
	}

}

func (r *BaseRepository[T]) FindField(entity *T, field, value string) error {

	if err := r.DB.Where(field+" = ?", value).First(&entity).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil
		}
		return err
	}
	return nil
}
