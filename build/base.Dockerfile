FROM golang:1.24

RUN apt-get update && \
    apt-get install -y tesseract-ocr libtesseract-dev libleptonica-dev

WORKDIR /app
COPY . .

RUN go mod download
RUN go build -o app .

CMD ["./app"]