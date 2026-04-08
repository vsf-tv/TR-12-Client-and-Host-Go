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
