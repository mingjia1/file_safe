package repository

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"password-timer-manager/internal/model"
)

type ApiKeyRepository struct {
	db *gorm.DB
}

func NewApiKeyRepository(db *gorm.DB) *ApiKeyRepository {
	return &ApiKeyRepository{db: db}
}

type ApiKeyWithUser struct {
	model.ApiKey
	User *model.User
}

func (r *ApiKeyRepository) Create(ctx context.Context, apiKey *model.ApiKey) error {
	return r.db.WithContext(ctx).Create(apiKey).Error
}

func (r *ApiKeyRepository) GetByHash(ctx context.Context, keyHash string) (*ApiKeyWithUser, error) {
	var result ApiKeyWithUser
	err := r.db.WithContext(ctx).
		Preload("User").
		Where("key_hash = ?", keyHash).
		First(&result.ApiKey).Error
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *ApiKeyRepository) ListByUserID(ctx context.Context, userID uuid.UUID) ([]model.ApiKey, error) {
	var apiKeys []model.ApiKey
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&apiKeys).Error
	if err != nil {
		return nil, err
	}
	return apiKeys, nil
}

func (r *ApiKeyRepository) Delete(ctx context.Context, userID, keyID uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&model.ApiKey{}, "user_id = ? AND id = ?", userID, keyID).Error
}

func (r *ApiKeyRepository) DeleteExpired(ctx context.Context) error {
	return r.db.WithContext(ctx).Delete(&model.ApiKey{}, "expires_at < NOW()").Error
}
