FROM golang:1.24-alpine

RUN apk add --no-cache git build-base

WORKDIR /app

# Copy go.mod & go.sum first for caching dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy only the needed source files (exclude .env)
COPY cmd/ ./cmd/
COPY db/ ./db/
COPY internal/ ./internal/
# COPY app ./app   # optional, only if app binary is needed

# Run the app directly
CMD ["go", "run", "cmd/main.go"]
