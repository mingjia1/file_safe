package repository

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"password-timer-manager/internal/model"
)

type AuditRepository struct {
	db *gorm.DB
}

func NewAuditRepository(db *gorm.DB) *AuditRepository {
	return &AuditRepository{db: db}
}

func (r *AuditRepository) Create(ctx context.Context, log *model.AuditLog) error {
	return r.db.WithContext(ctx).Create(log).Error
}

func (r *AuditRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.AuditLog, error) {
	var auditLog model.AuditLog
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&auditLog).Error
	if err != nil {
		return nil, err
	}
	return &auditLog, nil
}

func (r *AuditRepository) List(ctx context.Context, page, pageSize int, action, packageID, userID, ipAddress string, startTime, endTime *string) ([]model.AuditLog, int64, error) {
	var logs []model.AuditLog
	var total int64

	query := r.db.WithContext(ctx).Model(&model.AuditLog{})

	if action != "" {
		query = query.Where("action = ?", action)
	}
	if packageID != "" {
		query = query.Where("package_id = ?", packageID)
	}
	if userID != "" {
		query = query.Where("user_id = ?", userID)
	}
	if ipAddress != "" {
		query = query.Where("ip_address = ?", ipAddress)
	}
	if startTime != nil {
		query = query.Where("created_at >= ?", *startTime)
	}
	if endTime != nil {
		query = query.Where("created_at <= ?", *endTime)
	}

	query.Count(&total)

	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&logs).Error
	if err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}

func (r *AuditRepository) ListByPackageID(ctx context.Context, packageID uuid.UUID, page, pageSize int) ([]model.AuditLog, int64, error) {
	var logs []model.AuditLog
	var total int64

	query := r.db.WithContext(ctx).Model(&model.AuditLog{}).Where("package_id = ?", packageID)
	query.Count(&total)

	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&logs).Error
	if err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}

func (r *AuditRepository) ListByUserID(ctx context.Context, userID uuid.UUID, page, pageSize int) ([]model.AuditLog, int64, error) {
	var logs []model.AuditLog
	var total int64

	query := r.db.WithContext(ctx).Model(&model.AuditLog{}).Where("user_id = ?", userID)
	query.Count(&total)

	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&logs).Error
	if err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}

func (r *AuditRepository) CountByAction(ctx context.Context, action string) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.AuditLog{}).Where("action = ?", action).Count(&count).Error
	return count, err
}
