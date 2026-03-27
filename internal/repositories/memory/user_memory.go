package memory

import (
	"fmt"
	"sync"

	"backend_job_assistant/internal/models"
)

type UserMemoryRepository struct {
	mu     sync.Mutex
	byID   map[int64]*models.User
	byTgID map[int64]int64
	nextID int64
}

func NewUserMemoryRepository() *UserMemoryRepository {
	return &UserMemoryRepository{
		byID:   make(map[int64]*models.User),
		byTgID: make(map[int64]int64),
		nextID: 1,
	}
}

func (r *UserMemoryRepository) UpsertByTelegramID(tgID int64, firstName, lastName, username, languageCode string) (*models.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if id, ok := r.byTgID[tgID]; ok {
		u := r.byID[id]
		if u == nil {
			return nil, fmt.Errorf("user not found")
		}
		if u.FirstName == "" {
			u.FirstName = firstName
		}
		if u.LastName == "" {
			u.LastName = lastName
		}
		if u.Username == "" {
			u.Username = username
		}
		if u.LanguageCode == "" {
			u.LanguageCode = languageCode
		}
		if u.Nickname == "" {
			u.Nickname = chooseNickname(username, firstName, lastName)
		}
		return u, nil
	}

	id := r.nextID
	r.nextID++

	u := &models.User{
		ID:           id,
		TelegramID:   tgID,
		Nickname:     chooseNickname(username, firstName, lastName),
		Age:          25, // мок
		FirstName:    firstName,
		LastName:     lastName,
		Username:     username,
		LanguageCode: languageCode,
	}

	r.byID[id] = u
	r.byTgID[tgID] = id
	return u, nil
}

func (r *UserMemoryRepository) FindByID(id int64) (*models.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	u, ok := r.byID[id]
	if !ok {
		return nil, nil
	}
	return u, nil
}

func chooseNickname(username, firstName, lastName string) string {
	if username != "" {
		return username
	}
	if firstName != "" {
		return firstName
	}
	if lastName != "" {
		return lastName
	}
	return "user"
}
