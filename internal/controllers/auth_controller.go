package controllers

import (
	"backend_job_assistant/internal/repositories"
	"backend_job_assistant/internal/services"
	"backend_job_assistant/internal/views"
	"encoding/json"
	"net/http"
	"time"
)

type AuthController struct {
	TelegramAuth services.TelegramAuthService
	Sessions     services.SessionService
	Users        repositories.UserRepository
}

type authRequest struct {
	InitData string `json:"init_data"`
}

func (c *AuthController) TelegramLogin(w http.ResponseWriter, r *http.Request) {
	var req authRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.InitData == "" {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	verified, err := c.TelegramAuth.Verify(req.InitData)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	user, err := c.Users.UpsertByTelegramID(
		verified.User.ID,
		verified.User.FirstName,
		verified.User.LastName,
		verified.User.Username,
		verified.User.LanguageCode,
	)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	token, err := c.Sessions.Mint(user.ID)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(c.Sessions.TTL),
	})

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(views.AuthResponse{User: user})
}
