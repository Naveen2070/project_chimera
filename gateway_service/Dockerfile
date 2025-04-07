# Build stage with OpenJDK 21
FROM openjdk:21-jdk-slim AS build

WORKDIR /app

# Copy Maven wrapper and pom.xml
COPY mvnw .
COPY .mvn .mvn
COPY pom.xml .

# Download project dependencies (cache if unchanged)
RUN ./mvnw dependency:go-offline

# Copy the source code
COPY src ./src

# Build the application and skip tests
RUN ./mvnw clean package -DskipTests

# Runtime stage with minimal OpenJDK 21 image
FROM openjdk:21-jdk-slim AS runtime

# Set the working directory
WORKDIR /app

# Copy the packaged jar from the build stage
COPY --from=build /app/target/gateway_service-0.0.1-SNAPSHOT.jar /app/gateway-service.jar

# Create a non-root user for security
RUN useradd -ms /bin/bash appuser
USER appuser

# Expose the application port
EXPOSE 8080

# Set the entrypoint for the application
ENTRYPOINT ["java", "-jar", "gateway-service.jar"]

# Health check (optional but recommended for monitoring)
HEALTHCHECK --interval=30s --timeout=10s --retries=3 \
  CMD curl --fail http://host.docker.internal:8080/actuator/health || exit 1
