package services

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
)

var (
	ErrInvalidInitData = errors.New("invalid init data")
	ErrExpiredAuth     = errors.New("auth data expired")
)

type TelegramUser struct {
	ID           int64  `json:"id"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name,omitempty"`
	Username     string `json:"username,omitempty"`
	LanguageCode string `json:"language_code,omitempty"`
}

type VerifiedInitData struct {
	AuthDate time.Time
	User     TelegramUser
	Raw      map[string]string
}

type TelegramAuthService struct {
	BotToken string
	MaxAge   time.Duration
}

func (s TelegramAuthService) Verify(initData string) (*VerifiedInitData, error) {
	values, err := url.ParseQuery(initData)
	if err != nil {
		return nil, ErrInvalidInitData
	}

	receivedHash := values.Get("hash")
	if receivedHash == "" {
		return nil, ErrInvalidInitData
	}
	values.Del("hash")

	keys := make([]string, 0, len(values))
	for k := range values {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	lines := make([]string, 0, len(keys))
	for _, k := range keys {
		lines = append(lines, fmt.Sprintf("%s=%s", k, values.Get(k)))
	}
	dataCheckString := strings.Join(lines, "\n")

	mac := hmac.New(sha256.New, []byte("WebAppData"))
	_, _ = mac.Write([]byte(s.BotToken))
	secretKey := mac.Sum(nil)

	mac2 := hmac.New(sha256.New, secretKey)
	_, _ = mac2.Write([]byte(dataCheckString))
	expectedHash := hex.EncodeToString(mac2.Sum(nil))

	if !hmac.Equal([]byte(expectedHash), []byte(receivedHash)) {
		return nil, ErrInvalidInitData
	}

	authDateStr := values.Get("auth_date")
	if authDateStr == "" {
		return nil, ErrInvalidInitData
	}

	authUnix, err := strconv.ParseInt(authDateStr, 10, 64)
	if err != nil {
		return nil, ErrInvalidInitData
	}
	authTime := time.Unix(authUnix, 0)

	if s.MaxAge > 0 && time.Since(authTime) > s.MaxAge {
		return nil, ErrExpiredAuth
	}

	var user TelegramUser
	if err := json.Unmarshal([]byte(values.Get("user")), &user); err != nil {
		return nil, ErrInvalidInitData
	}

	raw := make(map[string]string, len(values))
	for k := range values {
		raw[k] = values.Get(k)
	}

	return &VerifiedInitData{
		AuthDate: authTime,
		User:     user,
		Raw:      raw,
	}, nil
}
