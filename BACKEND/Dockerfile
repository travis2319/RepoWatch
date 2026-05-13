# ---------- Builder ----------
FROM golang:1.24-alpine AS builder

# sqlite needs CGO
RUN apk add --no-cache git gcc musl-dev build-base

WORKDIR /app

# Cache Go deps
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Enable CGO for sqlite
ENV CGO_ENABLED=1

# Build binary
RUN go build -o app ./cmd/main.go


# ---------- Runtime ----------
FROM alpine:latest

# runtime libs for sqlite
RUN apk add --no-cache libc6-compat

WORKDIR /app

# Copy binary
COPY --from=builder /app/app .

# optional if you use sqlite db file
COPY --from=builder /app/db ./db

EXPOSE 4000

CMD ["./app"]