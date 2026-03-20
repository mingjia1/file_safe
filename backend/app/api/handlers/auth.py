from datetime import datetime
from typing import Optional
from uuid import UUID

from fastapi import APIRouter, Depends, HTTPException, status
from sqlalchemy.ext.asyncio import AsyncSession
from sqlalchemy import select

from app.core.database import get_db
from app.core.security import verify_password, get_password_hash, create_access_token
from app.models.user import User, UserStatus
from app.schemas.auth import LoginRequest, LoginResponse, UserCreate, UserResponse

router = APIRouter(prefix="/auth", tags=["auth"])


@router.post("/login", response_model=LoginResponse)
async def login(request: LoginRequest, db: AsyncSession = Depends(get_db)):
    result = await db.execute(select(User).where(User.username == request.username))
    user = result.scalar_one_or_none()
    
    if not user or not verify_password(request.password, user.password_hash):
        raise HTTPException(
            status_code=status.HTTP_401_UNAUTHORIZED,
            detail="invalid username or password"
        )
    
    if user.status != UserStatus.ACTIVE.value:
        raise HTTPException(
            status_code=status.HTTP_401_UNAUTHORIZED,
            detail="account is disabled"
        )
    
    token = create_access_token({
        "sub": str(user.id),
        "username": user.username,
        "role": user.role
    })
    
    from datetime import timedelta
    from app.core.config import settings
    expires_at = datetime.utcnow() + timedelta(hours=settings.JWT_EXPIRE_HOURS)
    
    return LoginResponse(token=token, expires_at=expires_at)


@router.post("/register", response_model=UserResponse, status_code=status.HTTP_201_CREATED)
async def register(request: UserCreate, db: AsyncSession = Depends(get_db)):
    result = await db.execute(select(User).where(User.username == request.username))
    if result.scalar_one_or_none():
        raise HTTPException(
            status_code=status.HTTP_409_CONFLICT,
            detail="username already exists"
        )
    
    result = await db.execute(select(User).where(User.email == request.email))
    if result.scalar_one_or_none():
        raise HTTPException(
            status_code=status.HTTP_409_CONFLICT,
            detail="email already exists"
        )
    
    user = User(
        username=request.username,
        email=request.email,
        password_hash=get_password_hash(request.password),
        role=request.role.value,
        status=UserStatus.ACTIVE.value,
    )
    
    db.add(user)
    await db.commit()
    await db.refresh(user)
    
    return UserResponse(
        id=user.id,
        username=user.username,
        email=user.email,
        role=user.role,
        status=user.status,
        created_at=user.created_at,
        last_login=user.last_login,
    )
