# Build stage with OpenJDK 21
FROM openjdk:21-jdk-slim AS build

WORKDIR /app

# Copy Maven wrapper and pom.xml
COPY mvnw ./
COPY .mvn .mvn
COPY pom.xml ./

# Download dependencies
RUN ./mvnw dependency:go-offline

# Copy source code
COPY src ./src

# Build the application and skip tests
RUN ./mvnw clean package -DskipTests

# Runtime stage with distroless image for security
FROM openjdk:21-jdk-slim AS runtime

WORKDIR /app

# Copy compiled JAR file from build stage
COPY --from=build /app/target/auth_service-0.0.1-SNAPSHOT.jar /app/auth-service.jar

# Create a non-root user for security
RUN useradd -ms /bin/bash appuser
USER appuser

# Expose application port
EXPOSE 8081

# Run the application
ENTRYPOINT ["java", "-jar", "auth-service.jar"]

# Health check
HEALTHCHECK --interval=30s --timeout=10s --retries=3 \
  CMD curl --fail http://localhost:8081/actuator/health || exit 1
