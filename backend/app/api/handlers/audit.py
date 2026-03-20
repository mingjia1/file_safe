from fastapi import APIRouter, Depends, Query
from sqlalchemy.ext.asyncio import AsyncSession
from sqlalchemy import select, func
from typing import List, Optional
from datetime import datetime
import uuid

from app.core.database import get_db
from app.models.audit import AuditLog, AuditAction
from app.models.user import User
from app.models.package import FilePackage

router = APIRouter(prefix="/audit", tags=["audit"])


@router.get("")
async def list_audit_logs(
    page: int = 1,
    page_size: int = 20,
    action: Optional[str] = None,
    package_id: Optional[str] = None,
    user_id: Optional[str] = None,
    start_time: Optional[datetime] = None,
    end_time: Optional[datetime] = None,
    ip_address: Optional[str] = None,
    db: AsyncSession = Depends(get_db),
):
    query = select(AuditLog).order_by(AuditLog.created_at.desc())
    count_query = select(func.count(AuditLog.id))
    
    if action:
        query = query.where(AuditLog.action == action)
        count_query = count_query.where(AuditLog.action == action)
    if package_id:
        query = query.where(AuditLog.package_id == uuid.UUID(package_id))
        count_query = count_query.where(AuditLog.package_id == uuid.UUID(package_id))
    if user_id:
        query = query.where(AuditLog.user_id == uuid.UUID(user_id))
        count_query = count_query.where(AuditLog.user_id == uuid.UUID(user_id))
    if ip_address:
        query = query.where(AuditLog.ip_address == ip_address)
        count_query = count_query.where(AuditLog.ip_address == ip_address)
    if start_time:
        query = query.where(AuditLog.created_at >= start_time)
        count_query = count_query.where(AuditLog.created_at >= start_time)
    if end_time:
        query = query.where(AuditLog.created_at <= end_time)
        count_query = count_query.where(AuditLog.created_at <= end_time)
    
    total_result = await db.execute(count_query)
    total = total_result.scalar()
    
    query = query.offset((page - 1) * page_size).limit(page_size)
    result = await db.execute(query)
    logs = result.scalars().all()
    
    items = []
    for log in logs:
        package_name = None
        username = None
        
        if log.package_id:
            pkg_result = await db.execute(
                select(FilePackage.name).where(FilePackage.id == log.package_id)
            )
            pkg_name = pkg_result.scalar_one_or_none()
            package_name = pkg_name
        
        if log.user_id:
            user_result = await db.execute(
                select(User.username).where(User.id == log.user_id)
            )
            username = user_result.scalar_one_or_none()
        
        items.append({
            "id": str(log.id),
            "action": log.action,
            "package_id": str(log.package_id) if log.package_id else None,
            "package_name": package_name,
            "user_id": str(log.user_id) if log.user_id else None,
            "username": username,
            "ip_address": log.ip_address,
            "user_agent": log.user_agent,
            "detail": log.detail,
            "created_at": log.created_at.isoformat() if log.created_at else None,
        })
    
    total_pages = (total + page_size - 1) // page_size if total else 0
    
    return {
        "items": items,
        "total": total,
        "page": page,
        "page_size": page_size,
        "total_pages": total_pages,
    }


@router.get("/package/{package_id}")
async def get_package_audit_logs(
    package_id: uuid.UUID,
    page: int = 1,
    page_size: int = 20,
    db: AsyncSession = Depends(get_db),
):
    query = select(AuditLog).where(AuditLog.package_id == package_id).order_by(AuditLog.created_at.desc())
    count_query = select(func.count(AuditLog.id)).where(AuditLog.package_id == package_id)
    
    total_result = await db.execute(count_query)
    total = total_result.scalar()
    
    query = query.offset((page - 1) * page_size).limit(page_size)
    result = await db.execute(query)
    logs = result.scalars().all()
    
    return {
        "items": [{"id": str(log.id), "action": log.action, "detail": log.detail, "created_at": log.created_at.isoformat()} for log in logs],
        "total": total,
        "page": page,
        "page_size": page_size,
    }


@router.get("/verify-fails")
async def get_verify_fail_logs(
    page: int = 1,
    page_size: int = 20,
    db: AsyncSession = Depends(get_db),
):
    query = select(AuditLog).where(
        AuditLog.action == AuditAction.VERIFY_FAIL.value
    ).order_by(AuditLog.created_at.desc())
    
    count_query = select(func.count(AuditLog.id)).where(AuditLog.action == AuditAction.VERIFY_FAIL.value)
    
    total_result = await db.execute(count_query)
    total = total_result.scalar()
    
    query = query.offset((page - 1) * page_size).limit(page_size)
    result = await db.execute(query)
    logs = result.scalars().all()
    
    return {
        "items": [{"id": str(log.id), "package_id": str(log.package_id), "ip_address": log.ip_address, "detail": log.detail, "created_at": log.created_at.isoformat()} for log in logs],
        "total": total,
        "page": page,
        "page_size": page_size,
    }


@router.get("/export")
async def export_audit_logs(
    format: str = "json",
    start_time: Optional[datetime] = None,
    end_time: Optional[datetime] = None,
    db: AsyncSession = Depends(get_db),
):
    query = select(AuditLog).order_by(AuditLog.created_at.desc())
    
    if start_time:
        query = query.where(AuditLog.created_at >= start_time)
    if end_time:
        query = query.where(AuditLog.created_at <= end_time)
    
    result = await db.execute(query)
    logs = result.scalars().all()
    
    if format == "csv":
        csv_content = "id,action,package_id,user_id,ip_address,created_at\n"
        for log in logs:
            csv_content += f"{log.id},{log.action},{log.package_id},{log.user_id},{log.ip_address},{log.created_at}\n"
        return {"content": csv_content, "content_type": "text/csv"}
    
    return {
        "items": [{"id": str(log.id), "action": log.action, "package_id": str(log.package_id), "user_id": str(log.user_id), "ip_address": log.ip_address, "detail": log.detail, "created_at": log.created_at.isoformat()} for log in logs],
    }
