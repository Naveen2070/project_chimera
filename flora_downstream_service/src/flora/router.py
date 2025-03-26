from fastapi import HTTPException, Depends
from sqlalchemy.ext.asyncio import AsyncSession
from src.db.postgres.postgres_connect import SessionLocal
from src.flora.service import get_floras, get_flora
from src.model.flora import Flora


# Dependency to get the database session
async def get_db():
    async with SessionLocal() as session:
        yield session


# Handler for processing requests based on the command
async def process_request(cmd: str, db: AsyncSession, data: dict = None) -> dict:
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
            return {"status": "success", "data": res or []}
        elif cmd == "get_flora_by_id":
            flora_id = data.get("param")
            if not flora_id:
                return {"status": "error", "message": "flora_id is required"}
            res = await get_flora(flora_id, db)
            return {"status": "success", "data": res or []}
        else:
            return {"status": "error", "message": f"Unknown command: {cmd}"}
    except Exception as e:
        print(e)
        return {"status": "error", "message": "Internal Server Error"}


# Example handlers
async def get_all_floras(db: AsyncSession = Depends(get_db)):
    """
    Fetch all flora records from the database.
    :param db: The database session.
    :return: A list of flora records.
    """
    try:
        res = await get_floras(db)
        return res or []
    except Exception as e:
        print(e)
        raise HTTPException(status_code=500, detail="Internal Server Error")


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
        raise HTTPException(status_code=500, detail="Internal Server Error")
