import uuid
from sqlalchemy import Column, String, Enum
from sqlalchemy.dialects.postgresql import UUID
from sqlalchemy.ext.declarative import declarative_base
from enum import Enum as en

Base = declarative_base()


class Type(en):
    PUBLIC = "public"
    PRIVATE = "private"


class FloraPG(Base):
    __tablename__ = "Flora"
    id = Column(UUID(as_uuid=True), primary_key=True, default=uuid.uuid4)
    user_id = Column(String, unique=True, index=True)
    common_name = Column(String, unique=True, index=True)
    scientific_name = Column(String, unique=True, index=True)
    type = Column(Enum(Type), nullable=False)
