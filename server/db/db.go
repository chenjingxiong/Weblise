package db

import (
	"database/sql"
	"embed"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

//go:embed schema.sql
var SchemaFS embed.FS

type Database struct {
	db *sql.DB
}

// Config holds database configuration
type Config struct {
	Path string // Path to SQLite database file
}

// New creates a new database connection
func New(cfg Config) (*Database, error) {
	db, err := sql.Open("sqlite3", cfg.Path)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	// Enable foreign keys
	if _, err := db.Exec("PRAGMA foreign_keys = ON"); err != nil {
		return nil, fmt.Errorf("enable foreign keys: %w", err)
	}

	// Run schema
	if err := runSchema(db); err != nil {
		return nil, fmt.Errorf("run schema: %w", err)
	}

	return &Database{db: db}, nil
}

// runSchema executes the database schema
func runSchema(db *sql.DB) error {
	schema, err := SchemaFS.ReadFile("schema.sql")
	if err != nil {
		return fmt.Errorf("read schema: %w", err)
	}

	_, err = db.Exec(string(schema))
	return err
}

// Close closes the database connection
func (d *Database) Close() error {
	return d.db.Close()
}

// Device operations

// CreateDevice creates a new device
func (d *Database) CreateDevice(id, key, name, osType string) error {
	now := time.Now().Unix()
	query := `
		INSERT INTO devices (id, key, name, os_type, last_seen, created_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`
	_, err := d.db.Exec(query, id, key, name, osType, now, now)
	return err
}

// GetDeviceByKey retrieves a device by its key
func (d *Database) GetDeviceByKey(key string) (*Device, error) {
	query := `
		SELECT id, key, name, os_type, last_seen, created_at
		FROM devices WHERE key = ?
	`
	row := d.db.QueryRow(query, key)

	var device Device
	var lastSeen, createdAt int64
	err := row.Scan(&device.ID, &device.Key, &device.Name, &device.OSType, &lastSeen, &createdAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	device.LastSeen = time.Unix(lastSeen, 0)
	device.CreatedAt = time.Unix(createdAt, 0)
	return &device, nil
}

// UpdateDeviceLastSeen updates the last seen timestamp for a device
func (d *Database) UpdateDeviceLastSeen(id string) error {
	query := `UPDATE devices SET last_seen = ? WHERE id = ?`
	_, err := d.db.Exec(query, time.Now().Unix(), id)
	return err
}

// Session operations

// CreateSession creates a new session
func (d *Database) CreateSession(id, deviceID, clientConnID, agentConnID string) error {
	query := `
		INSERT INTO sessions (id, device_id, client_conn_id, agent_conn_id, started_at)
		VALUES (?, ?, ?, ?, ?)
	`
	_, err := d.db.Exec(query, id, deviceID, clientConnID, agentConnID, time.Now().Unix())
	return err
}

// DeleteSession deletes a session by ID
func (d *Database) DeleteSession(id string) error {
	_, err := d.db.Exec("DELETE FROM sessions WHERE id = ?", id)
	return err
}

// Device represents a device record
type Device struct {
	ID        string
	Key       string
	Name      string
	OSType    string
	LastSeen  time.Time
	CreatedAt time.Time
}
