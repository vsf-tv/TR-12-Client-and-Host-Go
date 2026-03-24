package models

// Account represents a user account.
type Account struct {
	AccountID    string `json:"account_id"`
	Username     string `json:"username"`
	PasswordHash string `json:"-"`
	DisplayName  string `json:"display_name"`
	CreatedAt    string `json:"created_at"`
}

// RegisterRequest for POST /account/register.
type RegisterRequest struct {
	Username    string `json:"username"`
	Password    string `json:"password"`
	DisplayName string `json:"display_name"`
}

// LoginRequest for POST /account/login.
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// AuthTokenResponse returned on register/login.
type AuthTokenResponse struct {
	Account *Account `json:"account"`
	Token   string   `json:"token"`
}
