package model

import (
	"time"

	"github.com/google/uuid"
)

type ApiKey struct {
	ID        uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID    uuid.UUID  `gorm:"type:uuid;not null;index" json:"user_id"`
	Name      string     `gorm:"type:varchar(128);not null" json:"name"`
	KeyHash   string     `gorm:"type:varchar(255);not null;index" json:"-"`
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
	CreatedAt time.Time  `gorm:"not null;default:now()" json:"created_at"`

	User *User `gorm:"foreignKey:UserID" json:"-"`
}

func (ApiKey) TableName() string {
	return "api_keys"
}

type CreateApiKeyRequest struct {
	Name      string `json:"name" binding:"required,min=3,max=128"`
	ExpiresIn int    `json:"expires_in" binding:"omitempty,min=1"`
}

type ApiKeyResponse struct {
	ID        uuid.UUID  `json:"id"`
	Name      string     `json:"name"`
	ApiKey    string     `json:"api_key,omitempty"`
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
}
