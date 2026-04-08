// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
package db

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

// Store wraps the SQLite database connection.
type Store struct {
	DB *sql.DB
}

// New opens (or creates) the SQLite database and runs migrations.
func New(dbPath string) (*Store, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}
	// Enable WAL mode and foreign keys
	if _, err := db.Exec("PRAGMA journal_mode=WAL"); err != nil {
		return nil, fmt.Errorf("set WAL mode: %w", err)
	}
	if _, err := db.Exec("PRAGMA foreign_keys=ON"); err != nil {
		return nil, fmt.Errorf("enable foreign keys: %w", err)
	}

	s := &Store{DB: db}
	if err := s.migrate(); err != nil {
		return nil, fmt.Errorf("migrate: %w", err)
	}
	return s, nil
}

// Close closes the database connection.
func (s *Store) Close() error {
	return s.DB.Close()
}

func (s *Store) migrate() error {
	schema := `
CREATE TABLE IF NOT EXISTS service_config (
    key   TEXT PRIMARY KEY,
    value BLOB NOT NULL
);

CREATE TABLE IF NOT EXISTS accounts (
    account_id    TEXT PRIMARY KEY,
    username      TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    display_name  TEXT NOT NULL,
    created_at    TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS devices (
    device_id              TEXT PRIMARY KEY,
    account_id             TEXT NOT NULL DEFAULT '',
    device_type            TEXT NOT NULL,
    state                  TEXT NOT NULL DEFAULT 'PAIRING',
    registration           TEXT,
    desired_config         TEXT,
    actual_config          TEXT,
    status                 TEXT,
    online                 INTEGER NOT NULL DEFAULT 0,
    last_seen              TEXT,
    source_ip              TEXT,
    paired_at              TEXT NOT NULL,
    registration_expires_at TEXT,
    current_cert_pem       TEXT,
    previous_cert_pem      TEXT,
    cert_expires_at        TEXT,
    prev_cert_expires_at   TEXT,
    last_rotation_at       TEXT,
    csr_pem                TEXT,
    pairing_code           TEXT,
    access_code            TEXT,
    pairing_expires_at     TEXT,
    config_update_id       INTEGER NOT NULL DEFAULT 0,
    location_name          TEXT,
    device_name            TEXT,
    rotation_interval_days INTEGER NOT NULL DEFAULT 365
);

CREATE TABLE IF NOT EXISTS thumbnails (
    device_id    TEXT NOT NULL,
    source_id    TEXT NOT NULL,
    image_data   BLOB NOT NULL,
    timestamp    TEXT NOT NULL,
    image_type   TEXT NOT NULL,
    image_size_kb INTEGER NOT NULL,
    PRIMARY KEY (device_id, source_id)
);

CREATE TABLE IF NOT EXISTS device_logs (
    device_id   TEXT PRIMARY KEY,
    log_data    BLOB NOT NULL,
    uploaded_at TEXT NOT NULL,
    log_size_kb INTEGER NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_devices_account ON devices(account_id);
CREATE INDEX IF NOT EXISTS idx_devices_pairing_code ON devices(pairing_code);
`
	if _, err := s.DB.Exec(schema); err != nil {
		return err
	}
	// Additive migrations for existing databases
	migrations := []string{
		`ALTER TABLE devices ADD COLUMN location_name TEXT`,
		`ALTER TABLE devices ADD COLUMN device_name TEXT`,
		`ALTER TABLE devices ADD COLUMN rotation_interval_days INTEGER NOT NULL DEFAULT 365`,
	}
	for _, m := range migrations {
		if _, err := s.DB.Exec(m); err != nil {
			// Ignore "duplicate column" errors — column already exists
			if !isDuplicateColumnError(err) {
				return fmt.Errorf("migration failed: %s: %w", m, err)
			}
		}
	}
	return nil
}

func isDuplicateColumnError(err error) bool {
	return err != nil && (contains(err.Error(), "duplicate column") || contains(err.Error(), "already exists"))
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsStr(s, substr))
}

func containsStr(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
