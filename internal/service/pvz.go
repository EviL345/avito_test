package service

import (
	"github.com/EviL345/avito_test/internal/metrics"
	"github.com/EviL345/avito_test/internal/model/dto/response"
	"github.com/EviL345/avito_test/internal/model/entity"
	"time"
)

type PVZRepository interface {
	CreatePvz(pvz *entity.Pvz) (*entity.Pvz, error)
	GetPvz(startDate, endDate *time.Time, page, limit *int) ([]response.PvzInfo, error)
}

type PVZService struct {
	pvzRepo PVZRepository
}

func NewPVZService(pvzRepo PVZRepository) *PVZService {
	return &PVZService{
		pvzRepo: pvzRepo,
	}
}

func (s *PVZService) CreatePvz(pvz *entity.Pvz) (*entity.Pvz, error) {
	pvz, err := s.pvzRepo.CreatePvz(pvz)
	if err != nil {
		return nil, err
	}

	metrics.CreatePVZ()

	return pvz, nil
}

func (s *PVZService) GetPvz(startDate, endDate *time.Time, page, limit *int) ([]response.PvzInfo, error) {
	return s.pvzRepo.GetPvz(startDate, endDate, page, limit)
}
