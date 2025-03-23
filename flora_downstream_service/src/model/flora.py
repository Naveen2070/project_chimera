from typing import Any, Dict
from pydantic import BaseModel
from enum import Enum

from src.model.postgres_flora import Type


# Unified Flora Class
class Flora(BaseModel):
    id: str
    user_id: str
    common_name: str
    scientific_name: str
    type: Type
    image: bytes
    description: str
    origin: str
    other_details: Dict[str, Any]
