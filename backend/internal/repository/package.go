package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"password-timer-manager/internal/model"
	ptmerrors "password-timer-manager/internal/pkg/errors"
)

type PackageRepository struct {
	db *gorm.DB
}

func NewPackageRepository(db *gorm.DB) *PackageRepository {
	return &PackageRepository{db: db}
}

func (r *PackageRepository) Create(ctx context.Context, pkg *model.FilePackage) error {
	return r.db.WithContext(ctx).Create(pkg).Error
}

func (r *PackageRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.FilePackage, error) {
	var pkg model.FilePackage
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&pkg).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ptmerrors.ErrPackageNotFound
		}
		return nil, err
	}
	return &pkg, nil
}

func (r *PackageRepository) GetByIDWithPasswords(ctx context.Context, id uuid.UUID) (*model.FilePackage, error) {
	var pkg model.FilePackage
	err := r.db.WithContext(ctx).Preload("Passwords").Where("id = ?", id).First(&pkg).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ptmerrors.ErrPackageNotFound
		}
		return nil, err
	}
	return &pkg, nil
}

func (r *PackageRepository) Update(ctx context.Context, pkg *model.FilePackage) error {
	return r.db.WithContext(ctx).Save(pkg).Error
}

func (r *PackageRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("package_id = ?", id).Delete(&model.PasswordPolicy{}).Error; err != nil {
			return err
		}
		return tx.Delete(&model.FilePackage{}, "id = ?", id).Error
	})
}

func (r *PackageRepository) List(ctx context.Context, page, pageSize int, status, format string) ([]model.FilePackage, int64, error) {
	var packages []model.FilePackage
	var total int64

	query := r.db.WithContext(ctx).Model(&model.FilePackage{})

	if status != "" {
		query = query.Where("status = ?", status)
	}
	if format != "" {
		query = query.Where("format = ?", format)
	}

	query.Count(&total)

	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&packages).Error
	if err != nil {
		return nil, 0, err
	}

	return packages, total, nil
}

func (r *PackageRepository) CountByStatus(ctx context.Context, status string) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.FilePackage{}).Where("status = ?", status).Count(&count).Error
	return count, err
}
