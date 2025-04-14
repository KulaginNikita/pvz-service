FROM golang:1.23-alpine AS builder

WORKDIR /app
COPY . .
ENV CGO_ENABLED=0
RUN go mod download
RUN go build -o pvz-service ./cmd/http

FROM alpine:3.18

WORKDIR /app
COPY --from=builder /app/pvz-service .

RUN chmod +x pvz-service
CMD ["./pvz-service"]

