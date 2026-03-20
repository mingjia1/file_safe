from fastapi import APIRouter

from app.api.handlers.health import router as health_router
from app.api.handlers.auth import router as auth_router
from app.api.handlers.verify import router as verify_router

api_router = APIRouter()

api_router.include_router(health_router)
api_router.include_router(auth_router)
api_router.include_router(verify_router)
