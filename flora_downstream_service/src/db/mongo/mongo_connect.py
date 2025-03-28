# 	Copyright 2025 Naveen R
#
# 		Licensed under the Apache License, Version 2.0 (the "License");
# 		you may not use this file except in compliance with the License.
# 		You may obtain a copy of the License at
#
# 		http://www.apache.org/licenses/LICENSE-2.0
#
# 		Unless required by applicable law or agreed to in writing, software
# 		distributed under the License is distributed on an "AS IS" BASIS,
# 		WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# 		See the License for the specific language governing permissions and
# 		limitations under the License.

from odmantic import AIOEngine
from motor.motor_asyncio import AsyncIOMotorClient
import os

# Environment variable for MongoDB connection URL
MONGO_URL = os.getenv("DATABASE_URL_MONGO", "mongodb://localhost:27017")

# Initialize MongoDB engine
mongo_client = AsyncIOMotorClient(MONGO_URL)
mongo_engine = AIOEngine(client=mongo_client, database="chimera_flora")
