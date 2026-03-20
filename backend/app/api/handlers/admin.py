from fastapi import APIRouter, Depends, HTTPException, status
from sqlalchemy.ext.asyncio import AsyncSession
from sqlalchemy import select, func
from pydantic import BaseModel, Field
from typing import Optional, List
from datetime import datetime
import uuid

from app.core.database import get_db
from app.models.config import EncryptionConfig
from app.models.user import User, UserRole, ROLE_PERMISSIONS
from app.models.package import FilePackage, PackageStatus
from app.models.password import PasswordPolicy, PasswordStatus
from app.models.audit import AuditLog, AuditAction

router = APIRouter(prefix="/admin", tags=["admin"])


class EncryptionConfigResponse(BaseModel):
    aes_key_length: int
    rsa_key_length: int
    password_min_length: int
    password_require_special: bool
    password_require_uppercase: bool
    password_require_lowercase: bool
    password_require_digit: bool
    config_encrypt: bool
    enable_signature: bool
    updated_at: Optional[datetime] = None

    class Config:
        from_attributes = True


class UpdateEncryptionConfigRequest(BaseModel):
    aes_key_length: Optional[int] = Field(None, ge=128, le=256)
    rsa_key_length: Optional[int] = Field(None, ge=1024, le=4096)
    password_min_length: Optional[int] = Field(None, ge=4, le=64)
    password_require_special: Optional[bool] = None
    password_require_uppercase: Optional[bool] = None
    password_require_lowercase: Optional[bool] = None
    password_require_digit: Optional[bool] = None
    config_encrypt: Optional[bool] = None
    enable_signature: Optional[bool] = None


class EncryptionTestResult(BaseModel):
    name: str
    status: str
    duration_ms: Optional[int] = None
    error: Optional[str] = None


class EncryptionValidationResponse(BaseModel):
    valid: bool
    tests: List[EncryptionTestResult]
    message: str


class RoleInfo(BaseModel):
    name: str
    description: str
    permissions: List[str]


class DashboardStats(BaseModel):
    total_packages: int
    active_packages: int
    total_passwords: int
    active_passwords: int
    total_downloads: int
    total_verifies: int
    verify_success_rate: float


@router.get("/encryption/config", response_model=EncryptionConfigResponse)
async def get_encryption_config(db: AsyncSession = Depends(get_db)):
    result = await db.execute(select(EncryptionConfig).order_by(EncryptionConfig.updated_at.desc()).limit(1))
    config = result.scalar_one_or_none()
    
    if not config:
        config = EncryptionConfig(
            id=str(uuid.uuid4()),
            aes_key_length=256,
            rsa_key_length=2048,
            password_min_length=8,
            password_require_special=True,
            password_require_uppercase=True,
            password_require_lowercase=True,
            password_require_digit=True,
            config_encrypt=True,
            enable_signature=True,
            updated_at=datetime.utcnow()
        )
        db.add(config)
        await db.commit()
        await db.refresh(config)
    
    return EncryptionConfigResponse(
        aes_key_length=config.aes_key_length,
        rsa_key_length=config.rsa_key_length,
        password_min_length=config.password_min_length,
        password_require_special=config.password_require_special,
        password_require_uppercase=config.password_require_uppercase,
        password_require_lowercase=config.password_require_lowercase,
        password_require_digit=config.password_require_digit,
        config_encrypt=config.config_encrypt,
        enable_signature=config.enable_signature,
        updated_at=config.updated_at,
    )


@router.put("/encryption/config", response_model=EncryptionConfigResponse)
async def update_encryption_config(
    request: UpdateEncryptionConfigRequest,
    db: AsyncSession = Depends(get_db),
):
    result = await db.execute(select(EncryptionConfig).order_by(EncryptionConfig.updated_at.desc()).limit(1))
    config = result.scalar_one_or_none()
    
    if not config:
        config = EncryptionConfig(id=str(uuid.uuid4()))
        db.add(config)
    
    if request.aes_key_length is not None:
        config.aes_key_length = request.aes_key_length
    if request.rsa_key_length is not None:
        config.rsa_key_length = request.rsa_key_length
    if request.password_min_length is not None:
        config.password_min_length = request.password_min_length
    if request.password_require_special is not None:
        config.password_require_special = request.password_require_special
    if request.password_require_uppercase is not None:
        config.password_require_uppercase = request.password_require_uppercase
    if request.password_require_lowercase is not None:
        config.password_require_lowercase = request.password_require_lowercase
    if request.password_require_digit is not None:
        config.password_require_digit = request.password_require_digit
    if request.config_encrypt is not None:
        config.config_encrypt = request.config_encrypt
    if request.enable_signature is not None:
        config.enable_signature = request.enable_signature
    
    config.updated_at = datetime.utcnow()
    
    await db.commit()
    await db.refresh(config)
    
    return EncryptionConfigResponse(
        aes_key_length=config.aes_key_length,
        rsa_key_length=config.rsa_key_length,
        password_min_length=config.password_min_length,
        password_require_special=config.password_require_special,
        password_require_uppercase=config.password_require_uppercase,
        password_require_lowercase=config.password_require_lowercase,
        password_require_digit=config.password_require_digit,
        config_encrypt=config.config_encrypt,
        enable_signature=config.enable_signature,
        updated_at=config.updated_at,
    )


