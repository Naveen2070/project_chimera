from fastapi import APIRouter, HTTPException
from src.db.postgres.postgres_connect import SessionLocal
from fastapi import Depends
from sqlalchemy.ext.asyncio import AsyncSession

from src.flora.service import get_floras, get_flora
from src.model.flora import Flora

flora_router = APIRouter(
    prefix="/flora",
    tags=["flora"],
)


async def get_db():
    async with SessionLocal() as session:
        yield session


@flora_router.get("/", response_model=list[Flora], status_code=200)
async def get_all_floras(db: AsyncSession = Depends(get_db)):
    try:
        res = await get_floras(db)
        if res is None:
            return []
        return res
    except Exception as e:
        print(e)
        raise HTTPException(status_code=500, detail="Internal Server Error")


@flora_router.get("/{flora_id}", response_model=Flora, status_code=200)
async def get_flora_by_id(flora_id: str, db: AsyncSession = Depends(get_db)):
    try:
        return await get_flora(flora_id, db)
    except Exception as e:
        print(e)
        raise HTTPException(status_code=500, detail="Internal Server Error")
