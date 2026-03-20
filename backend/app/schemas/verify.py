from datetime import datetime
from typing import Optional, List
from pydantic import BaseModel, Field
from uuid import UUID


class VerifyRequest(BaseModel):
    package_id: UUID
    password: str = Field(..., min_length=1)


class VerifyResponse(BaseModel):
    valid: bool
    key: Optional[str] = None
    expires_at: Optional[datetime] = None


class VerifyErrorResponse(BaseModel):
    valid: bool = False
    message: Optional[str] = None
    remaining_attempts: Optional[int] = None
    valid_from: Optional[datetime] = None


class BatchVerifyRequest(BaseModel):
    package_id: UUID
    passwords: List[str] = Field(..., min_length=1)


class BatchVerifyResponse(BaseModel):
    valid: bool
    matched_password: Optional[str] = None
    key: Optional[str] = None


class PackageStatusResponse(BaseModel):
    package_id: UUID
    status: str
    current_password_count: int
    next_password_change: Optional[datetime] = None
    offline_mode_available: bool = True
