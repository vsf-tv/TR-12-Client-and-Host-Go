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
package api

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vsf-tv/TR-12-Client-and-Host-Go/host/internal/models"
	"github.com/vsf-tv/TR-12-Client-and-Host-Go/host/internal/service"
)

// AccountHandlers handles /account/* endpoints.
type AccountHandlers struct {
	accountSvc *service.AccountService
}

// NewAccountHandlers creates account handlers.
func NewAccountHandlers(accountSvc *service.AccountService) *AccountHandlers {
	return &AccountHandlers{accountSvc: accountSvc}
}

// Register handles POST /account/register.
func (h *AccountHandlers) Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body", "code": 400})
		return
	}
	acct, token, err := h.accountSvc.Register(req.Username, req.Password, req.DisplayName)
	if err != nil {
		if errors.Is(err, service.ErrConflict) {
			c.JSON(http.StatusConflict, gin.H{"error": "username already exists", "code": 409})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "code": 400})
		return
	}
	c.JSON(http.StatusOK, models.AuthTokenResponse{Account: acct, Token: token})
}

// Login handles POST /account/login.
func (h *AccountHandlers) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body", "code": 400})
		return
	}
	acct, token, err := h.accountSvc.Login(req.Username, req.Password)
	if err != nil {
		if errors.Is(err, service.ErrUnauthorized) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials", "code": 401})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "code": 500})
		return
	}
	c.JSON(http.StatusOK, models.AuthTokenResponse{Account: acct, Token: token})
}

// GetAccount handles GET /account.
func (h *AccountHandlers) GetAccount(c *gin.Context) {
	accountID := c.GetString("account_id")
	acct, err := h.accountSvc.GetAccount(accountID)
	if err != nil || acct == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "account not found", "code": 404})
		return
	}
	c.JSON(http.StatusOK, acct)
}
