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

import asyncio
from fastapi import FastAPI
from contextlib import asynccontextmanager
import consul
import uuid
from dotenv import load_dotenv
import os
from src.db.postgres.postgres_connect import database
from src.db.mongo.mongo_connect import mongo_client
from src.queue.rabbit_consumer import RpcConsumer
from src.actuator.router import actuator_router

# Load environment variables
load_dotenv()

CONSUL_HOST = os.getenv("CONSUL_HOST", "localhost")
CONSUL_PORT = int(os.getenv("CONSUL_PORT"))
APP_PORT = int(os.getenv("APP_PORT"))
AMQP_URL = os.getenv("RABBITMQ_URL")
QUEUE_NAME = os.getenv("QUEUE_NAME")


@asynccontextmanager
async def lifespan(app: FastAPI):
    consumer = RpcConsumer(AMQP_URL, QUEUE_NAME)
    await consumer.connect()
    c = consul.Consul(host=CONSUL_HOST, port=CONSUL_PORT)
    id = str(uuid.uuid4())
    service_id = (
        "flora-downstream-service" + "-" + id.split("-")[0] + "-" + id.split("-")[2]
    )  # Generate a unique service ID

    # Startup logic
    c.agent.service.register(
        name="flora-downstream-service",
        service_id=service_id,
        address="localhost",
        port=APP_PORT,
        check=consul.Check.http(
            url="http://host.docker.internal:" + str(APP_PORT) + "/actuator/health",
            interval="15s",
            timeout="10s",
            deregister="10s",
        ),
    )
    print(f"Service {service_id} in port {APP_PORT} registered with Consul!")
    await database.connect()
    print("Connected to Postgres database")
    asyncio.create_task(consumer.start())
    print(f"Started {QUEUE_NAME} RPC Consumer")

    try:
        # Yield control to the app
        yield
    finally:
        # Shutdown logic
        c.agent.service.deregister(service_id)
        print(f"Service {service_id} deregistered from Consul!")
        await database.disconnect()
        print("Disconnected from Postgres database")
        if mongo_client is not None:
            mongo_client.close()
            print("Disconnected from MongoDB")


app = FastAPI(lifespan=lifespan, openapi_url="/swagger/v1/openapi.json")

app.include_router(actuator_router)

if __name__ == "__main__":
    import uvicorn

    uvicorn.run(app, host="localhost", port=APP_PORT)
