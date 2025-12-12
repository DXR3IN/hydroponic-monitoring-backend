package repository

import (
	"errors"
	"time"

	models "github.com/DXR3IN/device-service-v2/internal/domain"
	"gorm.io/gorm"
)

type Telemetry struct {
	ID                       string  `gorm:"primaryKey;type:varchar(36);not null"`
	DeviceID                 string  `gorm:"type:varchar(36);not null"`
	Ppm                      float64 `gorm:"not null"`
	WaterLevelOnPlant        float64 `gorm:"not null"`
	WaterLevelOnCondenser    float64 `gorm:"not null"`
	WaterLevelOnNutrientTank float64 `gorm:"not null"`
	Humidity                 float64 `gorm:"not null"`
	CreatedAt                time.Time
}

func (t *Telemetry) ToDomain() *models.Telemetry {
	if t == nil {
		return nil
	}
	return &models.Telemetry{
		ID:                       t.ID,
		DeviceID:                 t.DeviceID,
		Ppm:                      t.Ppm,
		WaterLevelOnPlant:        t.WaterLevelOnPlant,
		WaterLevelOnCondenser:    t.WaterLevelOnCondenser,
		WaterLevelOnNutrientTank: t.WaterLevelOnNutrientTank,
		Humidity:                 t.Humidity,
		CreatedAt:                t.CreatedAt,
	}
}

func ToRepository(t *models.Telemetry) *Telemetry {
	if t == nil {
		return nil
	}
	return &Telemetry{
		DeviceID:                 t.DeviceID,
		Ppm:                      t.Ppm,
		WaterLevelOnPlant:        t.WaterLevelOnPlant,
		WaterLevelOnCondenser:    t.WaterLevelOnCondenser,
		WaterLevelOnNutrientTank: t.WaterLevelOnNutrientTank,
		Humidity:                 t.Humidity,
	}
}

type TelemetryRepository interface {
	TelemetryInserted(t *Telemetry) (*Telemetry, error)
	GetTelemetryByDeviceID(duration time.Duration, deviceID string) ([]*models.Telemetry, error)
	GetLatestTelemetryByDeviceID(deviceID string) (*models.Telemetry, error)
}

type telemetryRepo struct {
	db *gorm.DB
}

func NewTelemetryRepository(db *gorm.DB) TelemetryRepository {
	return &telemetryRepo{db: db}
}

func (r *telemetryRepo) TelemetryInserted(t *Telemetry) (*Telemetry, error) {
	if err := r.db.Create(t).Error; err != nil {
		return nil, err
	}
	return t, nil
}

func (r *telemetryRepo) GetTelemetryByDeviceID(duration time.Duration, deviceID string) ([]*models.Telemetry, error) {
	var telemetry []Telemetry
	timeStart := time.Now().Add(-duration)
	if err := r.db.
		Where("device_id = ? AND created_at >= ?", deviceID, timeStart).
		Order("created_at DESC").
		Find(&telemetry).Error; err != nil {

		return nil, err
	}

	var result []*models.Telemetry
	for _, t := range telemetry {
		result = append(result, t.ToDomain())
	}

	return result, nil
}

func (r *telemetryRepo) GetLatestTelemetryByDeviceID(deviceID string) (*models.Telemetry, error) {
	var telemetry Telemetry
	if err := r.db.
		Where("device_id = ?", deviceID).
		Order("created_at DESC").
		First(&telemetry).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return telemetry.ToDomain(), nil
}
