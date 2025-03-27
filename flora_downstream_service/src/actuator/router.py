from fastapi import APIRouter, HTTPException
from fastapi.responses import JSONResponse

from src.actuator.service import Is_Healthy
from src.model.actuator_models import HealthCheckResponse


actuator_router = APIRouter(
    prefix="/actuator",
    tags=["Actuator"],
)


@actuator_router.get(
    "/health",
    response_model=HealthCheckResponse,
    summary="Health Check Endpoint for Flora Downstream Service",
)
async def health_check():
    try:
        status = await Is_Healthy()
        if status.code == 200:
            return JSONResponse(content=status.model_dump(), status_code=status.code)
        else:
            raise HTTPException(status_code=status.code, detail=status)
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))
