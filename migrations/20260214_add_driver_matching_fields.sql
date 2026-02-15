-- Add new trip statuses and search tracking fields
ALTER TABLE trips MODIFY COLUMN status VARCHAR(30);

-- Add search tracking fields
ALTER TABLE trips 
  ADD COLUMN IF NOT EXISTS search_started_at TIMESTAMP NULL,
  ADD COLUMN IF NOT EXISTS search_ended_at TIMESTAMP NULL,
  ADD COLUMN IF NOT EXISTS cancelled_by UUID NULL,
  ADD COLUMN IF NOT EXISTS cancellation_reason VARCHAR(255) NULL;

-- Create driver availability table
CREATE TABLE IF NOT EXISTS driver_availability (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  driver_id UUID NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
  is_available BOOLEAN DEFAULT false,
  service_types VARCHAR(100)[] DEFAULT '{}',
  last_active_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Add index for searching trips efficiently
CREATE INDEX IF NOT EXISTS idx_trips_status_search ON trips(status, search_started_at);

-- Add index for driver availability queries
CREATE INDEX IF NOT EXISTS idx_driver_availability_status ON driver_availability(is_available, service_types);
