package db

import (
	"database/sql"

	"github.com/vsf-tv/TR-12-Client-and-Host-Go/host/internal/models"
)

// CreateAccount inserts a new account.
func (s *Store) CreateAccount(a *models.Account) error {
	_, err := s.DB.Exec(
		"INSERT INTO accounts (account_id, username, password_hash, display_name, created_at) VALUES (?, ?, ?, ?, ?)",
		a.AccountID, a.Username, a.PasswordHash, a.DisplayName, a.CreatedAt,
	)
	return err
}

// GetAccountByUsername looks up an account by username.
func (s *Store) GetAccountByUsername(username string) (*models.Account, error) {
	a := &models.Account{}
	err := s.DB.QueryRow(
		"SELECT account_id, username, password_hash, display_name, created_at FROM accounts WHERE username = ?",
		username,
	).Scan(&a.AccountID, &a.Username, &a.PasswordHash, &a.DisplayName, &a.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return a, nil
}

// GetAccountByID looks up an account by ID.
func (s *Store) GetAccountByID(accountID string) (*models.Account, error) {
	a := &models.Account{}
	err := s.DB.QueryRow(
		"SELECT account_id, username, password_hash, display_name, created_at FROM accounts WHERE account_id = ?",
		accountID,
	).Scan(&a.AccountID, &a.Username, &a.PasswordHash, &a.DisplayName, &a.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return a, nil
}
