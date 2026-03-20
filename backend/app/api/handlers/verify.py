import uuid
from datetime import datetime
from typing import Optional
from uuid import UUID
import bcrypt

from fastapi import APIRouter, Depends, HTTPException, status
from sqlalchemy.ext.asyncio import AsyncSession
from sqlalchemy import select

from app.core.database import get_db
from app.models.package import FilePackage, PackageStatus
from app.models.password import PasswordPolicy, PasswordStatus
from app.models.audit import AuditLog, AuditAction
from app.schemas.verify import VerifyRequest, VerifyResponse, BatchVerifyRequest, BatchVerifyResponse

router = APIRouter(prefix="/verify", tags=["verify"])


def is_password_active(policy: PasswordPolicy) -> bool:
    now = datetime.utcnow()
    
    if policy.status == PasswordStatus.DISABLED.value:
        return False
    
    if policy.valid_from and now < policy.valid_from:
        return False
    
    if policy.valid_until and now > policy.valid_until:
        return False
    
    return True


@router.post("", response_model=VerifyResponse)
async def verify_password(
    request: VerifyRequest,
    db: AsyncSession = Depends(get_db)
):
    result = await db.execute(select(FilePackage).where(FilePackage.id == str(request.package_id)))
    package = result.scalar_one_or_none()
    
    if not package:
        raise HTTPException(
            status_code=status.HTTP_404_NOT_FOUND,
            detail="package not found"
        )
    
    if package.status != PackageStatus.ACTIVE.value:
        raise HTTPException(
            status_code=status.HTTP_404_NOT_FOUND,
            detail="package is not active"
        )
    
    result = await db.execute(
        select(PasswordPolicy).where(PasswordPolicy.package_id == str(request.package_id))
    )
    passwords = result.scalars().all()
    
    now = datetime.utcnow()
    
    for pwd_policy in passwords:
        if not is_password_active(pwd_policy):
            continue
        
        try:
            if bcrypt.checkpw(request.password.encode('utf-8'), pwd_policy.password_hash.encode('utf-8')):
                key = f"{package.id}:{pwd_policy.id}"
                
                audit_log = AuditLog(
                    id=str(uuid.uuid4()),
                    action=AuditAction.VERIFY_SUCCESS.value,
                    package_id=str(package.id),
                    created_at=now,
                    detail={"password_id": str(pwd_policy.id), "verify_mode": "online"}
                )
                db.add(audit_log)
                await db.commit()
                
                return VerifyResponse(
                    valid=True,
                    key=key,
                    expires_at=now
                )
        except Exception:
            continue
    
    audit_log = AuditLog(
        id=str(uuid.uuid4()),
        action=AuditAction.VERIFY_FAIL.value,
        package_id=str(package.id),
        created_at=now,
        detail={"attempted_password": request.password[:10]}
    )
    db.add(audit_log)
    await db.commit()
    
    raise HTTPException(
        status_code=status.HTTP_422_UNPROCESSABLE_ENTITY,
        detail="password is incorrect"
    )


@router.post("/batch", response_model=BatchVerifyResponse)
async def batch_verify_passwords(
    request: BatchVerifyRequest,
    db: AsyncSession = Depends(get_db)
):
    for password in request.passwords:
        verify_req = VerifyRequest(package_id=request.package_id, password=password)
        try:
            result = await verify_password(verify_req, db)
            return BatchVerifyResponse(
                valid=True,
                matched_password=password,
                key=result.key
            )
        except HTTPException:
            continue
    
    return BatchVerifyResponse(valid=False)
