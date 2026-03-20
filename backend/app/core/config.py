from pydantic_settings import BaseSettings
from functools import lru_cache
from typing import Optional


class Settings(BaseSettings):
    APP_NAME: str = "Password Timer Manager"
    VERSION: str = "1.0.0"
    DEBUG: bool = True

    DATABASE_TYPE: str = "sqlite"
    DATABASE_URL: str = "sqlite+aiosqlite:///./ptm.db"
    DATABASE_URL_SYNC: str = "sqlite:///./ptm.db"
    
    MYSQL_HOST: str = "localhost"
    MYSQL_PORT: int = 3306
    MYSQL_USER: str = "root"
    MYSQL_PASSWORD: str = "password"
    MYSQL_DATABASE: str = "ptm"

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

    def get_database_url(self) -> str:
        if self.DATABASE_TYPE == "mysql":
            return f"mysql+aiomysql://{self.MYSQL_USER}:{self.MYSQL_PASSWORD}@{self.MYSQL_HOST}:{self.MYSQL_PORT}/{self.MYSQL_DATABASE}"
        return "sqlite+aiosqlite:///./ptm.db"
    
    def get_database_url_sync(self) -> str:
        if self.DATABASE_TYPE == "mysql":
            return f"mysql+pymysql://{self.MYSQL_USER}:{self.MYSQL_PASSWORD}@{self.MYSQL_HOST}:{self.MYSQL_PORT}/{self.MYSQL_DATABASE}"
        return "sqlite:///./ptm.db"


@lru_cache()
def get_settings() -> Settings:
    return Settings()


settings = get_settings()
