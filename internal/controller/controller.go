package controller

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/WebCraftersGH/Auth-service/internal/contracts"
	"github.com/WebCraftersGH/Auth-service/internal/controller/dto"
	"github.com/WebCraftersGH/Auth-service/internal/domain"
)

type Controller struct {
	authSVC contracts.AuthSVC
	logger  contracts.ILogger
}

func New(authSVC contracts.AuthSVC, logger contracts.ILogger) *Controller {
	return &Controller{authSVC: authSVC, logger: logger}
}

func (c *Controller) Register(mux *http.ServeMux) {
	mux.HandleFunc("/health", c.handleHealth)
	mux.HandleFunc("/api/v1/auth/start", c.handleStartAuth)
	mux.HandleFunc("/api/v1/auth/verify", c.handleVerifyOTP)
	mux.HandleFunc("/api/v1/auth/check", c.handleAuthCheck)
	mux.HandleFunc("/api/v1/auth/logout", c.handleLogout)
}

func (c *Controller) handleHealth(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, dto.MessageResponse{Message: "ok"})
}

func (c *Controller) handleStartAuth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, dto.ErrorResponse{Error: "method not allowed"})
		return
	}

	var req dto.StartAuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, dto.ErrorResponse{Error: "invalid json"})
		return
	}

	if err := c.authSVC.StartAuth(r.Context(), req.Email); err != nil {
		writeDomainError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, dto.MessageResponse{Message: "otp sent"})
}

func (c *Controller) handleVerifyOTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, dto.ErrorResponse{Error: "method not allowed"})
		return
	}

	var req dto.VerifyOTPRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, dto.ErrorResponse{Error: "invalid json"})
		return
	}

	token, err := c.authSVC.OTPCheck(r.Context(), req.Email, req.Code)
	if err != nil {
		writeDomainError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, dto.VerifyOTPResponse{Token: token.Value})
}

func (c *Controller) handleAuthCheck(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost && r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, dto.ErrorResponse{Error: "method not allowed"})
		return
	}

	token := extractBearerToken(r.Header.Get("Authorization"))
	if token == "" {
		writeJSON(w, http.StatusUnauthorized, dto.ErrorResponse{Error: domain.ErrUnauthorized.Error()})
		return
	}

	if err := c.authSVC.AuthCheck(r.Context(), token); err != nil {
		writeDomainError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, dto.MessageResponse{Message: "authorized"})
}

func (c *Controller) handleLogout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, dto.ErrorResponse{Error: "method not allowed"})
		return
	}

	var req dto.LogoutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, dto.ErrorResponse{Error: "invalid json"})
		return
	}

	if err := c.authSVC.Logout(r.Context(), req.Email); err != nil {
		writeDomainError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, dto.MessageResponse{Message: "logged out"})
}

func extractBearerToken(value string) string {
	const prefix = "Bearer "
	if !strings.HasPrefix(value, prefix) {
		return ""
	}
	return strings.TrimSpace(strings.TrimPrefix(value, prefix))
}

func writeDomainError(w http.ResponseWriter, err error) {
	status := http.StatusInternalServerError

	switch {
	case errors.Is(err, domain.ErrInvalidEmail), errors.Is(err, domain.ErrInvalidOTP), errors.Is(err, domain.ErrOTPExpired):
		status = http.StatusBadRequest
	case errors.Is(err, domain.ErrToManyOTPAttempts):
		status = http.StatusTooManyRequests
	case errors.Is(err, domain.ErrUnauthorized), errors.Is(err, domain.ErrTokenNotFound), errors.Is(err, domain.ErrTokenExpired):
		status = http.StatusUnauthorized
	case errors.Is(err, domain.ErrUserNotFound):
		status = http.StatusNotFound
	}

	writeJSON(w, status, dto.ErrorResponse{Error: err.Error()})
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
