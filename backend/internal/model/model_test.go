package model

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPasswordPolicy_CalculateStatus(t *testing.T) {
	now := time.Now()
	past := now.Add(-24 * time.Hour)
	future := now.Add(24 * time.Hour)

	tests := []struct {
		name     string
		policy   PasswordPolicy
		expected PasswordStatus
	}{
		{
			name: "disabled password returns disabled",
			policy: PasswordPolicy{
				Status: PasswordStatusDisabled,
			},
			expected: PasswordStatusDisabled,
		},
		{
			name: "valid_from in future returns pending",
			policy: PasswordPolicy{
				Status:    PasswordStatusActive,
				ValidFrom: &future,
			},
			expected: PasswordStatusPending,
		},
		{
			name: "valid_until in past returns expired",
			policy: PasswordPolicy{
				Status:     PasswordStatusActive,
				ValidUntil: &past,
			},
			expected: PasswordStatusExpired,
		},
		{
			name: "within valid range returns active",
			policy: PasswordPolicy{
				Status:     PasswordStatusPending,
				ValidFrom:  &past,
				ValidUntil: &future,
			},
			expected: PasswordStatusActive,
		},
		{
			name: "no time constraints returns active",
			policy: PasswordPolicy{
				Status:     PasswordStatusPending,
				ValidFrom:  nil,
				ValidUntil: nil,
			},
			expected: PasswordStatusActive,
		},
		{
			name: "valid_until exactly future returns active",
			policy: PasswordPolicy{
				Status:     PasswordStatusActive,
				ValidFrom:  nil,
				ValidUntil: &future,
			},
			expected: PasswordStatusActive,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.policy.CalculateStatus()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestUser_HasPermission(t *testing.T) {
	tests := []struct {
		name     string
		role     UserRole
		perm     Permission
		expected bool
	}{
		{
			name:     "super_admin has all permissions",
			role:     RoleSuperAdmin,
			perm:     PermissionPackageCreate,
			expected: true,
		},
		{
			name:     "super_admin has user management",
			role:     RoleSuperAdmin,
			perm:     PermissionUserManage,
			expected: true,
		},
		{
			name:     "admin has package permissions",
			role:     RoleAdmin,
			perm:     PermissionPackageCreate,
			expected: true,
		},
		{
			name:     "admin does not have user management",
			role:     RoleAdmin,
			perm:     PermissionUserManage,
			expected: false,
		},
		{
			name:     "operator has package read",
			role:     RoleOperator,
			perm:     PermissionPackageRead,
			expected: true,
		},
		{
			name:     "operator does not have audit read",
			role:     RoleOperator,
			perm:     PermissionAuditRead,
			expected: false,
		},
		{
			name:     "viewer has package read",
			role:     RoleViewer,
			perm:     PermissionPackageRead,
			expected: true,
		},
		{
			name:     "viewer does not have package create",
			role:     RoleViewer,
			perm:     PermissionPackageCreate,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := &User{Role: tt.role}
			result := user.HasPermission(tt.perm)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestUser_ToResponse(t *testing.T) {
	now := time.Now()
	user := &User{
		ID:        [16]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16},
		Username:  "testuser",
		Email:     "test@example.com",
		Role:      RoleAdmin,
		Status:    UserStatusActive,
		CreatedAt: now,
		LastLogin: &now,
	}

	response := user.ToResponse()

	assert.Equal(t, user.ID, response.ID)
	assert.Equal(t, user.Username, response.Username)
	assert.Equal(t, user.Email, response.Email)
	assert.Equal(t, user.Role, response.Role)
	assert.Equal(t, user.Status, response.Status)
	assert.Equal(t, user.CreatedAt, response.CreatedAt)
	assert.NotNil(t, response.LastLogin)
}

func TestPasswordPolicy_ToResponse(t *testing.T) {
	now := time.Now()
	policy := &PasswordPolicy{
		ID:         [16]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16},
		Priority:   1,
		ValidFrom:  &now,
		ValidUntil: &now,
		Status:     PasswordStatusActive,
		CreatedAt:  now,
	}

	t.Run("include password", func(t *testing.T) {
		response := policy.ToResponse(true)
		assert.Equal(t, policy.ID, response.ID)
		assert.Equal(t, policy.Priority, response.Priority)
		assert.Equal(t, policy.Status, response.Status)
	})

	t.Run("exclude password", func(t *testing.T) {
		response := policy.ToResponse(false)
		assert.Equal(t, policy.ID, response.ID)
		assert.Empty(t, response.Password)
	})
}
