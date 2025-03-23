from sqlalchemy import MetaData, Column, Integer, String, create_engine
from sqlalchemy.ext.asyncio import create_async_engine
from sqlalchemy.ext.asyncio import AsyncSession
from sqlalchemy.ext.declarative import declarative_base
from sqlalchemy.orm import sessionmaker
from databases import Database
import os

# Load PostgreSQL connection URL from environment variables
POSTGRES_URL = os.getenv("DATABASE_URL_POSTGRES")

# Create the Database and Engine
database = Database(POSTGRES_URL)
async_engine = create_async_engine(POSTGRES_URL, echo=True)

# Create the Base and Session
metadata = MetaData(schema="public")
Base = declarative_base()
SessionLocal = sessionmaker(
    bind=async_engine,
    class_=AsyncSession,  # Use AsyncSession here
    expire_on_commit=False,
)
