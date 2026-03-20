from datetime import datetime
from sqlalchemy import Column, String, DateTime, Integer, Text, BigInteger, Boolean, ForeignKey
from sqlalchemy.orm import relationship
import enum

from app.core.database import Base


class PackageFormat(str, enum.Enum):
    EXE = "exe"
    ZIP = "zip"


class PackageStatus(str, enum.Enum):
    ACTIVE = "active"
    ARCHIVED = "archived"


class FilePackage(Base):
    __tablename__ = "file_packages"

    id = Column(String(36), primary_key=True)
    name = Column(String(255), nullable=False)
    format = Column(String(16), nullable=False)
    description = Column(Text, nullable=True)
    status = Column(String(32), nullable=False, default=PackageStatus.ACTIVE.value)
    file_path = Column(String(512), nullable=False)
    file_hash = Column(String(128), nullable=False)
    file_size = Column(BigInteger, nullable=False)
    created_by = Column(String(36), ForeignKey("users.id"), nullable=False)
    created_at = Column(DateTime, nullable=False, default=datetime.utcnow)
    updated_at = Column(DateTime, nullable=True)

    creator = relationship("User", back_populates="packages")
    passwords = relationship("PasswordPolicy", back_populates="package", cascade="all, delete-orphan")
