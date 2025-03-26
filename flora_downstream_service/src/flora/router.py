from fastapi import HTTPException, Depends
from sqlalchemy.ext.asyncio import AsyncSession
from src.db.postgres.postgres_connect import SessionLocal
from src.flora.service import get_floras, get_flora
from src.model.common import ApiResponse
from src.model.flora import Flora


# Dependency to get the database session
async def get_db():
    async with SessionLocal() as session:
        yield session


# Handler for processing requests based on the command
async def process_request(cmd: str, db: AsyncSession, data: dict = None) -> ApiResponse:
    """
    Process the incoming request based on the command and route it to appropriate handlers.
    :param cmd: The command indicating the action to perform.
    :param db: The database session for handling requests.
    :param data: Additional data for the request.
    :return: A response dictionary.
    """
    try:
        if cmd == "get_all_floras":
            res = await get_floras(db)

            if res.code == 200:
                return ApiResponse(status="success", code=res.code, data=res.data)
            else:
                return ApiResponse(status="error", code=res.code, data=res.data)

        elif cmd == "get_flora_by_id":
            flora_id = data.get("param")
            if not flora_id:
                return ApiResponse(
                    status="error", code=400, data="Flora ID not provided"
                )

            res = await get_flora(flora_id, db)

            if res.code == 200:
                return ApiResponse(status="success", code=res.code, data=res.data)
            else:
                return ApiResponse(status="error", code=res.code, data=res.data)

        else:
            return ApiResponse(status="error", code=400, data="Invalid command")
    except Exception as e:
        print(e)
        return ApiResponse(status="error", code=500, data="Internal Server Error")


# Example handlers
async def get_all_floras(db: AsyncSession = Depends(get_db)):
    """
    Fetch all flora records from the database.
    :param db: The database session.
    :return: A list of flora records.
    """
    try:
        res = await get_floras(db)
        return res
    except Exception as e:
        print(e)
        return ApiResponse(status="error", code=500, data=f"Internal Server Error {e}")


async def get_flora_by_id(flora_id: str, db: AsyncSession = Depends(get_db)):
    """
    Fetch a flora record by its ID from the database.
    :param flora_id: The ID of the flora record.
    :param db: The database session.
    :return: A flora record.
    """
    try:
        return await get_flora(flora_id, db)
    except Exception as e:
        print(e + "from get_flora_by_id", flora_id)
        raise ApiResponse(status="error", code=500, data=f"Internal Server Error {e}")
