FROM golang:1.25.3 AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . . 
RUN go build -o /bin/app ./cmd/app

FROM alpine:3.20

WORKDIR /app

COPY --from=builder /bin/app .

EXPOSE 8080

CMD ["./app"]