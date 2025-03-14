# Используем базовый образ Golang для сборки
FROM golang:1.24 AS builder

WORKDIR /app

# Копируем файлы
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Компилируем приложение с правильной архитектурой
RUN GOOS=linux GOARCH=amd64 go build -o main .

# Минимальный образ для финального контейнера
FROM debian:bullseye-slim

WORKDIR /app
COPY --from=builder /app/main .

# Даем права на исполнение
RUN chmod +x /app/main

# Запускаем приложение
CMD ["/app/main"]
