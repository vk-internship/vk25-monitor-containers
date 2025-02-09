package services

import (
	"monitoring-backend/models"
	"monitoring-backend/repositories"
)

type PingService struct {
	repo *repositories.PingRepository
}

func NewPingService(repo *repositories.PingRepository) *PingService {
	return &PingService{repo: repo}
}

func (s *PingService) GetAllPings() ([]models.Ping, error) {
	return s.repo.GetAll()
}

func (s *PingService) CreateOrUpdatePing(result models.Ping) error {
	return s.repo.CreateOrUpdate(result)
}
