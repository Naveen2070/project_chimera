import base64
from io import BytesIO
from PIL import Image
from fastapi import HTTPException


# Convert Base64 to PIL image
def decode_image(base64_str: str) -> Image.Image:
    try:
        img_data = base64.b64decode(base64_str)
        return Image.open(BytesIO(img_data)).convert("RGB")
    except Exception:
        raise HTTPException(status_code=400, detail="Invalid image format")
