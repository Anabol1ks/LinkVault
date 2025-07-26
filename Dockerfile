# syntax=docker/dockerfile:1
FROM golang:1.24
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o linkvault ./cmd/main.go

CMD ["./linkvault"]
