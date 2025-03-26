import uuid
from sqlalchemy import Column, String, Enum
from sqlalchemy.dialects.postgresql import UUID
from sqlalchemy.ext.declarative import declarative_base
from enum import Enum as en

Base = declarative_base()


class Type(str, en):
    public = "public"
    private = "private"


class FloraPG(Base):
    __tablename__ = "Flora"
    id = Column(UUID(as_uuid=True), primary_key=True, default=uuid.uuid4)
    user_id = Column(String, unique=True, index=True)
    common_name = Column(String, unique=True, index=True)
    scientific_name = Column(String, unique=True, index=True)
    type = Column(Enum(Type), nullable=False)

    def __repr__(self):
        return (
            f"FloraPG(id={self.id}, user_id={self.user_id}, "
            f"common_name={self.common_name}, scientific_name={self.scientific_name}, "
            f"type={self.type})"
        )
