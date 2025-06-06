# Используем официальный образ Go
FROM golang:1.23.8-alpine AS builder

# Устанавливаем зависимости
RUN apk add --no-cache git

# Копируем исходный код
WORKDIR /app
COPY . .
RUN go mod download



# Собираем приложение
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/avito-tech ./cmd/avito_tech/main.go

# Финальный образ
FROM alpine:latest

#
#WORKDIR /app
# Копируем бинарник из builder
COPY --from=builder /app/avito-tech /app/avito-tech
COPY --from=builder /app/.env .
COPY --from=builder /app/migrations ./migrations/


#RUN chmod 644 /app/.env
# Открываем порт
EXPOSE 8080

# Запускаем сервер
CMD ["/app/avito-tech"]



#/avito_tech
#├── cmd
#│    └── avito_tech
#│          └── main.go
#├── internal
#│   ├── app
#│   ├── auth
#│   ├── middleware
#│   ├── storage
#│   └── service
#├── migrations
#├── .env
#├── docker-compose.yaml
#├── Dockerfile
