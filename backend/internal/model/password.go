package model

import (
	"time"

	"github.com/google/uuid"
)

type PasswordStatus string

const (
	PasswordStatusPending  PasswordStatus = "pending"
	PasswordStatusActive   PasswordStatus = "active"
	PasswordStatusExpired  PasswordStatus = "expired"
	PasswordStatusDisabled PasswordStatus = "disabled"
)

type PasswordPolicy struct {
	ID           uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	PackageID    uuid.UUID      `gorm:"type:uuid;not null;index" json:"package_id"`
	PasswordHash string         `gorm:"type:varchar(255);not null" json:"-"`
	Priority     int            `gorm:"not null;default:1" json:"priority"`
	ValidFrom    *time.Time     `json:"valid_from,omitempty"`
	ValidUntil   *time.Time     `json:"valid_until,omitempty"`
	Status       PasswordStatus `gorm:"type:varchar(32);not null;default:'pending'" json:"status"`
	CreatedAt    time.Time      `gorm:"not null;default:now()" json:"created_at"`

	Package *FilePackage `gorm:"foreignKey:PackageID" json:"-"`
}

func (PasswordPolicy) TableName() string {
	return "password_policies"
}

type CreatePasswordRequest struct {
	Password   string     `json:"password" binding:"required,min=4"`
	Priority   int        `json:"priority" binding:"omitempty,min=1"`
	ValidFrom  *time.Time `json:"valid_from"`
	ValidUntil *time.Time `json:"valid_until"`
}

type BatchCreatePasswordRequest struct {
	Passwords []CreatePasswordRequest `json:"passwords" binding:"required,min=1,dive"`
}

type UpdatePasswordRequest struct {
	Password   *string    `json:"password"`
	Priority   *int       `json:"priority" binding:"omitempty,min=1"`
	ValidFrom  *time.Time `json:"valid_from"`
	ValidUntil *time.Time `json:"valid_until"`
}

type PasswordResponse struct {
	ID         uuid.UUID      `json:"id"`
	Password   string         `json:"password,omitempty"`
	Priority   int            `json:"priority"`
	ValidFrom  *time.Time     `json:"valid_from,omitempty"`
	ValidUntil *time.Time     `json:"valid_until,omitempty"`
	Status     PasswordStatus `json:"status"`
	CreatedAt  time.Time      `json:"created_at"`
}

func (p *PasswordPolicy) ToResponse(includePassword bool) PasswordResponse {
	resp := PasswordResponse{
		ID:         p.ID,
		Priority:   p.Priority,
		ValidFrom:  p.ValidFrom,
		ValidUntil: p.ValidUntil,
		Status:     p.Status,
		CreatedAt:  p.CreatedAt,
	}
	if includePassword {
		resp.Password = ""
	}
	return resp
}

type CurrentPasswordResponse struct {
	ID         uuid.UUID      `json:"id"`
	Password   string         `json:"password"`
	ValidFrom  *time.Time     `json:"valid_from,omitempty"`
	ValidUntil *time.Time     `json:"valid_until,omitempty"`
	Status     PasswordStatus `json:"status"`
}

func (p *PasswordPolicy) CalculateStatus() PasswordStatus {
	now := time.Now()

	if p.Status == PasswordStatusDisabled {
		return PasswordStatusDisabled
	}

	if p.ValidFrom != nil && now.Before(*p.ValidFrom) {
		return PasswordStatusPending
	}

	if p.ValidUntil != nil && now.After(*p.ValidUntil) {
		return PasswordStatusExpired
	}

	return PasswordStatusActive
}
