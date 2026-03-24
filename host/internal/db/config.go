package db

// GetConfig retrieves a value from the service_config table.
func (s *Store) GetConfig(key string) ([]byte, error) {
	var value []byte
	err := s.DB.QueryRow("SELECT value FROM service_config WHERE key = ?", key).Scan(&value)
	if err != nil {
		return nil, err
	}
	return value, nil
}

// SetConfig upserts a value in the service_config table.
func (s *Store) SetConfig(key string, value []byte) error {
	_, err := s.DB.Exec(
		"INSERT INTO service_config (key, value) VALUES (?, ?) ON CONFLICT(key) DO UPDATE SET value = excluded.value",
		key, value,
	)
	return err
}
