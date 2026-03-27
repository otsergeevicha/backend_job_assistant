package views

import "backend_job_assistant/internal/models"

type AuthResponse struct {
	User *models.User `json:"user"`
}

type MeResponse struct {
	User *models.User `json:"user"`
}
