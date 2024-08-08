# Первый этап: сборка приложения
FROM golang:1.22 AS builder

WORKDIR /app

# Копируем go.mod для установки зависимостей
COPY go.mod go.sum ./

# копируем исходные файлы
COPY cmd/ cmd/
COPY internal/ internal/

# Устанавливаем зависимости
RUN go mod download

# Сборка приложения
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /todo-app cmd/task-tracker/main.go

# Второй этап: создание конечного образа
FROM alpine:3.20

# Установка доп. пакетов
RUN apk --no-cache add ca-certificates libc6-compat

WORKDIR /root/

# Копируем бинарный файл из предыдущего этапа
COPY --from=builder /todo-app .
COPY internal/database/migration migration/
COPY web/ web/

# Экспортируем порт, указанный в переменной окружения
ENV TODO_PORT 7540
EXPOSE ${TODO_PORT}

# Устанавливаем переменные окружения
ENV TODO_DBFILE /scheduler.db
ENV TODO_PASSWORD simple_password
ENV TODO_MIGRATION_PATH migration

# Запускаем приложение
CMD [ "./todo-app" ]