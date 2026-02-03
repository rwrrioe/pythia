# syntax=docker/dockerfile:1.6
FROM golang:1.25.3 AS builder

WORKDIR /app

COPY backend/go.mod backend/go.sum ./

WORKDIR /app/backend

RUN go mod download

COPY backend/. ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/app ./cmd/app

FROM alpine:3.20
WORKDIR /app
COPY --from=builder /app/bin/app .
EXPOSE 8080
CMD ["./app"]
