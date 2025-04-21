from fastapi import APIRouter
from src.flora.dto import ImagePayload
from src.flora.service import FloraService

flora_router = APIRouter(
    prefix="/actuator",
    tags=["Actuator"],
)

flora_service = FloraService()


@flora_router.post("/predict")
def predict(payload: ImagePayload):
    prediction = flora_service.predict(payload.image_base64)
    return {"prediction": prediction}
