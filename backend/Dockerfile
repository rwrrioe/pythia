# syntax=docker/dockerfile:1.5
FROM golang:1.25.3 AS builder

ENV GOPRIVATE=github.com/rwrrioe/*
WORKDIR /app

RUN mkdir -p /root/.ssh && chmod 700 /root/.ssh && \
    ssh-keyscan github.com >> /root/.ssh/known_hosts

RUN git config --global url."ssh://git@github.com/".insteadOf https://github.com/
COPY backend/go.mod backend/go.sum ./
WORKDIR /app/backend 

RUN --mount=type=ssh go mod download

COPY backend/. ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/app ./cmd/app

FROM alpine:3.20
WORKDIR /app
COPY --from=builder /app/bin/app .
EXPOSE 8080
CMD ["./app"]