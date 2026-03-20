from fastapi import APIRouter, Depends, HTTPException, status, UploadFile, File
from fastapi.responses import FileResponse
from sqlalchemy.ext.asyncio import AsyncSession
from sqlalchemy import select, func
from typing import Optional
import os
import uuid
import hashlib

from app.core.database import get_db
from app.core.config import settings
from app.models.user import User
from app.models.package import FilePackage, PackageStatus, PackageFormat
from app.models.password import PasswordPolicy, PasswordStatus
from app.models.audit import AuditLog, AuditAction
from app.schemas.package import (
    CreatePackageRequest, UpdatePackageRequest, 
    PackageResponse, PackageListResponse, DownloadURLResponse
)
from app.core.security import get_password_hash

router = APIRouter(prefix="/packages", tags=["packages"])


async def get_current_user_id() -> str:
    return "current-user-id"


@router.post("", response_model=PackageResponse, status_code=status.HTTP_201_CREATED)
async def create_package(
    name: str,
    format: PackageFormat,
    description: Optional[str] = None,
    file: UploadFile = File(...),
    db: AsyncSession = Depends(get_db),
    user_id: str = Depends(get_current_user_id)
):
    file_content = await file.read()
    file_size = len(file_content)
    file_hash = hashlib.sha256(file_content).hexdigest()
    
    package_id = str(uuid.uuid4())
    file_path = f"{settings.STORAGE_LOCAL_PATH}/packages/{package_id}"
    os.makedirs(os.path.dirname(file_path), exist_ok=True)
    with open(file_path, "wb") as f:
        f.write(file_content)
    
    package = FilePackage(
        id=package_id,
        name=name,
        format=format.value,
        description=description,
        status=PackageStatus.ACTIVE.value,
        file_path=file_path,
        file_hash=file_hash,
        file_size=file_size,
        created_by=user_id if user_id != "current-user-id" else str(uuid.uuid4()),
    )
    
    db.add(package)
    
    audit_log = AuditLog(
        id=str(uuid.uuid4()),
        action=AuditAction.PACKAGE_CREATE.value,
        package_id=package_id,
        created_at=package.created_at,
    )
    db.add(audit_log)
    
    await db.commit()
    await db.refresh(package)
    
    return PackageResponse(
        id=package.id,
        name=package.name,
        format=package.format,
        status=package.status,
        description=package.description,
        file_size=package.file_size,
        file_hash=package.file_hash,
        created_at=package.created_at,
        created_by=package.created_by,
        password_count=0,
    )


@router.get("", response_model=PackageListResponse)
async def list_packages(
    page: int = 1,
    page_size: int = 20,
    status_filter: Optional[PackageStatus] = None,
    format: Optional[PackageFormat] = None,
    db: AsyncSession = Depends(get_db),
):
    query = select(FilePackage)
    count_query = select(func.count(FilePackage.id))
    
    if status_filter:
        query = query.where(FilePackage.status == status_filter.value)
        count_query = count_query.where(FilePackage.status == status_filter.value)
    if format:
        query = query.where(FilePackage.format == format.value)
        count_query = count_query.where(FilePackage.format == format.value)
    
    total_result = await db.execute(count_query)
    total = total_result.scalar() or 0
    
    query = query.offset((page - 1) * page_size).limit(page_size).order_by(FilePackage.created_at.desc())
    result = await db.execute(query)
    packages = result.scalars().all()
    
    items = []
    for pkg in packages:
        pwd_count_result = await db.execute(
            select(func.count(PasswordPolicy.id)).where(PasswordPolicy.package_id == pkg.id)
        )
        pwd_count = pwd_count_result.scalar() or 0
        
        items.append(PackageResponse(
            id=pkg.id,
            name=pkg.name,
            format=pkg.format,
            status=pkg.status,
            description=pkg.description,
            file_size=pkg.file_size,
            file_hash=pkg.file_hash,
            created_at=pkg.created_at,
            updated_at=pkg.updated_at,
            created_by=pkg.created_by,
            password_count=pwd_count,
        ))
    
    total_pages = (total + page_size - 1) // page_size if total > 0 else 0
    
    return PackageListResponse(
        items=items,
        total=total,
        page=page,
        page_size=page_size,
        total_pages=total_pages,
    )


