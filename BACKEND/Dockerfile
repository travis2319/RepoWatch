FROM golang:1.24-alpine

RUN apk add --no-cache git build-base

WORKDIR /app

# Copy go mod files first
COPY go.mod go.sum ./
RUN go mod download

# Copy all source
COPY . .

# Run the app directly (no binary build)
CMD ["go", "run", "cmd/main.go"]
