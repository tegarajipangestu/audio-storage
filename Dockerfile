# Builder stage
FROM golang:1.23-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o /audio-storage cmd/api/main.go

# Final stage
FROM alpine:latest

RUN apk add --no-cache ffmpeg

WORKDIR /app

COPY --from=builder /audio-storage /audio-storage
RUN chmod +x /audio-storage

ENTRYPOINT ["/audio-storage"]
