# Используем базовый образ Golang для сборки
FROM golang:1.24 AS builder

WORKDIR /app/cmd

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN GOOS=linux GOARCH=amd64 go build -o main ./cmd

# Используем образ с обновлённой glibc, например Ubuntu 22.04
FROM ubuntu:22.04
RUN apt-get update && apt-get install -y cron curl nano && rm -rf /var/lib/apt/lists/*

WORKDIR /app/cmd

COPY --from=builder /app/cmd/main .
COPY ../config/mycron /etc/cron.d/mycron
RUN chmod 0644 /etc/cron.d/mycron && crontab /etc/cron.d/mycron

COPY ../config/entrypoint.sh .
RUN chmod +x /app/cmd/main /app/cmd/entrypoint.sh

CMD ["/app/cmd/entrypoint.sh"]