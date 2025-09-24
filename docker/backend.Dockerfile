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
# Debug: Show .env content during build
RUN echo "===== .env file content during build =====" && cat .env && echo "===== end of .env file ====="

# Build the application
RUN go build -o server ./cmd/...

# Run stage
FROM gcr.io/distroless/base-debian12
WORKDIR /app
COPY --from=builder /app/server .
COPY --from=builder /app/.env .
EXPOSE 8080
CMD ["./server"]