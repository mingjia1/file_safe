from fastapi import Request, HTTPException, status
from fastapi.security import HTTPBearer, HTTPAuthorizationCredentials
from typing import Optional

from app.core.security import decode_access_token
from app.models.user import User

security = HTTPBearer(auto_error=False)


async def get_current_user(request: Request) -> Optional[dict]:
    auth_header = request.headers.get("Authorization")
    
    if not auth_header or not auth_header.startswith("Bearer "):
        return None
    
    token = auth_header.replace("Bearer ", "")
    payload = decode_access_token(token)
    
    if not payload:
        return None
    
    return {
        "user_id": payload.get("sub"),
        "username": payload.get("username"),
        "role": payload.get("role"),
    }


async def require_auth(request: Request) -> dict:
    user = await get_current_user(request)
    if not user:
        raise HTTPException(
            status_code=status.HTTP_401_UNAUTHORIZED,
            detail="Not authenticated"
        )
    return user


def require_permission(permission: str):
    async def permission_checker(request: Request) -> dict:
        user = await require_auth(request)
        
        from app.models.user import has_permission, UserRole
        role = UserRole(user.get("role", "viewer"))
        
        if not has_permission(role, permission):
            raise HTTPException(
                status_code=status.HTTP_403_FORBIDDEN,
                detail="Insufficient permissions"
            )
        
        return user
    
    return permission_checker
