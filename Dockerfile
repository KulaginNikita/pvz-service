# === Этап сборки ===
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY . .

# ВАЖНО: билдим статически, чтобы работало в Alpine
ENV CGO_ENABLED=0
RUN go mod download
RUN go build -o pvz-service ./cmd/http

# === Финальный образ ===
FROM alpine:3.18

WORKDIR /app
COPY --from=builder /app/pvz-service .

# Делаем бинарник исполняемым
RUN chmod +x pvz-service

CMD ["./pvz-service"]

