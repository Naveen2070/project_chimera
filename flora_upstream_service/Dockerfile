# Build Stage
FROM node:22.14.0-slim AS build

# Set working directory
WORKDIR /app

# Install build dependencies (e.g., OpenSSL and libssl-dev for Prisma)
RUN apt-get update -y && apt-get install -y openssl libssl-dev

# Copy and install dependencies
COPY package*.json ./
RUN npm install

# Copy the rest of the app files
COPY . .

# Build app and generate prisma client
RUN npx prisma generate
RUN npm run build

# Production Stage
FROM node:22.14.0-slim

# Set working directory
WORKDIR /app

# Install runtime dependencies for Prisma (like OpenSSL)
RUN apt-get update -y && apt-get install -y openssl libssl-dev

# Copy and install production dependencies
COPY --from=build /app/package*.json ./
RUN npm install --only=production

# Copy build files
COPY --from=build /app/dist ./dist

# Copy the generated Prisma client
COPY --from=build /app/node_modules/@prisma/client ./node_modules/@prisma/client

# Copy the Prisma schema and other Prisma files (to the same path where they were in the build stage)
COPY --from=build /app/prisma ./prisma

# Add non-root user
RUN adduser --disabled-password --gecos "" appuser

# Set ownership of the app files to the non-root user
RUN chown -R appuser:appuser /app

# Switch to non-root user for security
USER appuser

# Expose port
EXPOSE 3030

# Set environment variables (you can add more if needed)
ENV DATABASE_URL="postgresql://postgres:naveen2007@host.docker.internal:5432/chimera_flora?schema=public"
ENV DATABASE_URL_MONGO="mongodb://host.docker.internal:27017/chimera_flora"
ENV PORT=3030
ENV RABBIT_MQ_URL="amqp://admin:naveen@2007@host.docker.internal:5672"
ENV RABBIT_QUEUE_FLORA="flora_upstream_queue"

# Start the application, generate Prisma client before starting the app in production mode
CMD ["sh", "-c", "npx prisma generate && npm run start:prod"]
