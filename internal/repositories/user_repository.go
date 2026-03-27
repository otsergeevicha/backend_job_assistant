package repositories

import "backend_job_assistant/internal/models"

type UserRepository interface {
	UpsertByTelegramID(tgID int64, firstName, lastName, username, languageCode string) (*models.User, error)
	FindByID(id int64) (*models.User, error)
}
