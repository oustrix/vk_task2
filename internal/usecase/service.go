package usecase

import (
	"vk_task2/internal/domain"
	"vk_task2/internal/repository/postgres/service"
)

type serviceUsecase struct {
	repository service.Repository
}

type Service interface {
	GetUserServices(int64) ([]*domain.Service, error)
	Create(*domain.Service) error
}

func NewServiceUsecase(repository service.Repository) Service {
	return &serviceUsecase{
		repository: repository,
	}
}

func (u *serviceUsecase) GetUserServices(userID int64) ([]*domain.Service, error) {
	return u.repository.GetByUserID(userID)
}

func (u *serviceUsecase) Create(service *domain.Service) error {
	return u.repository.Create(service)
}
