import uuid
from bson import Binary
import bson
from sqlalchemy.future import select
from src.db.mongo.mongo_connect import mongo_engine
from src.db.postgres.postgres_connect import SessionLocal
from sqlalchemy.orm import Session
from fastapi import Depends
from sqlalchemy.ext.asyncio import AsyncSession
from src.model.flora import Flora
from src.model.mongo_flora import FloraMongo
from src.model.postgres_flora import FloraPG


async def get_floras(db: AsyncSession) -> list[Flora] | None:
    result = await db.execute(select(FloraPG))  # Correct async query
    florasPg: list[FloraPG] = result.scalars().all()  # Extract result
    if florasPg is None:
        return None
    floras: list[Flora] = []

    for floraPg in florasPg:
        floraMongo = await mongo_engine.find_one(
            FloraMongo, FloraMongo.flora_id == str(floraPg.id)
        )
        floraRes = Flora(
            id=str(floraPg.id),
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

    return floras


async def get_flora(flora_id: str, db: AsyncSession) -> Flora | None:
    uid = uuid.UUID(flora_id)
    result = await db.execute(select(FloraPG).where(FloraPG.id == uid))
    floraPg: FloraPG = result.scalar()
    if floraPg is None:
        return None

    floraMongo = await mongo_engine.find_one(
        FloraMongo, FloraMongo.flora_id == str(flora_id)
    )
    floraRes = Flora(
        id=str(floraPg.id),
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
