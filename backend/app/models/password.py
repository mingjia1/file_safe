from datetime import datetime
from sqlalchemy import Column, String, DateTime, Integer, Text, ForeignKey
from sqlalchemy.orm import relationship
import enum

from app.core.database import Base


class PasswordStatus(str, enum.Enum):
    PENDING = "pending"
    ACTIVE = "active"
    EXPIRED = "expired"
    DISABLED = "disabled"


class PasswordPolicy(Base):
    __tablename__ = "password_policies"

    id = Column(String(36), primary_key=True)
    package_id = Column(String(36), ForeignKey("file_packages.id"), nullable=False, index=True)
    password_hash = Column(String(255), nullable=False)
    password_key_hash = Column(String(32), nullable=False)
    priority = Column(Integer, nullable=False, default=1)
    valid_from = Column(DateTime, nullable=True)
    valid_until = Column(DateTime, nullable=True)
    status = Column(String(32), nullable=False, default=PasswordStatus.PENDING.value)
    created_at = Column(DateTime, nullable=False, default=datetime.utcnow)

    package = relationship("FilePackage", back_populates="passwords")

    def calculate_status(self) -> PasswordStatus:
        now = datetime.utcnow()
        
        if self.status == PasswordStatus.DISABLED.value:
            return PasswordStatus.DISABLED
        
        if self.valid_from and now < self.valid_from:
            return PasswordStatus.PENDING
        
        if self.valid_until and now > self.valid_until:
            return PasswordStatus.EXPIRED
        
        return PasswordStatus.ACTIVE
