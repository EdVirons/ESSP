-- Initialize all databases for ESSP services
-- This script runs automatically when PostgreSQL container starts

CREATE DATABASE ssp_ims;
CREATE DATABASE ssp_school;
CREATE DATABASE ssp_devices;
CREATE DATABASE ssp_parts;
CREATE DATABASE ssp_hr;

-- Grant all privileges to the essp user on all databases
GRANT ALL PRIVILEGES ON DATABASE ssp_ims TO essp;
GRANT ALL PRIVILEGES ON DATABASE ssp_school TO essp;
GRANT ALL PRIVILEGES ON DATABASE ssp_devices TO essp;
GRANT ALL PRIVILEGES ON DATABASE ssp_parts TO essp;
GRANT ALL PRIVILEGES ON DATABASE ssp_hr TO essp;
