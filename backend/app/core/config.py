from pydantic_settings import BaseSettings
from functools import lru_cache


class Settings(BaseSettings):
    APP_NAME: str = "Password Timer Manager"
    VERSION: str = "1.0.0"
    DEBUG: bool = True

    DATABASE_URL: str = "postgresql+asyncpg://ptm_user:ptm_password@localhost:5432/ptm_db"
    DATABASE_URL_SYNC: str = "postgresql://ptm_user:ptm_password@localhost:5432/ptm_db"

    REDIS_HOST: str = "localhost"
    REDIS_PORT: int = 6379
    REDIS_DB: int = 0

    JWT_SECRET: str = "change-me-in-production"
    JWT_ALGORITHM: str = "HS256"
    JWT_EXPIRE_HOURS: int = 24

    STORAGE_TYPE: str = "local"
    STORAGE_LOCAL_PATH: str = "./storage"

    AES_KEY_LENGTH: int = 256
    RSA_KEY_LENGTH: int = 2048
    PASSWORD_MIN_LENGTH: int = 8

    class Config:
        env_file = ".env"
        case_sensitive = True


@lru_cache()
def get_settings() -> Settings:
    return Settings()


settings = get_settings()
