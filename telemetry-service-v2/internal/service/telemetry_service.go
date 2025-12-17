package service

import (
	"errors"
	"time"

	models "github.com/DXR3IN/telemetry-service-v2/internal/domain"
	"github.com/DXR3IN/telemetry-service-v2/internal/repository"
	"github.com/DXR3IN/telemetry-service-v2/internal/utils"
	"github.com/google/uuid"
)

var (
	ErrTelemetryNotFound = errors.New("telemetry not found")
)

type TelemetryService struct {
	repo   repository.TelemetryRepository
	jwt    *utils.JWTManager
	Broker *Broker
}

func NewTelemetryService(r repository.TelemetryRepository, jwt *utils.JWTManager) *TelemetryService {
	broker := NewBroker()
	return &TelemetryService{repo: r, jwt: jwt, Broker: broker}
}

func (s *TelemetryService) GetTelemetryByDeviceID(duration time.Duration, deviceID string) ([]*models.Telemetry, error) {
	// ex, err := s.deviceRepo.FindByID(deviceID)
	// if err != nil {
	// 	return nil, err
	// }
	// if ex != nil {
	// 	return nil, ErrDeviceExists
	// }
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
	// deviceID := t.DeviceID
	// ex, err := s.deviceRepo.FindByID(deviceID)
	// if err != nil {
	// 	return nil, err
	// }
	// if ex != nil {
	// 	return nil, ErrDeviceExists
	// }
	t.ID = uuid.New().String()
	repoData := repository.ToRepository(t)

	insertedRepoData, err := s.repo.TelemetryInserted(repoData)
	if err != nil {
		return nil, err
	}
	data := insertedRepoData.ToDomain()

	select {
	case s.Broker.Notifier <- data:
	default:
		// Opsional: 
	}

	return data, nil
}

var TelemetryStream = make(chan *models.Telemetry)

func (s *TelemetryService) GetLatestTelemetryByDeviceID(deviceID string) (*models.Telemetry, error) {
	// ex, err := s.deviceRepo.FindByID(deviceID)
	// if err != nil {
	// 	return nil, err
	// }
	// if ex != nil {
	// 	return nil, ErrDeviceNotFound
	// }
	telemetries, err := s.repo.GetLatestTelemetryByDeviceID(deviceID)
	if err != nil {
		return nil, err
	}
	return telemetries, nil
}
