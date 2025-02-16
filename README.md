# Avito Merchandise Shop Service

Сервис позволяет пользователям авторизоваться, отправлять монеты другим пользователям, совершать покупки мерча и получать информацию о балансе, истории транзакций и инвентаре.


## **Содержание**
1. [Требования](#требования)
2. [Установка](#установка)
3. [Запуск](#запуск)
4. [О проекте](#о-проекте)


## **Требования**

Для запуска проекта необходимы следующие компоненты:
- Docker
- Docker Compose
- Go (версия 1.23)
- PostgreSQL (используется в контейнере)


## **Установка**

1. **Клонируйте репозиторий:**
   ```bash
   git clone https://github.com/forzeyy/avito-internship-service.git
   cd avito-internship-service
   
2. **Установите зависимости**
    Если вы хотите собрать проект локально:
    ```bash
    go mod download


## **Запуск**
1. **Запуск через Docker Compose:**
 - Создайте файл .env с переменными окружения в ./cmd/avito (пример ниже).
 - Запустите контейнеры:
    ```bash
    docker-compose up --build
 - Сервис будет доступен по адресу http://localhost:8080.
2. **Пример файла .env (должен лежать в ./cmd/avito):**
   ```env
   DB_USER=postgres
   DB_PASSWORD=postgres
   DB_HOST=db
   DB_PORT=5432
   DB_NAME=avito
   JWT_SECRET=secretsecret
3. **Остановка сервиса:**
    ```bash
    docker-compose down


# О проекте
Проект разработан в рамках отбора на стажировку в Avito.
