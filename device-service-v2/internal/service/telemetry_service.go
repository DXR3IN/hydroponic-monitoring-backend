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
	Broker     *Broker
}

func NewTelemetryService(r repository.TelemetryRepository, jwt *utils.JWTManager, ds repository.DeviceRepository) *TelemetryService {
	broker := NewBroker()
	return &TelemetryService{repo: r, jwt: jwt, deviceRepo: ds, Broker: broker}
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
	data := insertedRepoData.ToDomain()

	// gave info to the broker that there is a new data on database
	s.Broker.Notifier <- data

	return data, nil
}

var TelemetryStream = make(chan *models.Telemetry)

func (s *TelemetryService) GetLatestTelemetryByDeviceID(deviceID string) (*models.Telemetry, error) {
	telemetries, err := s.repo.GetLatestTelemetryByDeviceID(deviceID)
	if err != nil {
		return nil, err
	}
	return telemetries, nil
}
