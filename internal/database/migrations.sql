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
('Halte Pinang Ranti', -6.2593, 106.8789, 50),
('Halte Cawang UKI', -6.2426, 106.8585, 50),
('Halte Pancoran Tugu', -6.2253, 106.8401, 50),
('Halte Pertamburan', -6.1679, 106.8038, 50),
('Halte Pluit', -6.1250, 106.7942, 50);