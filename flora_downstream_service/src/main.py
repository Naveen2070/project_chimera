from fastapi import FastAPI
from contextlib import asynccontextmanager
import consul
import uuid
from dotenv import load_dotenv
import os
from src.db.postgres.postgres_connect import database
from src.db.mongo.mongo_connect import mongo_client
from src.flora.router import flora_router

# Load environment variables
load_dotenv()

CONSUL_HOST = os.getenv("CONSUL_HOST", "localhost")
CONSUL_PORT = int(os.getenv("CONSUL_PORT", 8500))
APP_PORT = int(os.getenv("APP_PORT", 6000))


@asynccontextmanager
async def lifespan(app: FastAPI):
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
    )
    print(f"Service {service_id} registered with Consul!")
    await database.connect()
    print("Connected to Postgres database")

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


app = FastAPI(lifespan=lifespan)

app.include_router(flora_router)

if __name__ == "__main__":
    import uvicorn

    uvicorn.run(app, host="localhost", port=APP_PORT)
