package controllers

import (
	"backend_job_assistant/internal/repositories"
	"backend_job_assistant/internal/views"
	"encoding/json"
	"net/http"
)

type UserController struct {
	Users repositories.UserRepository
}

// @Summary Текущий пользователь
// @Description Возвращает текущего пользователя
// @Tags users
// @Produce json
// @Success 200 {object} views.MeResponse
// @Router /me [get]
func (c *UserController) Me(w http.ResponseWriter, r *http.Request, userID int64) {
	user, err := c.Users.FindByID(userID)
	if err != nil || user == nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(views.MeResponse{User: user})
}
