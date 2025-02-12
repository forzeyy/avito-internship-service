#!/bin/sh

# Ждем, пока PostgreSQL станет доступным
./wait-for-it.sh db:5432 --timeout=30 --strict -- echo "PostgreSQL is up"

# Применяем миграции
if /usr/local/bin/sql-migrate up -config=/app/migrate.yaml; then
    echo "Migrations applied successfully"
else
    echo "Failed to apply migrations"
    exit 1
fi

# Запускаем основное приложение
echo "Starting application..."
exec "$@"