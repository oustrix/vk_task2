package service

import "vk_task2/internal/domain"

type Service struct {
	ID       int `gorm:"primaryKey"`
	Name     string
	UserID   int64 `gorm:"not null"`
	Login    string
	Password string
}

func toDomain(s *Service) *domain.Service {
	return &domain.Service{
		ID:       s.ID,
		Name:     s.Name,
		UserID:   s.UserID,
		Login:    s.Login,
		Password: s.Password,
	}
}

func toDomainList(services []*Service) []*domain.Service {
	result := make([]*domain.Service, 0, len(services))
	for _, s := range services {
		result = append(result, toDomain(s))
	}
	return result
}

func fromDomain(s *domain.Service) *Service {
	return &Service{
		ID:       s.ID,
		Name:     s.Name,
		UserID:   s.UserID,
		Login:    s.Login,
		Password: s.Password,
	}
}
