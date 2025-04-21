from pydantic import BaseModel


class ImagePayload(BaseModel):
    image_base64: str
