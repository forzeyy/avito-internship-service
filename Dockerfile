# Stage 1: Builder
FROM golang:1.23-alpine AS builder
WORKDIR /app

# Копируем зависимости
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходники
COPY . .

# Строим бинарник приложения
RUN go build -o /app/bin/app ./cmd/avito/main.go

# Устанавливаем sql-migrate
RUN go install github.com/rubenv/sql-migrate/...@latest

# Stage 2: Final Image
FROM alpine:latest

WORKDIR /app

# Устанавливаем необходимые пакеты
RUN apk add --no-cache bash wget git

# Качаем и делаем исполнимым wait-for-it.sh
RUN wget https://raw.githubusercontent.com/vishnubob/wait-for-it/master/wait-for-it.sh && \
    chmod +x wait-for-it.sh

# Копируем файлы из builder
COPY --from=builder /app/bin/app /app/bin/app
COPY migrations/ /app/migrations/
COPY migrate.yaml /app/migrate.yaml
COPY entrypoint.sh /app/entrypoint.sh

# Делаем entrypoint.sh исполняемым
RUN chmod +x /app/entrypoint.sh

# Убедимся, что sql-migrate доступен в финальном контейнере
COPY --from=builder /go/bin/sql-migrate /usr/local/bin/sql-migrate

# Указываем точку входа и команду
ENTRYPOINT ["/app/entrypoint.sh"]
CMD ["/app/bin/app"]