from datetime import datetime
from sqlalchemy import Column, String, Integer, DateTime, Boolean

from app.core.database import Base


class EncryptionConfig(Base):
    __tablename__ = "encryption_configs"

    id = Column(String(36), primary_key=True)
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
    updated_by = Column(String(36), nullable=True)
