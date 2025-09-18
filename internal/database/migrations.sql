-- Remove existing tables if they exist
DROP TABLE IF EXISTS vehicle_locations;
DROP TABLE IF EXISTS geofence_areas;

CREATE TABLE vehicle_locations (
    id SERIAL PRIMARY KEY,                    -- Auto increment ID
    vehicle_id VARCHAR(50) NOT NULL,          -- Plat nomor bus (B1234XYZ)
    latitude DOUBLE PRECISION NOT NULL,       -- Koordinat latitude
    longitude DOUBLE PRECISION NOT NULL,      -- Koordinat longitude
    timestamp BIGINT NOT NULL,                -- Unix timestamp
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_vehicle_id ON vehicle_locations(vehicle_id);
CREATE INDEX idx_timestamp ON vehicle_locations(timestamp);

CREATE TABLE IF NOT EXISTS geofence_areas (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,               -- Nama area (Terminal Blok M)
    center_latitude DOUBLE PRECISION NOT NULL, 
    center_longitude DOUBLE PRECISION NOT NULL,
    radius_meters INTEGER NOT NULL,           -- Radius dalam meter
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Insert sample geofence areas
INSERT INTO geofence_areas (name, center_latitude, center_longitude, radius_meters) VALUES
('Terminal Blok M', -6.2431, 106.8018, 50),
('Halte Bundaran HI', -6.1944, 106.8229, 50),
('Terminal Kampung Melayu', -6.2241, 106.8664, 50);