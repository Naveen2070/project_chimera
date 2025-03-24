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

    def __repr__(self):
        return (
            f"Flora(id={self.id}, user_id={self.user_id}, "
            f"common_name={self.common_name}, scientific_name={self.scientific_name}, "
            f"type={self.type}, description={self.description}, "
            f"origin={self.origin}, other_details={self.other_details})"
        )
