from sqlalchemy import Select
from src.db.mongo.mongo_connect import mongo_engine
from src.db.postgres.postgres_connect import SessionLocal
from sqlalchemy.orm import Session
from fastapi import Depends
from sqlalchemy.ext.asyncio import AsyncSession
from src.model.flora import Flora
from src.model.mongo_flora import FloraMongo
from src.model.postgres_flora import FloraPG


async def get_db():
    async with SessionLocal() as session:
        yield session


async def get_floras(db: AsyncSession = Depends(get_db)) -> list[Flora] | None:
    res = await db.execute(Select(FloraPG))
    florasPg: list[FloraPG] = await db.query(FloraPG).all()
    if florasPg is None:
        return None
    floras: list[Flora] = []

    print(florasPg)
    for floraPg in florasPg:
        floraMongo = await mongo_engine.find_one(
            FloraMongo, FloraMongo.flora_id == floraPg.id
        )
        floraRes = Flora(
            id=floraPg.id,
            user_id=floraPg.user_id,
            common_name=floraPg.common_name,
            scientific_name=floraPg.scientific_name,
            type=floraPg.type,
            image=floraMongo.Image,
            description=floraMongo.Description,
            origin=floraMongo.Origin,
            other_details=floraMongo.OtherDetails,
        )
        floras.append(floraRes)

    print(floras + "+")

    return floras


async def get_flora(flora_id: str, db: Session = Depends(get_db)) -> Flora | None:
    floraPg: FloraPG = await db.query(FloraPG).filter(FloraPG.id == flora_id).first()
    if floraPg is None:
        return None

    floraMongo = await mongo_engine.find_one(
        FloraMongo, FloraMongo.flora_id == flora_id
    )
    floraRes = Flora(
        id=floraPg.id,
        user_id=floraPg.user_id,
        common_name=floraPg.common_name,
        scientific_name=floraPg.scientific_name,
        type=floraPg.type,
        image=floraMongo.Image,
        description=floraMongo.Description,
        origin=floraMongo.Origin,
        other_details=floraMongo.OtherDetails,
    )
    print(floraRes)
    return floraRes
