DROP TABLE IF EXISTS vehicle_locations;
DROP TABLE IF EXISTS geofence_areas;

CREATE TABLE vehicle_locations (
    id SERIAL PRIMARY KEY,
    vehicle_id VARCHAR(50) NOT NULL,
    latitude DOUBLE PRECISION NOT NULL,
    longitude DOUBLE PRECISION NOT NULL,
    timestamp BIGINT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_vehicle_id ON vehicle_locations(vehicle_id);
CREATE INDEX idx_timestamp ON vehicle_locations(timestamp);

CREATE TABLE IF NOT EXISTS geofence_areas (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    center_latitude DOUBLE PRECISION NOT NULL, 
    center_longitude DOUBLE PRECISION NOT NULL,
    radius_meters INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Insert sample geofence areas
INSERT INTO geofence_areas (name, center_latitude, center_longitude, radius_meters) VALUES
('Halte Bundaran HI', -6.1944, 106.8229, 50),
('Halte Cililitan (PGC)', -6.2460, 106.8990, 50),
('Halte Cawang', -6.2380, 106.8550, 50),
('Halte Blok M', -6.2431, 106.8018, 50),
('Halte Pulo Gadung', -6.1920, 106.8950, 50);