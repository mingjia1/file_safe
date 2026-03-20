package model

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type AuditAction string

const (
	AuditActionDownload         AuditAction = "DOWNLOAD"
	AuditActionVerifySuccess    AuditAction = "VERIFY_SUCCESS"
	AuditActionVerifyFail       AuditAction = "VERIFY_FAIL"
	AuditActionPolicyCreate     AuditAction = "POLICY_CREATE"
	AuditActionPolicyUpdate     AuditAction = "POLICY_UPDATE"
	AuditActionPolicyDelete     AuditAction = "POLICY_DELETE"
	AuditActionPolicyActivate   AuditAction = "POLICY_ACTIVATE"
	AuditActionPolicyDeactivate AuditAction = "POLICY_DEACTIVATE"
	AuditActionPackageCreate    AuditAction = "PACKAGE_CREATE"
	AuditActionPackageUpdate    AuditAction = "PACKAGE_UPDATE"
	AuditActionPackageDelete    AuditAction = "PACKAGE_DELETE"
	AuditActionUserLogin        AuditAction = "USER_LOGIN"
	AuditActionUserLogout       AuditAction = "USER_LOGOUT"
	AuditActionUserCreate       AuditAction = "USER_CREATE"
	AuditActionUserUpdate       AuditAction = "USER_UPDATE"
	AuditActionUserDelete       AuditAction = "USER_DELETE"
	AuditActionConfigUpdate     AuditAction = "CONFIG_UPDATE"
)

type AuditLog struct {
	ID        uuid.UUID       `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Action    AuditAction     `gorm:"type:varchar(64);not null;index" json:"action"`
	PackageID *uuid.UUID      `gorm:"type:uuid;index" json:"package_id,omitempty"`
	UserID    *uuid.UUID      `gorm:"type:uuid;index" json:"user_id,omitempty"`
	IPAddress string          `gorm:"type:varchar(64)" json:"ip_address,omitempty"`
	UserAgent string          `gorm:"type:text" json:"user_agent,omitempty"`
	Detail    json.RawMessage `gorm:"type:jsonb" json:"detail,omitempty"`
	CreatedAt time.Time       `gorm:"not null;default:now();index" json:"created_at"`

	Package *FilePackage `gorm:"foreignKey:PackageID" json:"-"`
	User    *User        `gorm:"foreignKey:UserID" json:"-"`
}

func (AuditLog) TableName() string {
	return "audit_logs"
}

type AuditLogResponse struct {
	ID          uuid.UUID       `json:"id"`
	Action      AuditAction     `json:"action"`
	PackageID   *uuid.UUID      `json:"package_id,omitempty"`
	PackageName string          `json:"package_name,omitempty"`
	UserID      *uuid.UUID      `json:"user_id,omitempty"`
	Username    string          `json:"username,omitempty"`
	IPAddress   string          `json:"ip_address,omitempty"`
	UserAgent   string          `json:"user_agent,omitempty"`
	Detail      json.RawMessage `json:"detail,omitempty"`
	CreatedAt   time.Time       `json:"created_at"`
}

type AuditListRequest struct {
	Page      int         `form:"page" binding:"omitempty,min=1"`
	PageSize  int         `form:"page_size" binding:"omitempty,min=1,max=100"`
	Action    AuditAction `form:"action" binding:"omitempty"`
	PackageID string      `form:"package_id" binding:"omitempty,uuid"`
	UserID    string      `form:"user_id" binding:"omitempty,uuid"`
	StartTime *time.Time  `form:"start_time"`
	EndTime   *time.Time  `form:"end_time"`
	IPAddress string      `form:"ip_address"`
}

type VerifySuccessDetail struct {
	PasswordID       string `json:"password_id"`
	PasswordPriority int    `json:"password_priority"`
	VerifyMode       string `json:"verify_mode"`
}

type VerifyFailDetail struct {
	AttemptedPassword string `json:"attempted_password"`
	RemainingAttempts int    `json:"remaining_attempts"`
}

type DownloadDetail struct {
	Format          string `json:"format"`
	FileSize        int64  `json:"file_size"`
	DownloadURLType string `json:"download_url_type"`
}

type PolicyChangeDetail struct {
	PolicyID string         `json:"policy_id"`
	Changes  map[string]any `json:"changes"`
}

type LoginDetail struct {
	LoginType string `json:"login_type"`
	MFAUsed   bool   `json:"mfa_used"`
}
