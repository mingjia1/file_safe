package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"password-timer-manager/internal/model"
	ptmerrors "password-timer-manager/internal/pkg/errors"
	"password-timer-manager/internal/repository"
)

type VerifyService struct {
	passwordRepo *repository.PasswordRepository
	packageRepo  *repository.PackageRepository
	auditRepo    *repository.AuditRepository
}

func NewVerifyService(
	passwordRepo *repository.PasswordRepository,
	packageRepo *repository.PackageRepository,
	auditRepo *repository.AuditRepository,
) *VerifyService {
	return &VerifyService{
		passwordRepo: passwordRepo,
		packageRepo:  packageRepo,
		auditRepo:    auditRepo,
	}
}

func (s *VerifyService) VerifyPassword(ctx context.Context, packageID string, password string) (*model.VerifyResponse, error) {
	pkgUUID, err := uuid.Parse(packageID)
	if err != nil {
		return nil, ptmerrors.ErrPackageNotFound
	}

	pkg, err := s.packageRepo.GetByID(ctx, pkgUUID)
	if err != nil {
		return nil, err
	}

	if pkg.Status != model.PackageStatusActive {
		return nil, ptmerrors.ErrPackageNotFound
	}

	passwords, err := s.passwordRepo.GetByPackageID(ctx, pkgUUID)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	for _, pwdPolicy := range passwords {
		if !s.isPasswordActive(&pwdPolicy, now) {
			continue
		}

		if err := bcrypt.CompareHashAndPassword([]byte(pwdPolicy.PasswordHash), []byte(password)); err == nil {
			s.auditRepo.Create(ctx, &model.AuditLog{
				ID:        uuid.New(),
				Action:    model.AuditActionVerifySuccess,
				PackageID: &pkgUUID,
				CreatedAt: now,
			})

			return &model.VerifyResponse{
				Valid:     true,
				Key:       s.generateDecryptionKey(pkgUUID, pwdPolicy.ID),
				ExpiresAt: now.Add(5 * time.Minute),
			}, nil
		}
	}

	s.auditRepo.Create(ctx, &model.AuditLog{
		ID:        uuid.New(),
		Action:    model.AuditActionVerifyFail,
		PackageID: &pkgUUID,
		CreatedAt: now,
	})

	return &model.VerifyResponse{
		Valid: false,
	}, ptmerrors.ErrPasswordIncorrect
}

func (s *VerifyService) isPasswordActive(policy *model.PasswordPolicy, now time.Time) bool {
	if policy.Status == model.PasswordStatusDisabled {
		return false
	}

	if policy.ValidFrom != nil && now.Before(*policy.ValidFrom) {
		return false
	}

	if policy.ValidUntil != nil && now.After(*policy.ValidUntil) {
		return false
	}

	return true
}

func (s *VerifyService) generateDecryptionKey(packageID, passwordID uuid.UUID) string {
	return packageID.String() + ":" + passwordID.String()
}

func (s *VerifyService) BatchVerify(ctx context.Context, packageID string, passwords []string) (*model.BatchVerifyResponse, error) {
	for _, pwd := range passwords {
		result, err := s.VerifyPassword(ctx, packageID, pwd)
		if err == nil && result.Valid {
			return &model.BatchVerifyResponse{
				Valid:           true,
				MatchedPassword: pwd,
				Key:             result.Key,
			}, nil
		}
	}

	return &model.BatchVerifyResponse{
		Valid: false,
	}, ptmerrors.ErrPasswordIncorrect
}

func (s *VerifyService) GetPackageStatus(ctx context.Context, packageID string) (*model.PackageStatusResponse, error) {
	pkgUUID, err := uuid.Parse(packageID)
	if err != nil {
		return nil, ptmerrors.ErrPackageNotFound
	}

	pkg, err := s.packageRepo.GetByID(ctx, pkgUUID)
	if err != nil {
		return nil, err
	}

	activeCount, err := s.passwordRepo.CountActiveByPackageID(ctx, pkgUUID)
	if err != nil {
		return nil, err
	}

	status := "active"
	if pkg.Status == model.PackageStatusArchived {
		status = "archived"
	} else if activeCount == 0 {
		status = "no_active_password"
	}

	var nextChange *time.Time
	if activeCount > 0 {
		if activePolicy, err := s.passwordRepo.GetActiveByPackageID(ctx, pkgUUID); err == nil {
			nextChange = activePolicy.ValidUntil
		}
	}

	return &model.PackageStatusResponse{
		PackageID:            pkgUUID,
		Status:               status,
		CurrentPasswordCount: activeCount,
		NextPasswordChange:   nextChange,
		OfflineModeAvailable: true,
	}, nil
}
