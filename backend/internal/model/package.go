package model

import (
	"time"

	"github.com/google/uuid"
)

type PackageFormat string

const (
	FormatEXE PackageFormat = "exe"
	FormatZIP PackageFormat = "zip"
)

type PackageStatus string

const (
	PackageStatusActive   PackageStatus = "active"
	PackageStatusArchived PackageStatus = "archived"
)

type FilePackage struct {
	ID          uuid.UUID     `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Name        string        `gorm:"type:varchar(255);not null" json:"name"`
	Format      PackageFormat `gorm:"type:varchar(16);not null" json:"format"`
	Description string        `gorm:"type:text" json:"description,omitempty"`
	Status      PackageStatus `gorm:"type:varchar(32);not null;default:'active'" json:"status"`
	FilePath    string        `gorm:"type:varchar(512);not null" json:"-"`
	FileHash    string        `gorm:"type:varchar(128);not null" json:"file_hash"`
	FileSize    int64         `gorm:"not null" json:"file_size"`
	CreatedBy   uuid.UUID     `gorm:"type:uuid;not null" json:"created_by"`
	CreatedAt   time.Time     `gorm:"not null;default:now()" json:"created_at"`
	UpdatedAt   *time.Time    `json:"updated_at,omitempty"`

	Creator   *User            `gorm:"foreignKey:CreatedBy" json:"creator,omitempty"`
	Passwords []PasswordPolicy `gorm:"foreignKey:PackageID" json:"passwords,omitempty"`
}

func (FilePackage) TableName() string {
	return "file_packages"
}

type CreatePackageRequest struct {
	Name        string        `json:"name" binding:"required"`
	Format      PackageFormat `json:"format" binding:"required,oneof=exe zip"`
	Description string        `json:"description"`
}

type UpdatePackageRequest struct {
	Name        *string        `json:"name"`
	Description *string        `json:"description"`
	Status      *PackageStatus `json:"status" binding:"omitempty,oneof=active archived"`
}

type PackageResponse struct {
	ID              uuid.UUID          `json:"id"`
	Name            string             `json:"name"`
	Format          PackageFormat      `json:"format"`
	Status          PackageStatus      `json:"status"`
	Description     string             `json:"description,omitempty"`
	FileSize        int64              `json:"file_size"`
	FileHash        string             `json:"file_hash"`
	CreatedAt       time.Time          `json:"created_at"`
	UpdatedAt       *time.Time         `json:"updated_at,omitempty"`
	CreatedBy       string             `json:"created_by"`
	PasswordCount   int                `json:"password_count"`
	CurrentPassword string             `json:"current_password,omitempty"`
	Passwords       []PasswordResponse `json:"passwords,omitempty"`
}

type PackageListRequest struct {
	Page     int           `form:"page" binding:"omitempty,min=1"`
	PageSize int           `form:"page_size" binding:"omitempty,min=1,max=100"`
	Status   PackageStatus `form:"status" binding:"omitempty,oneof=active archived"`
	Format   PackageFormat `form:"format" binding:"omitempty,oneof=exe zip"`
}

type DownloadURLResponse struct {
	DownloadURL string    `json:"download_url"`
	ExpiresAt   time.Time `json:"expires_at"`
}
