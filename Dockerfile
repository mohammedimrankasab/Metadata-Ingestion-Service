FROM golang:1.24-alpine

# Create a non-root user for safety
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /app

USER appuser