from datetime import datetime
from sqlalchemy import Column, String, DateTime, Text, Index
from sqlalchemy import JSON
import enum

from app.core.database import Base


class AuditAction(str, enum.Enum):
    DOWNLOAD = "DOWNLOAD"
    VERIFY_SUCCESS = "VERIFY_SUCCESS"
    VERIFY_FAIL = "VERIFY_FAIL"
    POLICY_CREATE = "POLICY_CREATE"
    POLICY_UPDATE = "POLICY_UPDATE"
    POLICY_DELETE = "POLICY_DELETE"
    POLICY_ACTIVATE = "POLICY_ACTIVATE"
    POLICY_DEACTIVATE = "POLICY_DEACTIVATE"
    PACKAGE_CREATE = "PACKAGE_CREATE"
    PACKAGE_UPDATE = "PACKAGE_UPDATE"
    PACKAGE_DELETE = "PACKAGE_DELETE"
    USER_LOGIN = "USER_LOGIN"
    USER_LOGOUT = "USER_LOGOUT"
    USER_CREATE = "USER_CREATE"
    USER_UPDATE = "USER_UPDATE"
    USER_DELETE = "USER_DELETE"
    CONFIG_UPDATE = "CONFIG_UPDATE"


class AuditLog(Base):
    __tablename__ = "audit_logs"

    id = Column(String(36), primary_key=True)
    action = Column(String(64), nullable=False, index=True)
    package_id = Column(String(36), nullable=True, index=True)
    user_id = Column(String(36), nullable=True, index=True)
    ip_address = Column(String(64), nullable=True)
    user_agent = Column(Text, nullable=True)
    detail = Column(JSON, nullable=True)
    created_at = Column(DateTime, nullable=False, default=datetime.utcnow, index=True)

    __table_args__ = (
        Index("idx_audit_logs_package_user", "package_id", "user_id"),
    )
