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

# Install ca-certificates and download static ffmpeg binary to avoid dependency issues
RUN apt-get update && \
    apt-get install -y --no-install-recommends ca-certificates wget xz-utils && \
    wget -q https://johnvansickle.com/ffmpeg/releases/ffmpeg-release-amd64-static.tar.xz && \
    tar xf ffmpeg-release-amd64-static.tar.xz && \
    mv ffmpeg-*-amd64-static/ffmpeg /usr/local/bin/ && \
    mv ffmpeg-*-amd64-static/ffprobe /usr/local/bin/ && \
    rm -rf ffmpeg-* && \
    apt-get remove -y wget xz-utils && \
    apt-get autoremove -y && \
    rm -rf /var/lib/apt/lists/*

COPY --from=builder /app/server .
COPY --from=builder /app/.env .
EXPOSE 8080
CMD ["./server"]