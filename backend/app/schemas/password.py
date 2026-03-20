from datetime import datetime
from typing import Optional, List
from pydantic import BaseModel, Field
from uuid import UUID


class CreatePasswordRequest(BaseModel):
    password: str = Field(..., min_length=4)
    priority: int = Field(1, ge=1)
    valid_from: Optional[datetime] = None
    valid_until: Optional[datetime] = None


class BatchCreatePasswordRequest(BaseModel):
    passwords: List[CreatePasswordRequest] = Field(..., min_length=1)


class UpdatePasswordRequest(BaseModel):
    password: Optional[str] = Field(None, min_length=4)
    priority: Optional[int] = Field(None, ge=1)
    valid_from: Optional[datetime] = None
    valid_until: Optional[datetime] = None


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


class CurrentPasswordResponse(BaseModel):
    id: UUID
    password: str
    valid_from: Optional[datetime] = None
    valid_until: Optional[datetime] = None
    status: str
