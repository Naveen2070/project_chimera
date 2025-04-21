from contextlib import asynccontextmanager
import os
import uuid
import consul
from dotenv import load_dotenv
from fastapi import FastAPI

from src.flora.router import flora_router


load_dotenv()

CONSUL_HOST = os.getenv("CONSUL_HOST", "localhost")
CONSUL_PORT = int(os.getenv("CONSUL_PORT"))
APP_PORT = int(os.getenv("APP_PORT"))


@asynccontextmanager
async def lifespan(app: FastAPI):
    c = consul.Consul(host=CONSUL_HOST, port=CONSUL_PORT)
    id = str(uuid.uuid4())
    service_id = (
        "classification-service" + "-" + id.split("-")[0] + "-" + id.split("-")[2]
    )  # Generate a unique service ID

    # Startup logic
    c.agent.service.register(
        name="classification-service",
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

    try:
        # Yield control to the app
        yield
    finally:
        # Shutdown logic
        c.agent.service.deregister(service_id)
        print(f"Service {service_id} deregistered from Consul!")


app = FastAPI(lifespan=lifespan, openapi_url="/swagger/v1/openapi.json")


app.include_router(flora_router)

if __name__ == "__main__":
    import uvicorn

    uvicorn.run(app, host="localhost", port=APP_PORT)
