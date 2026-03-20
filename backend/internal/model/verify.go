package model

import (
	"time"

	"github.com/google/uuid"
)

type VerifyRequest struct {
	PackageID string `json:"package_id" binding:"required,uuid"`
	Password  string `json:"password" binding:"required"`
}

type VerifyResponse struct {
	Valid     bool      `json:"valid"`
	Key       string    `json:"key,omitempty"`
	ExpiresAt time.Time `json:"expires_at,omitempty"`
}

type VerifyErrorResponse struct {
	Valid             bool       `json:"valid"`
	Message           string     `json:"message,omitempty"`
	RemainingAttempts int        `json:"remaining_attempts,omitempty"`
	ValidFrom         *time.Time `json:"valid_from,omitempty"`
}

type BatchVerifyRequest struct {
	PackageID string   `json:"package_id" binding:"required,uuid"`
	Passwords []string `json:"passwords" binding:"required,min=1"`
}

type BatchVerifyResponse struct {
	Valid           bool   `json:"valid"`
	MatchedPassword string `json:"matched_password,omitempty"`
	Key             string `json:"key,omitempty"`
}

type PackageStatusResponse struct {
	PackageID            uuid.UUID  `json:"package_id"`
	Status               string     `json:"status"`
	CurrentPasswordCount int        `json:"current_password_count"`
	NextPasswordChange   *time.Time `json:"next_password_change,omitempty"`
	OfflineModeAvailable bool       `json:"offline_mode_available"`
}

type OfflinePolicy struct {
	PasswordHash string    `json:"password_hash"`
	ValidFrom    time.Time `json:"valid_from"`
	ValidUntil   time.Time `json:"valid_until"`
	KeyEncrypted string    `json:"key_encrypted"`
}

type OfflineConfigResponse struct {
	Version         int             `json:"version"`
	Algorithm       string          `json:"algorithm"`
	EncryptedConfig string          `json:"encrypted_config"`
	Signature       string          `json:"signature"`
	Policies        []OfflinePolicy `json:"policies"`
}
