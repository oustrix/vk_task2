package service

import (
	"gorm.io/gorm"
	"vk_task2/internal/domain"
)

type serviceRepository struct {
	db *gorm.DB
}

type Repository interface {
	GetByUserID(int64) ([]*domain.Service, error)
	Create(*domain.Service) error
}

func NewRepository(db *gorm.DB) Repository {
	return &serviceRepository{db: db}
}

func (r *serviceRepository) GetByUserID(userID int64) ([]*domain.Service, error) {
	var services []*Service
	err := r.db.Where("user_id = ?", userID).Find(&services).Error
	if err != nil {
		return nil, err
	}

	return toDomainList(services), nil
}

func (r *serviceRepository) Create(service *domain.Service) error {
	s := fromDomain(service)
	return r.db.Create(s).Error
}
