import json
from pydantic import BaseModel, Field
from typing import Dict, Any


class ErrorDTO(BaseModel):
    type: str = Field(..., description="Type of the error")
    status: str = Field(..., description="Status of the response, e.g., 'error'")
    code: int = Field(..., description="HTTP or application-specific status code")
    data: Dict[str, Any] = Field(
        default_factory=dict, description="Additional error details"
    )

    @classmethod
    def from_raw(cls, data: dict) -> "ErrorDTO":
        try:
            return cls(**data)
        except Exception as e:
            print(f"Invalid error data: {e}")
            raise ValueError(f"Invalid error data: {e}")

    def to_dict(self) -> Dict[str, Any]:
        """
        Returns the object as a standard dictionary.
        """
        return self.model_dump()
