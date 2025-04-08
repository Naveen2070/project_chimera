# Use a more specific version of the official Python image
FROM python:3.13.2-slim

# Set environment variables for non-interactive installation
ENV PYTHONUNBUFFERED=1 \
    PIP_NO_CACHE_DIR=off \
    ENVIRONMENT=production \
    CONSUL_HOST=host.docker.internal \
    CONSUL_PORT=8500 \
    APP_PORT=8100 \
    DATABASE_URL_MONGO=mongodb://host.docker.internal:27017 \
    DATABASE_URL_POSTGRES=postgresql+asyncpg://postgres:naveen2007@host.docker.internal:5432/chimera_flora \
    RABBITMQ_URL=amqp://admin:naveen@2007@host.docker.internal:5672 \
    QUEUE_NAME=flora_downstream_queue

# Create and use a non-root user for security
RUN adduser --disabled-password --gecos '' appuser

# Set the working directory inside the container
WORKDIR /app

# Copy only the requirements file first to leverage Docker's cache
COPY requirements.txt /app/requirements.txt

# Install dependencies securely, and avoid caching
RUN pip install --no-cache-dir --trusted-host pypi.org --trusted-host pypi.python.org -r requirements.txt

# Copy the rest of the application files after dependencies are installed
COPY . /app/

# Change the user to the non-root user
USER appuser

# Expose the app's port (ensure this matches the `APP_PORT` you use)
EXPOSE 8100

# Command to run the FastAPI app using Uvicorn
CMD ["fastapi","run","src/main.py","--port","8100"]
