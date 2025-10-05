FROM golang:1.24-alpine

RUN apk add --no-cache git build-base

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY cmd/ ./cmd/
COPY db/ ./db/
COPY internal/ ./internal/

# In BACKEND/Dockerfile, after copying backend source
COPY build/ ./build/


# The frontend build will be mounted at runtime.
CMD ["go", "run", "cmd/main.go"]
