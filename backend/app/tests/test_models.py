import pytest
from datetime import datetime, timedelta
from app.models.user import User, UserRole, has_permission, ROLE_PERMISSIONS
from app.models.password import PasswordPolicy, PasswordStatus


class TestUserRolePermissions:
    def test_super_admin_has_all_permissions(self):
        assert has_permission(UserRole.SUPER_ADMIN, "package:create") is True
        assert has_permission(UserRole.SUPER_ADMIN, "user:manage") is True
        assert has_permission(UserRole.SUPER_ADMIN, "encryption:manage") is True

    def test_admin_has_package_permissions(self):
        assert has_permission(UserRole.ADMIN, "package:create") is True
        assert has_permission(UserRole.ADMIN, "audit:read") is True

    def test_admin_no_user_management(self):
        assert has_permission(UserRole.ADMIN, "user:manage") is False

    def test_operator_has_package_read(self):
        assert has_permission(UserRole.OPERATOR, "package:read") is True
        assert has_permission(UserRole.OPERATOR, "password:manage") is True

    def test_operator_no_audit_read(self):
        assert has_permission(UserRole.OPERATOR, "audit:read") is False

    def test_viewer_only_package_read(self):
        assert has_permission(UserRole.VIEWER, "package:read") is True
        assert has_permission(UserRole.VIEWER, "package:create") is False


class TestPasswordPolicyStatus:
    def test_disabled_password_returns_disabled(self):
        policy = PasswordPolicy(status=PasswordStatus.DISABLED.value)
        assert policy.calculate_status() == PasswordStatus.DISABLED

    def test_future_valid_from_returns_pending(self):
        future = datetime.utcnow() + timedelta(days=1)
        policy = PasswordPolicy(
            status=PasswordStatus.ACTIVE.value,
            valid_from=future
        )
        assert policy.calculate_status() == PasswordStatus.PENDING

    def test_past_valid_until_returns_expired(self):
        past = datetime.utcnow() - timedelta(days=1)
        policy = PasswordPolicy(
            status=PasswordStatus.ACTIVE.value,
            valid_until=past
        )
        assert policy.calculate_status() == PasswordStatus.EXPIRED

    def test_within_valid_range_returns_active(self):
        past = datetime.utcnow() - timedelta(days=1)
        future = datetime.utcnow() + timedelta(days=1)
        policy = PasswordPolicy(
            status=PasswordStatus.PENDING.value,
            valid_from=past,
            valid_until=future
        )
        assert policy.calculate_status() == PasswordStatus.ACTIVE

    def test_no_time_constraints_returns_active(self):
        policy = PasswordPolicy(
            status=PasswordStatus.PENDING.value,
            valid_from=None,
            valid_until=None
        )
        assert policy.calculate_status() == PasswordStatus.ACTIVE
