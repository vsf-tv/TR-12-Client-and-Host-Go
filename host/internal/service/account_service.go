package service

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/vsf-tv/TR-12-Client-and-Host-Go/host/internal/db"
	"github.com/vsf-tv/TR-12-Client-and-Host-Go/host/internal/models"
	"golang.org/x/crypto/bcrypt"
)

// AccountService handles account registration, login, and JWT validation.
type AccountService struct {
	store     *db.Store
	jwtSecret []byte
	expiry    time.Duration
}

// Claims represents JWT token claims.
type Claims struct {
	AccountID string `json:"account_id"`
	Username  string `json:"username"`
	jwt.RegisteredClaims
}

// NewAccountService creates a new AccountService.
func NewAccountService(store *db.Store, jwtSecret []byte, expiryHours int) *AccountService {
	return &AccountService{
		store:     store,
		jwtSecret: jwtSecret,
		expiry:    time.Duration(expiryHours) * time.Hour,
	}
}

// Register creates a new account and returns it with a JWT.
func (s *AccountService) Register(username, password, displayName string) (*models.Account, string, error) {
	if len(username) < 3 || len(username) > 64 {
		return nil, "", errors.New("username must be 3-64 characters")
	}
	if len(password) < 8 {
		return nil, "", errors.New("password must be at least 8 characters")
	}

	existing, err := s.store.GetAccountByUsername(username)
	if err != nil {
		return nil, "", err
	}
	if existing != nil {
		return nil, "", ErrConflict
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, "", err
	}

	acct := &models.Account{
		AccountID:    generateAccountID(),
		Username:     username,
		PasswordHash: string(hash),
		DisplayName:  displayName,
		CreatedAt:    time.Now().UTC().Format(time.RFC3339),
	}
	if err := s.store.CreateAccount(acct); err != nil {
		return nil, "", err
	}

	token, err := s.generateToken(acct)
	if err != nil {
		return nil, "", err
	}
	return acct, token, nil
}

// Login validates credentials and returns a JWT.
func (s *AccountService) Login(username, password string) (*models.Account, string, error) {
	acct, err := s.store.GetAccountByUsername(username)
	if err != nil {
		return nil, "", err
	}
	if acct == nil {
		return nil, "", ErrUnauthorized
	}
	if err := bcrypt.CompareHashAndPassword([]byte(acct.PasswordHash), []byte(password)); err != nil {
		return nil, "", ErrUnauthorized
	}
	token, err := s.generateToken(acct)
	if err != nil {
		return nil, "", err
	}
	return acct, token, nil
}

// ValidateToken parses and validates a JWT, returning the claims.
func (s *AccountService) ValidateToken(tokenStr string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return s.jwtSecret, nil
	})
	if err != nil || !token.Valid {
		return nil, ErrUnauthorized
	}
	return claims, nil
}

// GetAccount returns an account by ID.
func (s *AccountService) GetAccount(accountID string) (*models.Account, error) {
	return s.store.GetAccountByID(accountID)
}

func (s *AccountService) generateToken(acct *models.Account) (string, error) {
	claims := &Claims{
		AccountID: acct.AccountID,
		Username:  acct.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.expiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(s.jwtSecret)
}

func generateAccountID() string {
	b := make([]byte, 4)
	rand.Read(b)
	return "acc_" + hex.EncodeToString(b)
}
