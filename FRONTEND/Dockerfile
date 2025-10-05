# Stage 1: Build
FROM node:lts-trixie-slim AS builder

WORKDIR /app

COPY package*.json ./
RUN npm i

COPY . .
RUN npm run build

# Stage 2: Run
FROM node:lts-trixie-slim

WORKDIR /app

COPY --from=builder /app/package*.json ./
COPY --from=builder /app/build ./build

RUN npm i --omit=dev

EXPOSE 3000
CMD ["node", "build"]
