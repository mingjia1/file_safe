package model

import (
	"time"

	"github.com/google/uuid"
)

type UserRole string

const (
	RoleSuperAdmin UserRole = "super_admin"
	RoleAdmin      UserRole = "admin"
	RoleOperator   UserRole = "operator"
	RoleViewer     UserRole = "viewer"
)

type UserStatus string

const (
	UserStatusActive   UserStatus = "active"
	UserStatusDisabled UserStatus = "disabled"
)

type User struct {
	ID           uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Username     string     `gorm:"type:varchar(64);uniqueIndex;not null" json:"username"`
	Email        string     `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"`
	PasswordHash string     `gorm:"type:varchar(255);not null" json:"-"`
	Role         UserRole   `gorm:"type:varchar(32);not null;default:'operator'" json:"role"`
	Status       UserStatus `gorm:"type:varchar(32);not null;default:'active'" json:"status"`
	CreatedAt    time.Time  `gorm:"not null;default:now()" json:"created_at"`
	LastLogin    *time.Time `json:"last_login,omitempty"`
}

func (User) TableName() string {
	return "users"
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8"`
}

type CreateUserRequest struct {
	Username string   `json:"username" binding:"required,min=3,max=64"`
	Email    string   `json:"email" binding:"required,email"`
	Password string   `json:"password" binding:"required,min=8"`
	Role     UserRole `json:"role" binding:"required,oneof=super_admin admin operator viewer"`
}

type UserResponse struct {
	ID        uuid.UUID  `json:"id"`
	Username  string     `json:"username"`
	Email     string     `json:"email"`
	Role      UserRole   `json:"role"`
	Status    UserStatus `json:"status"`
	CreatedAt time.Time  `json:"created_at"`
	LastLogin *time.Time `json:"last_login,omitempty"`
}

func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		Role:      u.Role,
		Status:    u.Status,
		CreatedAt: u.CreatedAt,
		LastLogin: u.LastLogin,
	}
}

type Permission string

const (
	PermissionPackageCreate    Permission = "package:create"
	PermissionPackageRead      Permission = "package:read"
	PermissionPackageUpdate    Permission = "package:update"
	PermissionPackageDelete    Permission = "package:delete"
	PermissionPasswordManage   Permission = "password:manage"
	PermissionPasswordActivate Permission = "password:activate"
	PermissionAuditRead        Permission = "audit:read"
	PermissionUserManage       Permission = "user:manage"
	PermissionEncryptionManage Permission = "encryption:manage"
)

var RolePermissions = map[UserRole][]Permission{
	RoleSuperAdmin: {
		PermissionPackageCreate, PermissionPackageRead, PermissionPackageUpdate, PermissionPackageDelete,
		PermissionPasswordManage, PermissionPasswordActivate,
		PermissionAuditRead, PermissionUserManage, PermissionEncryptionManage,
	},
	RoleAdmin: {
		PermissionPackageCreate, PermissionPackageRead, PermissionPackageUpdate, PermissionPackageDelete,
		PermissionPasswordManage, PermissionPasswordActivate,
		PermissionAuditRead, PermissionEncryptionManage,
	},
	RoleOperator: {
		PermissionPackageCreate, PermissionPackageRead, PermissionPackageUpdate,
		PermissionPasswordManage, PermissionPasswordActivate,
	},
	RoleViewer: {
		PermissionPackageRead,
	},
}

func (u *User) HasPermission(perm Permission) bool {
	perms, ok := RolePermissions[u.Role]
	if !ok {
		return false
	}
	for _, p := range perms {
		if p == perm {
			return true
		}
	}
	return false
}
