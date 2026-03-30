-- server/db/schema.sql

-- Devices table
CREATE TABLE IF NOT EXISTS devices (
    id TEXT PRIMARY KEY,
    key TEXT UNIQUE NOT NULL,
    name TEXT,
    os_type TEXT,
    last_seen INTEGER,
    created_at INTEGER
);

-- Sessions table
CREATE TABLE IF NOT EXISTS sessions (
    id TEXT PRIMARY KEY,
    device_id TEXT,
    client_conn_id TEXT,
    agent_conn_id TEXT,
    started_at INTEGER,
    FOREIGN KEY (device_id) REFERENCES devices(id)
);

-- Index for faster device lookups by key
CREATE INDEX IF NOT EXISTS idx_devices_key ON devices(key);
