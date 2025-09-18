package repositories

import (
	"github.com/asaaitika/fleetmgm-tst/internal/models"
	"github.com/jmoiron/sqlx"
)

// VehicleRepository handles database operations
type VehicleRepository struct {
	db *sqlx.DB
}

// NewVehicleRepository creates new repository instance
func NewVehicleRepository(db *sqlx.DB) *VehicleRepository {
	return &VehicleRepository{db: db}
}

// InsertLocation saves vehicle location to database
func (r *VehicleRepository) InsertLocation(payload *models.MQTTPayload) error {
	query := `
        INSERT INTO vehicle_locations (vehicle_id, latitude, longitude, timestamp)
        VALUES ($1, $2, $3, $4)
    `
	_, err := r.db.Exec(query, payload.VehicleID, payload.Latitude, payload.Longitude, payload.Timestamp)
	return err
}

// GetLastLocation retrieves the last known location of a vehicle
func (r *VehicleRepository) GetLastLocation(vehicleID string) (*models.VehicleLocation, error) {
	var location models.VehicleLocation

	query := `
        SELECT id, vehicle_id, latitude, longitude, timestamp, created_at
        FROM vehicle_locations
        WHERE vehicle_id = $1
        ORDER BY timestamp DESC
        LIMIT 1
    `

	err := r.db.Get(&location, query, vehicleID)
	if err != nil {
		return nil, err
	}

	return &location, nil
}

// GetLocationHistory retrieves location history of a vehicle within a time range
func (r *VehicleRepository) GetLocationHistory(vehicleID string, start, end int64) ([]models.VehicleLocation, error) {
	var locations []models.VehicleLocation

	query := `
        SELECT id, vehicle_id, latitude, longitude, timestamp, created_at
        FROM vehicle_locations
        WHERE vehicle_id = $1 AND timestamp >= $2 AND timestamp <= $3
        ORDER BY timestamp ASC
    `

	err := r.db.Select(&locations, query, vehicleID, start, end)
	if err != nil {
		return nil, err
	}

	return locations, nil
}

// GetGeofenceAreas retrieves all geofence areas from the database
func (r *VehicleRepository) GetGeofenceAreas() ([]models.GeofenceArea, error) {
	var areas []models.GeofenceArea

	query := `
        SELECT id, name, center_latitude, center_longitude, radius_meters
        FROM geofence_areas
    `

	err := r.db.Select(&areas, query)
	return areas, err
}
