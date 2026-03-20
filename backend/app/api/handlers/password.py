from fastapi import APIRouter, Depends, HTTPException, status
from sqlalchemy.ext.asyncio import AsyncSession
from sqlalchemy import select, func
from typing import List
import uuid
import hashlib

from app.core.database import get_db
from app.models.password import PasswordPolicy, PasswordStatus
from app.models.audit import AuditLog, AuditAction
from app.schemas.password import (
    CreatePasswordRequest, BatchCreatePasswordRequest,
    UpdatePasswordRequest, PasswordResponse, CurrentPasswordResponse
)
from app.core.security import get_password_hash

router = APIRouter(prefix="/passwords", tags=["passwords"])


@router.post("/packages/{package_id}/passwords", response_model=PasswordResponse, status_code=status.HTTP_201_CREATED)
async def create_password(
    package_id: uuid.UUID,
    request: CreatePasswordRequest,
    db: AsyncSession = Depends(get_db),
):
    password_hash = get_password_hash(request.password)
    password_key_hash = hashlib.sha256(request.password.encode('utf-8')).hexdigest()[:16]
    package_id_str = str(package_id)
    
    policy = PasswordPolicy(
        id=str(uuid.uuid4()),
        package_id=package_id_str,
        password_hash=password_hash,
        password_key_hash=password_key_hash,
        priority=request.priority,
        valid_from=request.valid_from,
        valid_until=request.valid_until,
        status=PasswordStatus.ACTIVE.value,
    )
    
    db.add(policy)
    
    audit_log = AuditLog(
        id=str(uuid.uuid4()),
        action=AuditAction.POLICY_CREATE.value,
        package_id=package_id_str,
    )
    db.add(audit_log)
    
    await db.commit()
    await db.refresh(policy)
    
    return PasswordResponse(
        id=policy.id,
        password=request.password,
        priority=policy.priority,
        valid_from=policy.valid_from,
        valid_until=policy.valid_until,
        status=policy.status,
        created_at=policy.created_at,
    )


@router.post("/packages/{package_id}/passwords/batch", status_code=status.HTTP_201_CREATED)
async def batch_create_passwords(
    package_id: uuid.UUID,
    request: BatchCreatePasswordRequest,
    db: AsyncSession = Depends(get_db),
):
    package_id_str = str(package_id)
    policies = []
    
    for pwd_req in request.passwords:
        password_hash = get_password_hash(pwd_req.password)
        password_key_hash = hashlib.sha256(pwd_req.password.encode('utf-8')).hexdigest()[:16]
        
        policy = PasswordPolicy(
            id=str(uuid.uuid4()),
            package_id=package_id_str,
            password_hash=password_hash,
            password_key_hash=password_key_hash,
            priority=pwd_req.priority,
            valid_from=pwd_req.valid_from,
            valid_until=pwd_req.valid_until,
            status=PasswordStatus.ACTIVE.value,
        )
        policies.append(policy)
    
    db.add_all(policies)
    
    audit_log = AuditLog(
        id=str(uuid.uuid4()),
        action=AuditAction.POLICY_CREATE.value,
        package_id=package_id_str,
        detail={"count": len(policies)}
    )
    db.add(audit_log)
    
    await db.commit()
    
    return {"created_count": len(policies)}


@router.get("/packages/{package_id}/passwords", response_model=List[PasswordResponse])
async def list_passwords(
    package_id: uuid.UUID,
    db: AsyncSession = Depends(get_db),
):
    result = await db.execute(
        select(PasswordPolicy)
        .where(PasswordPolicy.package_id == str(package_id))
        .order_by(PasswordPolicy.priority)
    )
    passwords = result.scalars().all()
    
    return [
        PasswordResponse(
            id=pwd.id,
            priority=pwd.priority,
            valid_from=pwd.valid_from,
            valid_until=pwd.valid_until,
            status=pwd.status,
            created_at=pwd.created_at,
        )
        for pwd in passwords
    ]


@router.get("/packages/{package_id}/passwords/current", response_model=CurrentPasswordResponse)
async def get_current_password(
    package_id: uuid.UUID,
    db: AsyncSession = Depends(get_db),
):
    from datetime import datetime
    
    now = datetime.utcnow()
    result = await db.execute(
        select(PasswordPolicy)
        .where(
            PasswordPolicy.package_id == str(package_id),
            PasswordPolicy.status == PasswordStatus.ACTIVE.value,
        )
        .order_by(PasswordPolicy.priority)
        .limit(1)
    )
    password = result.scalar_one_or_none()
    
    if not password:
        raise HTTPException(status_code=status.HTTP_404_NOT_FOUND, detail="no active password")
    
    return CurrentPasswordResponse(
        id=password.id,
        password="***",
        valid_from=password.valid_from,
        valid_until=password.valid_until,
        status=password.status,
    )


