from sqlalchemy.ext.asyncio import create_async_engine, AsyncSession, async_sessionmaker
from sqlalchemy.orm import DeclarativeBase
from sqlalchemy import create_engine
from sqlalchemy.pool import NullPool
from typing import AsyncGenerator

from app.core.config import settings


class Base(DeclarativeBase):
    pass


if settings.DATABASE_TYPE == "mysql":
    engine = create_async_engine(
        settings.get_database_url(),
        echo=settings.DEBUG,
        pool_pre_ping=True,
        pool_recycle=3600,
    )
    sync_engine = create_engine(
        settings.get_database_url_sync(),
        echo=settings.DEBUG,
        pool_pre_ping=True,
        pool_recycle=3600,
        poolclass=NullPool,
    )
else:
    engine = create_async_engine(settings.DATABASE_URL, echo=settings.DEBUG)
    sync_engine = create_engine(settings.DATABASE_URL_SYNC, echo=settings.DEBUG, connect_args={"check_same_thread": False})

async_session_maker = async_sessionmaker(engine, class_=AsyncSession, expire_on_commit=False)


async def get_db() -> AsyncGenerator[AsyncSession, None]:
    async with async_session_maker() as session:
        try:
            yield session
            await session.commit()
        except Exception:
            await session.rollback()
            raise
        finally:
            await session.close()


async def init_db():
    async with engine.begin() as conn:
        await conn.run_sync(Base.metadata.create_all)
