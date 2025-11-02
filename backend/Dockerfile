FROM golang:1.25.3 AS builder

WORKDIR /app
COPY go.work go.work.sum ./
COPY backend ./backend
COPY shared ./shared

RUN go work sync

WORKDIR /app/backend 
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/app ./cmd/app

FROM alpine:3.20

WORKDIR /app

COPY --from=builder /app/bin/app .

EXPOSE 8080

CMD ["./app"]