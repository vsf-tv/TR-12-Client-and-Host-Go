package service

import (
	"strings"
	"testing"
	"time"

	"github.com/vsf-tv/TR-12-Client-and-Host-Go/host/internal/db"
)

func newTestAccountService(t *testing.T) *AccountService {
	t.Helper()
	store, err := db.New(":memory:")
	if err != nil {
		t.Fatalf("db.New: %v", err)
	}
	t.Cleanup(func() { store.Close() })
	return NewAccountService(store, []byte("test-jwt-secret-32bytes-long!!!!"), 24)
}

func TestRegister_Success(t *testing.T) {
	svc := newTestAccountService(t)
	acct, token, err := svc.Register("testuser", "password123", "Test User")
	if err != nil {
		t.Fatalf("Register: %v", err)
	}
	if acct == nil {
		t.Fatal("expected account")
	}
	if !strings.HasPrefix(acct.AccountID, "acc_") {
		t.Fatalf("expected acc_ prefix, got %q", acct.AccountID)
	}
	if acct.Username != "testuser" {
		t.Fatalf("expected testuser, got %q", acct.Username)
	}
	if token == "" {
		t.Fatal("expected non-empty token")
	}
}

func TestRegister_DuplicateUsername(t *testing.T) {
	svc := newTestAccountService(t)
	svc.Register("testuser", "password123", "Test User")
	_, _, err := svc.Register("testuser", "password456", "Another User")
	if err != ErrConflict {
		t.Fatalf("expected ErrConflict, got %v", err)
	}
}

func TestRegister_ShortUsername(t *testing.T) {
	svc := newTestAccountService(t)
	_, _, err := svc.Register("ab", "password123", "Short")
	if err == nil || !strings.Contains(err.Error(), "username") {
		t.Fatalf("expected username validation error, got %v", err)
	}
}

func TestRegister_ShortPassword(t *testing.T) {
	svc := newTestAccountService(t)
	_, _, err := svc.Register("testuser", "short", "Test")
	if err == nil || !strings.Contains(err.Error(), "password") {
		t.Fatalf("expected password validation error, got %v", err)
	}
}

func TestLogin_Success(t *testing.T) {
	svc := newTestAccountService(t)
	svc.Register("testuser", "password123", "Test User")

	acct, token, err := svc.Login("testuser", "password123")
	if err != nil {
		t.Fatalf("Login: %v", err)
	}
	if acct == nil || token == "" {
		t.Fatal("expected account and token")
	}
}

func TestLogin_WrongPassword(t *testing.T) {
	svc := newTestAccountService(t)
	svc.Register("testuser", "password123", "Test User")

	_, _, err := svc.Login("testuser", "wrongpassword")
	if err != ErrUnauthorized {
		t.Fatalf("expected ErrUnauthorized, got %v", err)
	}
}

func TestLogin_NonexistentUser(t *testing.T) {
	svc := newTestAccountService(t)
	_, _, err := svc.Login("nobody", "password123")
	if err != ErrUnauthorized {
		t.Fatalf("expected ErrUnauthorized, got %v", err)
	}
}

func TestValidateToken_Valid(t *testing.T) {
	svc := newTestAccountService(t)
	_, token, _ := svc.Register("testuser", "password123", "Test User")

	claims, err := svc.ValidateToken(token)
	if err != nil {
		t.Fatalf("ValidateToken: %v", err)
	}
	if claims.Username != "testuser" {
		t.Fatalf("expected testuser, got %q", claims.Username)
	}
	if !strings.HasPrefix(claims.AccountID, "acc_") {
		t.Fatalf("expected acc_ prefix, got %q", claims.AccountID)
	}
}

func TestValidateToken_Expired(t *testing.T) {
	store, _ := db.New(":memory:")
	defer store.Close()
	// Create service with 0 hour expiry so tokens expire immediately
	svc := NewAccountService(store, []byte("test-jwt-secret-32bytes-long!!!!"), 0)
	_, token, _ := svc.Register("testuser", "password123", "Test User")

	// Token should be expired (0 hour expiry means it expires at issuance)
	time.Sleep(2 * time.Second)
	_, err := svc.ValidateToken(token)
	if err != ErrUnauthorized {
		t.Fatalf("expected ErrUnauthorized for expired token, got %v", err)
	}
}

func TestValidateToken_Tampered(t *testing.T) {
	svc := newTestAccountService(t)
	_, token, _ := svc.Register("testuser", "password123", "Test User")

	// Tamper with the token
	tampered := token + "x"
	_, err := svc.ValidateToken(tampered)
	if err != ErrUnauthorized {
		t.Fatalf("expected ErrUnauthorized for tampered token, got %v", err)
	}
}

func TestValidateToken_WrongSecret(t *testing.T) {
	svc := newTestAccountService(t)
	_, token, _ := svc.Register("testuser", "password123", "Test User")

	// Create a different service with a different secret
	store2, _ := db.New(":memory:")
	defer store2.Close()
	svc2 := NewAccountService(store2, []byte("different-secret-32bytes-long!!!"), 24)
	_, err := svc2.ValidateToken(token)
	if err != ErrUnauthorized {
		t.Fatalf("expected ErrUnauthorized for wrong secret, got %v", err)
	}
}

func TestGetAccount(t *testing.T) {
	svc := newTestAccountService(t)
	acct, _, _ := svc.Register("testuser", "password123", "Test User")

	got, err := svc.GetAccount(acct.AccountID)
	if err != nil {
		t.Fatalf("GetAccount: %v", err)
	}
	if got == nil || got.Username != "testuser" {
		t.Fatalf("expected testuser, got %+v", got)
	}
}

func TestGenerateAccountID_Format(t *testing.T) {
	id := generateAccountID()
	if !strings.HasPrefix(id, "acc_") {
		t.Fatalf("expected acc_ prefix, got %q", id)
	}
	if len(id) != 12 { // "acc_" + 8 hex chars
		t.Fatalf("expected length 12, got %d (%q)", len(id), id)
	}
}
