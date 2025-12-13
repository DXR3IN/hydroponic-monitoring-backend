package models

import "time"

type Telemetry struct {
	ID                       string    `json:"id"`
	DeviceID                 string    `json:"device_id"`
	Ppm                      float64   `json:"ppm"`
	WaterLevelOnPlant        float64   `json:"water_level_on_plant"`
	WaterLevelOnCondenser    float64   `json:"water_level_on_condenser"`
	WaterLevelOnNutrientTank float64   `json:"water_level_on_nutrient_tank"`
	Humidity                 float64   `json:"humidity"`
	CreatedAt                time.Time `json:"created_at"`
}
