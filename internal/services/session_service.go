package services

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"backend_job_assistant/internal/models"
)

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrExpiredToken = errors.New("expired token")
)

type SessionService struct {
	Secret []byte
	TTL    time.Duration
}

func (s SessionService) Mint(userID int64) (string, error) {
	sid, err := randomID(16)
	if err != nil {
		return "", err
	}

	claims := models.SessionClaims{
		SessionID: sid,
		UserID:    userID,
		ExpiresAt: time.Now().Add(s.TTL),
	}

	payload, err := json.Marshal(claims)
	if err != nil {
		return "", err
	}

	sig := sign(payload, s.Secret)
	return base64.RawURLEncoding.EncodeToString(payload) + "." + sig, nil
}

func (s SessionService) Verify(token string) (*models.SessionClaims, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 2 {
		return nil, ErrInvalidToken
	}

	payload, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return nil, ErrInvalidToken
	}

	expectedSig := sign(payload, s.Secret)
	if !hmac.Equal([]byte(expectedSig), []byte(parts[1])) {
		return nil, ErrInvalidToken
	}

	var claims models.SessionClaims
	if err := json.Unmarshal(payload, &claims); err != nil {
		return nil, ErrInvalidToken
	}

	if time.Now().After(claims.ExpiresAt) {
		return nil, ErrExpiredToken
	}

	return &claims, nil
}

func sign(payload, secret []byte) string {
	mac := hmac.New(sha256.New, secret)
	_, _ = mac.Write(payload)
	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}

func randomID(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", b), nil
}
