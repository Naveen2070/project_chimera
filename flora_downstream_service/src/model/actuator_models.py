from pydantic import BaseModel


class HealthCheckResponse(BaseModel):
    mongo: str
    postgres: str
    rabbitmq: str
    status: str
    code: int

    def to_dict(self) -> dict:
        """
        Converts the instance to a dictionary (optional in Pydantic, as `dict()` already works).
        """
        return self.model_dump()
