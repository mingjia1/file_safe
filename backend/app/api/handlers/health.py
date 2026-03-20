from fastapi import APIRouter, Depends, HTTPException, status
from sqlalchemy.ext.asyncio import AsyncSession

from app.core.database import get_db

router = APIRouter(tags=["health"])


@router.get("/health")
async def health_check():
    return {"status": "healthy", "version": "1.0.0"}


@router.get("/api/v1/health")
async def api_health_check():
    return {"code": 0, "message": "success", "data": {"status": "healthy", "version": "1.0.0"}}
