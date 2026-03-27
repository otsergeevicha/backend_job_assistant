package models

type User struct {
	ID           int64  `json:"id"`
	TelegramID   int64  `json:"telegram_id"`
	Nickname     string `json:"nickname"`
	Age          int    `json:"age"`
	FirstName    string `json:"first_name,omitempty"`
	LastName     string `json:"last_name,omitempty"`
	Username     string `json:"username,omitempty"`
	LanguageCode string `json:"language_code,omitempty"`
}
