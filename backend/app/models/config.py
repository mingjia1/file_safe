import uuid
from datetime import datetime
from sqlalchemy import Column, Integer, DateTime, Boolean, ForeignKey
from sqlalchemy.dialects.postgresql import UUID

from app.core.database import Base


class EncryptionConfig(Base):
    __tablename__ = "encryption_configs"

    id = Column(UUID(as_uuid=True), primary_key=True, default=uuid.uuid4)
    aes_key_length = Column(Integer, nullable=False, default=256)
    rsa_key_length = Column(Integer, nullable=False, default=2048)
    password_min_length = Column(Integer, nullable=False, default=8)
    password_require_special = Column(Boolean, nullable=False, default=True)
    password_require_uppercase = Column(Boolean, nullable=False, default=True)
    password_require_lowercase = Column(Boolean, nullable=False, default=True)
    password_require_digit = Column(Boolean, nullable=False, default=True)
    config_encrypt = Column(Boolean, nullable=False, default=True)
    enable_signature = Column(Boolean, nullable=False, default=True)
    updated_at = Column(DateTime, nullable=False, default=datetime.utcnow)
    updated_by = Column(UUID(as_uuid=True), ForeignKey("users.id"), nullable=True)
