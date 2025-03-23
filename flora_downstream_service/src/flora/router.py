from fastapi import APIRouter, HTTPException


from src.flora.service import get_floras
from src.model.flora import Flora

flora_router = APIRouter(
    prefix="/flora",
    tags=["flora"],
)


@flora_router.get("/", response_model=list[Flora], status_code=200)
async def get_all_floras():
    try:
        res = await get_floras()
        if res is None:
            return []
        return res
    except Exception as e:
        print(e)
        raise HTTPException(status_code=500, detail="Internal Server Error")
