from fastapi import APIRouter

from app.api.handlers.health import router as health_router
from app.api.handlers.auth import router as auth_router
from app.api.handlers.verify import router as verify_router
from app.api.handlers.package import router as package_router
from app.api.handlers.password import router as password_router
from app.api.handlers.audit import router as audit_router
from app.api.handlers.admin import router as admin_router

api_router = APIRouter()

api_router.include_router(health_router)
api_router.include_router(auth_router)
api_router.include_router(verify_router)
api_router.include_router(package_router)
api_router.include_router(password_router)
api_router.include_router(audit_router)
api_router.include_router(admin_router)
