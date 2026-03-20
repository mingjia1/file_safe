package repository

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"password-timer-manager/internal/model"
)

type ConfigRepository struct {
	db *gorm.DB
}

func NewConfigRepository(db *gorm.DB) *ConfigRepository {
	return &ConfigRepository{db: db}
}

func (r *ConfigRepository) GetEncryptionConfig(ctx context.Context) (*model.EncryptionConfig, error) {
	var config model.EncryptionConfig
	err := r.db.WithContext(ctx).Order("updated_at DESC").First(&config).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return r.createDefaultEncryptionConfig(ctx)
		}
		return nil, err
	}
	return &config, nil
}

func (r *ConfigRepository) createDefaultEncryptionConfig(ctx context.Context) (*model.EncryptionConfig, error) {
	config := &model.EncryptionConfig{
		ID:                       uuid.New(),
		AESKeyLength:             256,
		RSAKeyLength:             2048,
		PasswordMinLength:        8,
		PasswordRequireSpecial:   true,
		PasswordRequireUppercase: true,
		PasswordRequireLowercase: true,
		PasswordRequireDigit:     true,
		ConfigEncrypt:            true,
		EnableSignature:          true,
	}

	if err := r.db.WithContext(ctx).Create(config).Error; err != nil {
		return nil, err
	}

	return config, nil
}

func (r *ConfigRepository) UpdateEncryptionConfig(ctx context.Context, config *model.EncryptionConfig, userID uuid.UUID) error {
	config.UpdatedBy = &userID
	return r.db.WithContext(ctx).Save(config).Error
}
