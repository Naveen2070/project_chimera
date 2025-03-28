import os
from aio_pika import connect
from src.db.postgres.postgres_connect import database
from src.db.mongo.mongo_connect import mongo_engine
from src.model.actuator_models import HealthCheckResponse
from src.model.mongo_flora import FloraMongo


async def Is_Healthy() -> HealthCheckResponse:
    try:
        # Check MongoDB connection
        await mongo_engine.find_one(FloraMongo)  # Try a simple query
        mongo_status = "healthy"
    except Exception as e:
        mongo_status = f"error: {str(e)}"

    try:
        # Check PostgreSQL connection
        if not database.is_connected:
            await database.connect()
        postgres_status = "healthy"
    except Exception as e:
        postgres_status = f"error: {str(e)}"

    try:
        # Check RabbitMQ connection
        connection = await connect(os.getenv("RABBITMQ_URL"))
        await connection.close()
        rabbitmq_status = "healthy"
    except Exception as e:
        rabbitmq_status = f"error: {str(e)}"

    statusCode = 200
    if (
        mongo_status != "healthy"
        or postgres_status != "healthy"
        or rabbitmq_status != "healthy"
    ):
        statusCode = 500

    return HealthCheckResponse(
        mongo=mongo_status,
        postgres=postgres_status,
        rabbitmq=rabbitmq_status,
        status="healthy",
        code=statusCode,
    )
