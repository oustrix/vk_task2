package usecase

import (
	"vk_task2/internal/domain"
	"vk_task2/internal/repository/postgres/service"
)

type serviceUsecase struct {
	repository service.Repository
}

type Service interface {
	GetUserServices(*domain.User) ([]*domain.Service, error)
	Create(*domain.Service) error
}

func NewServiceUsecase(repository service.Repository) Service {
	return &serviceUsecase{
		repository: repository,
	}
}

func (u *serviceUsecase) GetUserServices(user *domain.User) ([]*domain.Service, error) {
	return u.repository.GetByUserID(user.ID)
}

func (u *serviceUsecase) Create(service *domain.Service) error {
	return u.repository.Create(service)
}
