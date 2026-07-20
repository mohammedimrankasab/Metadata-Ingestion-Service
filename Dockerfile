# -----------------------------------------------------------------------------
# Stage 1 - Build
# -----------------------------------------------------------------------------
FROM golang:1.26-alpine AS builder

# Install certificates and git (required for downloading some Go modules)
RUN apk add --no-cache git ca-certificates

WORKDIR /src

# Copy dependency files first for better layer caching
COPY go.mod .
COPY go.sum .

RUN go mod download

# Copy application source
COPY . .

# Build a statically linked binary
ENV CGO_ENABLED=0 \
    GOOS=linux

RUN go build \
    -ldflags="-s -w" \
    -o metadata-ingestion-service \
    ./cmd/app

# -----------------------------------------------------------------------------
# Stage 2 - Runtime
# -----------------------------------------------------------------------------
FROM alpine:3.22

# Install CA certificates for HTTPS communication
RUN apk --no-cache add ca-certificates

# Create a non-root user
RUN addgroup -S appgroup && \
    adduser -S appuser -G appgroup

WORKDIR /app

# Copy the compiled binary
COPY --from=builder /src/metadata-ingestion-service .

# Switch to non-root user
USER appuser

# Expose application port
EXPOSE 8080

# Run the application
ENTRYPOINT ["./metadata-ingestion-service"]