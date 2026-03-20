package repository

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"password-timer-manager/internal/model"
	ptmerrors "password-timer-manager/internal/pkg/errors"
)

type PasswordRepository struct {
	db *gorm.DB
}

func NewPasswordRepository(db *gorm.DB) *PasswordRepository {
	return &PasswordRepository{db: db}
}

func (r *PasswordRepository) Create(ctx context.Context, policy *model.PasswordPolicy) error {
	return r.db.WithContext(ctx).Create(policy).Error
}

func (r *PasswordRepository) CreateBatch(ctx context.Context, policies []model.PasswordPolicy) error {
	return r.db.WithContext(ctx).Create(&policies).Error
}

func (r *PasswordRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.PasswordPolicy, error) {
	var policy model.PasswordPolicy
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&policy).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ptmerrors.ErrPasswordIncorrect
		}
		return nil, err
	}
	return &policy, nil
}

func (r *PasswordRepository) GetByPackageID(ctx context.Context, packageID uuid.UUID) ([]model.PasswordPolicy, error) {
	var policies []model.PasswordPolicy
	err := r.db.WithContext(ctx).Where("package_id = ?", packageID).Order("priority ASC").Find(&policies).Error
	if err != nil {
		return nil, err
	}
	return policies, nil
}

func (r *PasswordRepository) GetActiveByPackageID(ctx context.Context, packageID uuid.UUID) (*model.PasswordPolicy, error) {
	var policy model.PasswordPolicy
	now := time.Now()
	err := r.db.WithContext(ctx).
		Where("package_id = ? AND status = ? AND (valid_from IS NULL OR valid_from <= ?) AND (valid_until IS NULL OR valid_until > ?)",
			packageID, model.PasswordStatusActive, now, now).
		Order("priority ASC").
		First(&policy).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ptmerrors.ErrNoActivePassword
		}
		return nil, err
	}
	return &policy, nil
}

func (r *PasswordRepository) Update(ctx context.Context, policy *model.PasswordPolicy) error {
	return r.db.WithContext(ctx).Save(policy).Error
}

func (r *PasswordRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&model.PasswordPolicy{}, "id = ?", id).Error
}

func (r *PasswordRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status model.PasswordStatus) error {
	return r.db.WithContext(ctx).Model(&model.PasswordPolicy{}).Where("id = ?", id).Update("status", status).Error
}

func (r *PasswordRepository) CountActiveByPackageID(ctx context.Context, packageID uuid.UUID) (int, error) {
	var count int64
	now := time.Now()
	err := r.db.WithContext(ctx).
		Model(&model.PasswordPolicy{}).
		Where("package_id = ? AND status = ? AND (valid_from IS NULL OR valid_from <= ?) AND (valid_until IS NULL OR valid_until > ?)",
			packageID, model.PasswordStatusActive, now, now).
		Count(&count).Error
	return int(count), err
}
