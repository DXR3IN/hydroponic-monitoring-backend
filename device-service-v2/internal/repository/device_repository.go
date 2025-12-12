package repository

import (
	"errors"
	"time"

	models "github.com/DXR3IN/device-service-v2/internal/domain"
	"gorm.io/gorm"
)

type Device struct {
	ID         string `gorm:"primaryKey;type:varchar(36);not null"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeviceName string `gorm:"not null"`
	OwnerID    string `gorm:"type:uuid;not null"`
}

func (d *Device) ToDomain() *models.Device {
	if d == nil {
		return nil
	}

	return &models.Device{
		ID:         d.ID,
		DeviceName: d.DeviceName,
		OwnerID:    d.OwnerID,
		CreatedAt:  d.CreatedAt,
		UpdatedAt:  d.UpdatedAt,
	}
}

type DeviceRepository interface {
	Create(d *Device) error
	FindByID(id string) (*models.Device, error)
	FindByOwnerID(ownerID string) ([]*models.Device, error)
	DeleteByOwnerID(ownerID string) error
	UpdateDevice(ownerID, deviceID, newDeviceID string, deviceName string) (*models.Device, error)
}

type deviceRepo struct {
	db *gorm.DB
}

func NewDeviceRepository(db *gorm.DB) DeviceRepository {
	return &deviceRepo{db: db}
}

func (r *deviceRepo) Create(d *Device) error {
	return r.db.Create(d).Error
}

func (r *deviceRepo) FindByID(id string) (*models.Device, error) {
	var device Device
	if err := r.db.First(&device, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return device.ToDomain(), nil
}

func (r *deviceRepo) FindByOwnerID(ownerID string) ([]*models.Device, error) {
	var devices []Device
	if err := r.db.Where("owner_id = ?", ownerID).Find(&devices).Error; err != nil {
		return nil, err
	}
	var result []*models.Device
	for _, d := range devices {
		result = append(result, d.ToDomain())
	}
	return result, nil
}

func (r *deviceRepo) DeleteByOwnerID(ownerID string) error {
	return r.db.Where("owner_id = ?", ownerID).Delete(&Device{}).Error
}

func (r *deviceRepo) UpdateDevice(ownerID, deviceID, newDeviceID, deviceName string) (*models.Device, error) {
	var device Device
	if err := r.db.First(&device, "id = ? AND owner_id =?", deviceID, ownerID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	device.DeviceName = deviceName
	device.UpdatedAt = time.Time{}
	if err := r.db.Save(&device).Error; err != nil {
		return nil, err
	}
	return device.ToDomain(), nil
}
