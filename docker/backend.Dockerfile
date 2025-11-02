# Build stage
FROM golang:1.25 AS builder
WORKDIR /app

# go.mod + go.sum eerst kopieren
COPY backend/go.mod backend/go.sum ./
RUN go mod download

# Kopieer source
COPY backend/ .

# .env file
COPY .env .

# Build the application
RUN go build -o server ./cmd/api

# Run stage - using debian-slim instead of distroless to install ffmpeg
FROM debian:bookworm-slim
WORKDIR /app

# Install ffmpeg (includes ffprobe) and clean up
RUN apt-get update && \
    apt-get install -y --no-install-recommends ffmpeg ca-certificates && \
    rm -rf /var/lib/apt/lists/*

COPY --from=builder /app/server .
COPY --from=builder /app/.env .
EXPOSE 8080
CMD ["./server"]