@router.get("/{package_id}", response_model=PackageResponse)
async def get_package(
    package_id: uuid.UUID,
    db: AsyncSession = Depends(get_db),
):
    result = await db.execute(select(FilePackage).where(FilePackage.id == str(package_id)))
    package = result.scalar_one_or_none()
    
    if not package:
        raise HTTPException(status_code=status.HTTP_404_NOT_FOUND, detail="package not found")
    
    pwd_result = await db.execute(
        select(PasswordPolicy).where(PasswordPolicy.package_id == str(package_id)).order_by(PasswordPolicy.priority)
    )
    passwords = pwd_result.scalars().all()
    
    pwd_count_result = await db.execute(
        select(func.count(PasswordPolicy.id)).where(PasswordPolicy.package_id == str(package_id))
    )
    pwd_count = pwd_count_result.scalar() or 0
    
    return PackageResponse(
        id=package.id,
        name=package.name,
        format=package.format,
        status=package.status,
        description=package.description,
        file_size=package.file_size,
        file_hash=package.file_hash,
        created_at=package.created_at,
        updated_at=package.updated_at,
        created_by=package.created_by,
        password_count=pwd_count,
    )


@router.put("/{package_id}", response_model=PackageResponse)
async def update_package(
    package_id: uuid.UUID,
    request: UpdatePackageRequest,
    db: AsyncSession = Depends(get_db),
):
    result = await db.execute(select(FilePackage).where(FilePackage.id == str(package_id)))
    package = result.scalar_one_or_none()
    
    if not package:
        raise HTTPException(status_code=status.HTTP_404_NOT_FOUND, detail="package not found")
    
    if request.name is not None:
        package.name = request.name
    if request.description is not None:
        package.description = request.description
    if request.status is not None:
        package.status = request.status.value
    
    await db.commit()
    await db.refresh(package)
    
    return PackageResponse(
        id=package.id,
        name=package.name,
        format=package.format,
        status=package.status,
        description=package.description,
        file_size=package.file_size,
        file_hash=package.file_hash,
        created_at=package.created_at,
        updated_at=package.updated_at,
        created_by=package.created_by,
        password_count=0,
    )


@router.delete("/{package_id}", status_code=status.HTTP_204_NO_CONTENT)
async def delete_package(
    package_id: uuid.UUID,
    db: AsyncSession = Depends(get_db),
):
    result = await db.execute(select(FilePackage).where(FilePackage.id == str(package_id)))
    package = result.scalar_one_or_none()
    
    if not package:
        raise HTTPException(status_code=status.HTTP_404_NOT_FOUND, detail="package not found")
    
    package_id_str = str(package_id)
    await db.delete(package)
    
    audit_log = AuditLog(
        id=str(uuid.uuid4()),
        action=AuditAction.PACKAGE_DELETE.value,
        package_id=package_id_str,
    )
    db.add(audit_log)
    
    await db.commit()


@router.get("/{package_id}/download")
async def download_package(
    package_id: uuid.UUID,
    db: AsyncSession = Depends(get_db),
):
    result = await db.execute(select(FilePackage).where(FilePackage.id == str(package_id)))
    package = result.scalar_one_or_none()
    
    if not package:
        raise HTTPException(status_code=status.HTTP_404_NOT_FOUND, detail="package not found")
    
    if not os.path.exists(package.file_path):
        raise HTTPException(status_code=status.HTTP_404_NOT_FOUND, detail="file not found")
    
    audit_log = AuditLog(
        id=str(uuid.uuid4()),
        action=AuditAction.DOWNLOAD.value,
        package_id=str(package_id),
    )
    db.add(audit_log)
    await db.commit()
    
    return FileResponse(
        path=package.file_path,
        filename=f"{package.name}.{package.format}",
        media_type="application/octet-stream"
    )


@router.post("/{package_id}/download-url", response_model=DownloadURLResponse)
async def get_download_url(
    package_id: uuid.UUID,
    expires_in: int = 3600,
):
    from datetime import datetime, timedelta
    
    return DownloadURLResponse(
        download_url=f"/api/v1/packages/{package_id}/download",
        expires_at=datetime.utcnow() + timedelta(seconds=expires_in)
    )
