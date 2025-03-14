# Используем базовый образ Golang для сборки
FROM golang:1.24 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN GOOS=linux GOARCH=amd64 go build -o main .

# Используем образ с обновлённой glibc, например Ubuntu 22.04
FROM ubuntu:22.04

WORKDIR /app
COPY --from=builder /app/main .

RUN chmod +x /app/main

CMD ["/app/main"]