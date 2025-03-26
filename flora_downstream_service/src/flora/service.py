from typing import Any, Dict, Optional
import uuid
from sqlalchemy.future import select
from src.db.mongo.mongo_connect import mongo_engine
from src.db.postgres.postgres_connect import SessionLocal
from sqlalchemy.orm import Session
from fastapi import Depends
from sqlalchemy.ext.asyncio import AsyncSession
from src.model.common import FloraResponse
from src.model.flora import Flora
from src.model.mongo_flora import FloraMongo
from src.model.postgres_flora import FloraPG


async def get_floras(db: AsyncSession) -> Optional[FloraResponse]:
    try:
        result = await db.execute(select(FloraPG))
        florasPg: list[FloraPG] = result.scalars().all()

        if not florasPg:
            return FloraResponse(code=200, data=None)

        floras: list[Dict[str, Any]] = []
        for floraPg in florasPg:
            try:
                floraMongo = await mongo_engine.find_one(
                    FloraMongo, FloraMongo.flora_id == str(floraPg.id)
                )
            except Exception as e:
                print(f"An error occurred when retrieving Mongo data: {e}")
                return FloraResponse(
                    code=500,
                    data=f"An error occurred when retrieving Mongo data: {e}",
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
            floras.append(floraRes.model_dump_json())

        return FloraResponse(code=200, data=floras)

    except Exception as e:
        # Handle exceptions and log or raise them accordingly
        print(f"An error occurred: {e}")
        return FloraResponse(code=500, data=f"An error occurred: {e}")


async def get_flora(flora_id: str, db: AsyncSession) -> Optional[FloraResponse]:
    try:
        uid = uuid.UUID(flora_id)
    except ValueError:
        return FloraResponse(code=400, data="Invalid flora ID")

    result = await db.execute(select(FloraPG).where(FloraPG.id == uid))
    floraPg: FloraPG = result.scalar()
    if floraPg is None:
        return FloraResponse(code=404, data="Flora not found")

    floras: list[Dict[str, Any]] = []

    try:
        floraMongo = await mongo_engine.find_one(
            FloraMongo, FloraMongo.flora_id == str(flora_id)
        )
    except Exception as e:
        print(f"An error occurred when retrieving Mongo data: {e}")
        return FloraResponse(
            code=500,
            data=f"An error occurred when retrieving Mongo data: {e}",
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

    floras.append(floraRes.model_dump_json())

    return FloraResponse(code=200, data=floras)
