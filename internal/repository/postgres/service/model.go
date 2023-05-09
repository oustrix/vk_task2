package service

import "vk_task2/internal/domain"

type Service struct {
	ID       int `gorm:"primaryKey"`
	Name     string
	UserID   int64 `gorm:"not null"`
	Login    string
	Password string
}

func ToDomain(s *Service) *domain.Service {
	return &domain.Service{
		ID:       s.ID,
		Name:     s.Name,
		UserID:   s.UserID,
		Login:    s.Login,
		Password: s.Password,
	}
}

func ToDomainList(services []*Service) []*domain.Service {
	result := make([]*domain.Service, 0, len(services))
	for _, s := range services {
		result = append(result, ToDomain(s))
	}
	return result
}
