from odmantic import AIOEngine
from motor.motor_asyncio import AsyncIOMotorClient
import os

# Environment variable for MongoDB connection URL
MONGO_URL = os.getenv("DATABASE_URL_MONGO", "mongodb://localhost:27017")

# Initialize MongoDB engine
mongo_client = AsyncIOMotorClient(MONGO_URL)
mongo_engine = AIOEngine(client=mongo_client, database="chimera_flora")
