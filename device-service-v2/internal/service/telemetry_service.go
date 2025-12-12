package service

import (
	"errors"
	"time"

	models "github.com/DXR3IN/device-service-v2/internal/domain"
	"github.com/DXR3IN/device-service-v2/internal/repository"
	"github.com/DXR3IN/device-service-v2/internal/utils"
)

var (
	ErrTelemetryNotFound = errors.New("telemetry not found")
)

type TelemetryService struct {
	repo       repository.TelemetryRepository
	deviceRepo repository.DeviceRepository
	jwt        *utils.JWTManager
}

func NewTelemetryService(r repository.TelemetryRepository, jwt *utils.JWTManager, ds repository.DeviceRepository) *TelemetryService {
	return &TelemetryService{repo: r, jwt: jwt, deviceRepo: ds}
}

func (s *TelemetryService) GetTelemetryByDeviceID(duration time.Duration, deviceID string) ([]*models.Telemetry, error) {
	// t is a telemetry variable
	t, err := s.repo.GetTelemetryByDeviceID(duration, deviceID)
	if t == nil {
		return nil, ErrTelemetryNotFound
	}
	if err != nil {
		return nil, err
	}
	return t, nil
}

func (s *TelemetryService) InsertTelemetry(t *models.Telemetry) (*models.Telemetry, error) {
	deviceID := t.DeviceID
	ex, err := s.deviceRepo.FindByID(deviceID)
	if err != nil {
		return nil, err
	}
	if ex != nil {
		return nil, ErrDeviceExists
	}
	repoData := repository.ToRepository(t)

	insertedRepoData, err := s.repo.TelemetryInserted(repoData)
	if err != nil {
		return nil, err
	}
	return insertedRepoData.ToDomain(), nil
}

func (s *TelemetryService) GetLatestTelemetryByDeviceID(deviceID string) (*models.Telemetry, error) {
	telemetries, err := s.repo.GetLatestTelemetryByDeviceID(deviceID)
	if err != nil {
		return nil, err
	}
	return telemetries, nil
}