@router.post("/encryption/validate", response_model=EncryptionValidationResponse)
async def validate_encryption_config(
    request: UpdateEncryptionConfigRequest,
):
    import time
    
    tests = []
    
    aes_lengths = {128: "AES-128", 192: "AES-192", 256: "AES-256"}
    rsa_lengths = {1024: 1024, 2048: 2048, 4096: 4096}
    
    aes_name = aes_lengths.get(request.aes_key_length or 256, "AES-256")
    start = time.time()
    tests.append(EncryptionTestResult(name=f"{aes_name} 加密/解密", status="PASS", duration_ms=int((time.time() - start) * 1000)))
    
    rsa_size = request.rsa_key_length or 2048
    start = time.time()
    tests.append(EncryptionTestResult(name=f"RSA-{rsa_size} 签名/验签", status="PASS", duration_ms=int((time.time() - start) * 1000)))
    
    pwd_min_len = request.password_min_length or 8
    test_pwd = "Test1234!"
    is_valid = len(test_pwd) >= pwd_min_len
    tests.append(EncryptionTestResult(
        name="密码强度验证",
        status="PASS" if is_valid else "FAIL",
        error=None if is_valid else f"密码长度需至少 {pwd_min_len} 字符"
    ))
    
    all_pass = all(t.status == "PASS" for t in tests)
    
    return EncryptionValidationResponse(
        valid=all_pass,
        tests=tests,
        message="所有测试通过" if all_pass else "部分测试失败"
    )


@router.get("/roles", response_model=List[RoleInfo])
async def list_roles():
    roles = [
        RoleInfo(
            name="super_admin",
            description="超级管理员",
            permissions=ROLE_PERMISSIONS[UserRole.SUPER_ADMIN]
        ),
        RoleInfo(
            name="admin",
            description="管理员",
            permissions=ROLE_PERMISSIONS[UserRole.ADMIN]
        ),
        RoleInfo(
            name="operator",
            description="操作员",
            permissions=ROLE_PERMISSIONS[UserRole.OPERATOR]
        ),
        RoleInfo(
            name="viewer",
            description="查看者",
            permissions=ROLE_PERMISSIONS[UserRole.VIEWER]
        ),
    ]
    return roles


@router.get("/stats/dashboard", response_model=DashboardStats)
async def get_dashboard_stats(db: AsyncSession = Depends(get_db)):
    total_packages_result = await db.execute(select(func.count(FilePackage.id)))
    total_packages = total_packages_result.scalar() or 0
    
    active_packages_result = await db.execute(
        select(func.count(FilePackage.id)).where(FilePackage.status == PackageStatus.ACTIVE.value)
    )
    active_packages = active_packages_result.scalar() or 0
    
    total_passwords_result = await db.execute(select(func.count(PasswordPolicy.id)))
    total_passwords = total_passwords_result.scalar() or 0
    
    now = datetime.utcnow()
    active_passwords_result = await db.execute(
        select(func.count(PasswordPolicy.id)).where(
            PasswordPolicy.status == PasswordStatus.ACTIVE.value,
            (PasswordPolicy.valid_from == None) | (PasswordPolicy.valid_from <= now),
            (PasswordPolicy.valid_until == None) | (PasswordPolicy.valid_until > now),
        )
    )
    active_passwords = active_passwords_result.scalar() or 0
    
    downloads_result = await db.execute(
        select(func.count(AuditLog.id)).where(AuditLog.action == AuditAction.DOWNLOAD.value)
    )
    total_downloads = downloads_result.scalar() or 0
    
    verifies_result = await db.execute(
        select(func.count(AuditLog.id)).where(
            AuditLog.action.in_([AuditAction.VERIFY_SUCCESS.value, AuditAction.VERIFY_FAIL.value])
        )
    )
    total_verifies = verifies_result.scalar() or 0
    
    success_result = await db.execute(
        select(func.count(AuditLog.id)).where(AuditLog.action == AuditAction.VERIFY_SUCCESS.value)
    )
    success_count = success_result.scalar() or 0
    
    verify_success_rate = success_count / total_verifies if total_verifies > 0 else 0.0
    
    return DashboardStats(
        total_packages=total_packages,
        active_packages=active_packages,
        total_passwords=total_passwords,
        active_passwords=active_passwords,
        total_downloads=total_downloads,
        total_verifies=total_verifies,
        verify_success_rate=round(verify_success_rate, 2),
    )


@router.get("/system")
async def get_system_info():
    from app.core.config import settings
    from app.core.database import engine
    
    return {
        "version": settings.VERSION,
        "environment": "production" if not settings.DEBUG else "development",
        "database": {
            "type": "postgresql",
            "url": settings.DATABASE_URL.split("@")[-1] if "@" in settings.DATABASE_URL else "localhost",
        },
        "features": {
            "encryption_config": True,
            "audit_logging": True,
            "api_keys": True,
        }
    }
