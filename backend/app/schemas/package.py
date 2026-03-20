from datetime import datetime
from typing import Optional, List
from pydantic import BaseModel, Field
from uuid import UUID

from app.models.package import PackageFormat, PackageStatus


class CreatePackageRequest(BaseModel):
    name: str = Field(..., min_length=1, max_length=255)
    format: PackageFormat
    description: Optional[str] = None


class UpdatePackageRequest(BaseModel):
    name: Optional[str] = Field(None, min_length=1, max_length=255)
    description: Optional[str] = None
    status: Optional[PackageStatus] = None


class PasswordResponse(BaseModel):
    id: UUID
    password: Optional[str] = None
    priority: int
    valid_from: Optional[datetime] = None
    valid_until: Optional[datetime] = None
    status: str
    created_at: datetime

    class Config:
        from_attributes = True


class PackageResponse(BaseModel):
    id: UUID
    name: str
    format: PackageFormat
    status: PackageStatus
    description: Optional[str] = None
    file_size: int
    file_hash: str
    created_at: datetime
    updated_at: Optional[datetime] = None
    created_by: UUID
    password_count: int = 0
    current_password: Optional[str] = None
    passwords: Optional[List[PasswordResponse]] = None

    class Config:
        from_attributes = True


class PackageListResponse(BaseModel):
    items: List[PackageResponse]
    total: int
    page: int
    page_size: int
    total_pages: int


class DownloadURLResponse(BaseModel):
    download_url: str
    expires_at: datetime
