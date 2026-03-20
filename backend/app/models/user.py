import uuid
from datetime import datetime
from sqlalchemy import Column, String, DateTime, Enum, ForeignKey, Integer, BigInteger, Boolean, Text
from sqlalchemy.dialects.postgresql import UUID, JSONB
from sqlalchemy.orm import relationship
import enum

from app.core.database import Base


class UserRole(str, enum.Enum):
    SUPER_ADMIN = "super_admin"
    ADMIN = "admin"
    OPERATOR = "operator"
    VIEWER = "viewer"


class UserStatus(str, enum.Enum):
    ACTIVE = "active"
    DISABLED = "disabled"


class User(Base):
    __tablename__ = "users"

    id = Column(UUID(as_uuid=True), primary_key=True, default=uuid.uuid4)
    username = Column(String(64), unique=True, nullable=False, index=True)
    email = Column(String(255), unique=True, nullable=False, index=True)
    password_hash = Column(String(255), nullable=False)
    role = Column(String(32), nullable=False, default=UserRole.OPERATOR.value)
    status = Column(String(32), nullable=False, default=UserStatus.ACTIVE.value)
    created_at = Column(DateTime, nullable=False, default=datetime.utcnow)
    last_login = Column(DateTime, nullable=True)

    packages = relationship("FilePackage", back_populates="creator")
    api_keys = relationship("ApiKey", back_populates="user")


ROLE_PERMISSIONS = {
    UserRole.SUPER_ADMIN: [
        "package:create", "package:read", "package:update", "package:delete",
        "password:manage", "password:activate",
        "audit:read", "user:manage", "encryption:manage"
    ],
    UserRole.ADMIN: [
        "package:create", "package:read", "package:update", "package:delete",
        "password:manage", "password:activate",
        "audit:read", "encryption:manage"
    ],
    UserRole.OPERATOR: [
        "package:create", "package:read", "package:update",
        "password:manage", "password:activate"
    ],
    UserRole.VIEWER: ["package:read"],
}


def has_permission(role: UserRole, permission: str) -> bool:
    return permission in ROLE_PERMISSIONS.get(role, [])
