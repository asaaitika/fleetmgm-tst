package models

import (
    "time"
)

// VehicleLocation is the main struct for vehicle location data
type VehicleLocation struct {
    ID        int       `json:"id" db:"id"`
    VehicleID string    `json:"vehicle_id" db:"vehicle_id"` 
    Latitude  float64   `json:"latitude" db:"latitude"`
    Longitude float64   `json:"longitude" db:"longitude"`
    Timestamp int64     `json:"timestamp" db:"timestamp"`
    CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// VehicleStatus is used for tracking last known status of a vehicle
type MQTTPayload struct {
    VehicleID string  `json:"vehicle_id"`
    Latitude  float64 `json:"latitude"`
    Longitude float64 `json:"longitude"`
    Timestamp int64   `json:"timestamp"`
}

// VehicleStatus tracks the last known location and geofence status of a vehicle
type GeofenceArea struct {
    ID              int     `json:"id" db:"id"`
    Name            string  `json:"name" db:"name"`
    CenterLatitude  float64 `json:"center_latitude" db:"center_latitude"`
    CenterLongitude float64 `json:"center_longitude" db:"center_longitude"`
    RadiusMeters    int     `json:"radius_meters" db:"radius_meters"`
}

// GeofenceEvent sended when a vehicle enters or exits a geofence area
type GeofenceEvent struct {
    VehicleID string    `json:"vehicle_id"`
    Event     string    `json:"event"`     // "geofence_entry" or "geofence_exit"
    Location  Location  `json:"location"`
    Timestamp int64     `json:"timestamp"`
    AreaName  string    `json:"area_name,omitempty"`
}

// Location is used in GeofenceEvent
type Location struct {
    Latitude  float64 `json:"latitude"`
    Longitude float64 `json:"longitude"`
}

// HistoryRequest for querying vehicle location history
type HistoryRequest struct {
    VehicleID string `json:"-"`
    Start     int64  `form:"start" binding:"required"`
    End       int64  `form:"end" binding:"required"`
}