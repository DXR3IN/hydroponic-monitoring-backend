package service

import (
	"errors"
	"time"

	models "github.com/DXR3IN/device-service-v2/internal/domain"
	"github.com/DXR3IN/device-service-v2/internal/repository"
	"github.com/DXR3IN/device-service-v2/internal/utils"
)

var (
	ErrDeviceExists   = errors.New("device already exists")
	ErrDeviceNotFound = errors.New("device not found")
)

type DeviceService struct {
	repo   repository.DeviceRepository
	jwt    *utils.JWTManager
	Broker *Broker
}

func NewDeviceService(r repository.DeviceRepository, jwt *utils.JWTManager) *DeviceService {
	return &DeviceService{repo: r, jwt: jwt}
}

func (s *DeviceService) CreateDevice(deviceID string, deviceName, ownerID string) (*repository.Device, error) {
	ex, err := s.repo.FindByID(deviceID)
	if err != nil {
		return nil, err
	}
	if ex != nil {
		return nil, ErrDeviceExists
	}

	// d is a device variabel
	d := &repository.Device{ID: deviceID, DeviceName: deviceName, OwnerID: ownerID}
	if err := s.repo.Create(d); err != nil {
		return nil, err
	}

	return d, nil
}

func (s *DeviceService) GetDeviceWithID(deviceID string) (*models.Device, error) {
	// d is a device variable
	d, err := s.repo.FindByID(deviceID)
	if d == nil {
		return nil, ErrDeviceNotFound
	}
	if err != nil {
		return nil, err
	}
	return d, nil
}

func (s *DeviceService) GetAllDeviceWithOwnerID(ownerID string) ([]*models.Device, error) {
	// d is a device variable
	d, err := s.repo.FindByOwnerID(ownerID)
	if d == nil {
		return nil, ErrDeviceNotFound
	}
	if err != nil {
		return nil, err
	}

	return d, nil
}

func (s *DeviceService) DeleteDevicesByOwnerID(ownerID string) error {
	return s.repo.DeleteByOwnerID(ownerID)
}

func (s *DeviceService) UpdateDevice(ownerID, deviceID, deviceName string) (*models.Device, error) {
	_, err := s.repo.FindByID(deviceID)
	if err != nil {
		return nil, err
	}

	d, err := s.repo.UpdateDevice(ownerID, deviceID, deviceName)
	if err != nil {
		return nil, err
	}
	return d, nil
}

func (s *DeviceService) UpdateDeviceStatusByID(deviceID, status string) (*models.Device, error) {
	_, err := s.repo.FindByID(deviceID)
	if err != nil {
		return nil, err
	}

	d, err := s.repo.UpdateStatusByID(deviceID, status)
	if err != nil {
		return nil, err
	}
	if s.Broker != nil {
		notification := struct {
			ID     string    `json:"device_id"`
			Status string    `json:"status"`
			Time   time.Time `json:"time"`
		}{
			ID:     d.ID,
			Status: d.Status,
			Time:   time.Now(),
		}

		s.Broker.Notifier <- notification
	}
	return d, nil
}

func (s *DeviceService) VerifyToken(token string) (string, error) {
	claims, err := s.jwt.Verify(token)
	if err != nil {
		return "", err
	}

	return claims.Subject, nil
}
