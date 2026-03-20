package model

import (
	"time"

	"github.com/google/uuid"
)

type EncryptionConfig struct {
	ID                       uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	AESKeyLength             int        `gorm:"not null;default:256" json:"aes_key_length"`
	RSAKeyLength             int        `gorm:"not null;default:2048" json:"rsa_key_length"`
	PasswordMinLength        int        `gorm:"not null;default:8" json:"password_min_length"`
	PasswordRequireSpecial   bool       `gorm:"not null;default:true" json:"password_require_special"`
	PasswordRequireUppercase bool       `gorm:"not null;default:true" json:"password_require_uppercase"`
	PasswordRequireLowercase bool       `gorm:"not null;default:true" json:"password_require_lowercase"`
	PasswordRequireDigit     bool       `gorm:"not null;default:true" json:"password_require_digit"`
	ConfigEncrypt            bool       `gorm:"not null;default:true" json:"config_encrypt"`
	EnableSignature          bool       `gorm:"not null;default:true" json:"enable_signature"`
	UpdatedAt                time.Time  `gorm:"not null;default:now()" json:"updated_at"`
	UpdatedBy                *uuid.UUID `gorm:"type:uuid" json:"updated_by,omitempty"`

	Updater *User `gorm:"foreignKey:UpdatedBy" json:"-"`
}

func (EncryptionConfig) TableName() string {
	return "encryption_configs"
}

type UpdateEncryptionConfigRequest struct {
	AESKeyLength             *int  `json:"aes_key_length" binding:"omitempty,oneof=128 192 256"`
	RSAKeyLength             *int  `json:"rsa_key_length" binding:"omitempty,oneof=1024 2048 4096"`
	PasswordMinLength        *int  `json:"password_min_length" binding:"omitempty,min=4,max=64"`
	PasswordRequireSpecial   *bool `json:"password_require_special"`
	PasswordRequireUppercase *bool `json:"password_require_uppercase"`
	PasswordRequireLowercase *bool `json:"password_require_lowercase"`
	PasswordRequireDigit     *bool `json:"password_require_digit"`
	ConfigEncrypt            *bool `json:"config_encrypt"`
	EnableSignature          *bool `json:"enable_signature"`
}

type ValidateEncryptionConfigRequest struct {
	AESKeyLength           int  `json:"aes_key_length" binding:"required,oneof=128 192 256"`
	RSAKeyLength           int  `json:"rsa_key_length" binding:"required,oneof=1024 2048 4096"`
	PasswordMinLength      int  `json:"password_min_length" binding:"required,min=4,max=64"`
	PasswordRequireSpecial bool `json:"password_require_special"`
	ConfigEncrypt          bool `json:"config_encrypt"`
	EnableSignature        bool `json:"enable_signature"`
}

type EncryptionTestResult struct {
	Name       string `json:"name"`
	Status     string `json:"status"`
	DurationMs int64  `json:"duration_ms,omitempty"`
	Error      string `json:"error,omitempty"`
}

type EncryptionValidationResponse struct {
	Valid   bool                   `json:"valid"`
	Tests   []EncryptionTestResult `json:"tests"`
	Message string                 `json:"message"`
}

type RegenerateKeysRequest struct {
	Reason string `json:"reason"`
}

type RegenerateKeysResponse struct {
	OldKeyExpiresAt time.Time `json:"old_key_expires_at"`
	NewKeyActiveAt  time.Time `json:"new_key_active_at"`
	Message         string    `json:"message"`
}