@router.get("/{password_id}", response_model=PasswordResponse)
async def get_password(
    password_id: uuid.UUID,
    db: AsyncSession = Depends(get_db),
):
    result = await db.execute(
        select(PasswordPolicy).where(PasswordPolicy.id == str(password_id))
    )
    password = result.scalar_one_or_none()
    
    if not password:
        raise HTTPException(status_code=status.HTTP_404_NOT_FOUND, detail="password not found")
    
    return PasswordResponse(
        id=password.id,
        priority=password.priority,
        valid_from=password.valid_from,
        valid_until=password.valid_until,
        status=password.status,
        created_at=password.created_at,
    )


@router.put("/{password_id}", response_model=PasswordResponse)
async def update_password(
    password_id: uuid.UUID,
    request: UpdatePasswordRequest,
    db: AsyncSession = Depends(get_db),
):
    result = await db.execute(
        select(PasswordPolicy).where(PasswordPolicy.id == str(password_id))
    )
    password = result.scalar_one_or_none()
    
    if not password:
        raise HTTPException(status_code=status.HTTP_404_NOT_FOUND, detail="password not found")
    
    if request.password:
        password.password_hash = get_password_hash(request.password)
    if request.priority is not None:
        password.priority = request.priority
    if request.valid_from is not None:
        password.valid_from = request.valid_from
    if request.valid_until is not None:
        password.valid_until = request.valid_until
    
    audit_log = AuditLog(
        id=str(uuid.uuid4()),
        action=AuditAction.POLICY_UPDATE.value,
        package_id=password.package_id,
        detail={"password_id": str(password_id)}
    )
    db.add(audit_log)
    
    await db.commit()
    await db.refresh(password)
    
    return PasswordResponse(
        id=password.id,
        priority=password.priority,
        valid_from=password.valid_from,
        valid_until=password.valid_until,
        status=password.status,
        created_at=password.created_at,
    )


@router.delete("/{password_id}", status_code=status.HTTP_204_NO_CONTENT)
async def delete_password(
    password_id: uuid.UUID,
    db: AsyncSession = Depends(get_db),
):
    result = await db.execute(
        select(PasswordPolicy).where(PasswordPolicy.id == str(password_id))
    )
    password = result.scalar_one_or_none()
    
    if not password:
        raise HTTPException(status_code=status.HTTP_404_NOT_FOUND, detail="password not found")
    
    package_id = password.package_id
    
    await db.delete(password)
    
    audit_log = AuditLog(
        id=str(uuid.uuid4()),
        action=AuditAction.POLICY_DELETE.value,
        package_id=package_id,
    )
    db.add(audit_log)
    
    await db.commit()


@router.post("/{password_id}/activate", response_model=PasswordResponse)
async def activate_password(
    password_id: uuid.UUID,
    db: AsyncSession = Depends(get_db),
):
    result = await db.execute(
        select(PasswordPolicy).where(PasswordPolicy.id == str(password_id))
    )
    password = result.scalar_one_or_none()
    
    if not password:
        raise HTTPException(status_code=status.HTTP_404_NOT_FOUND, detail="password not found")
    
    password.status = PasswordStatus.ACTIVE.value
    
    audit_log = AuditLog(
        id=str(uuid.uuid4()),
        action=AuditAction.POLICY_ACTIVATE.value,
        package_id=password.package_id,
    )
    db.add(audit_log)
    
    await db.commit()
    await db.refresh(password)
    
    return PasswordResponse(
        id=password.id,
        priority=password.priority,
        valid_from=password.valid_from,
        valid_until=password.valid_until,
        status=password.status,
        created_at=password.created_at,
    )


@router.post("/{password_id}/deactivate", response_model=PasswordResponse)
async def deactivate_password(
    password_id: uuid.UUID,
    db: AsyncSession = Depends(get_db),
):
    result = await db.execute(
        select(PasswordPolicy).where(PasswordPolicy.id == str(password_id))
    )
    password = result.scalar_one_or_none()
    
    if not password:
        raise HTTPException(status_code=status.HTTP_404_NOT_FOUND, detail="password not found")
    
    password.status = PasswordStatus.DISABLED.value
    
    audit_log = AuditLog(
        id=str(uuid.uuid4()),
        action=AuditAction.POLICY_DEACTIVATE.value,
        package_id=password.package_id,
    )
    db.add(audit_log)
    
    await db.commit()
    await db.refresh(password)
    
    return PasswordResponse(
        id=password.id,
        priority=password.priority,
        valid_from=password.valid_from,
        valid_until=password.valid_until,
        status=password.status,
        created_at=password.created_at,
    )
