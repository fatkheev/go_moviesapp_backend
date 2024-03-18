FROM golang:1.21-alpine

# Установка make, git и curl для скачивания и установки migrate
RUN apk update && apk add --no-cache make git curl

# Установка migrate
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.14.1/migrate.linux-amd64.tar.gz | tar xvz && \
    mv migrate.linux-amd64 /usr/local/bin/migrate

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o go_moviesapp_backend ./cmd/main.go

EXPOSE 8081

# Перед стартом приложения выполняем миграции через make migrate-up
CMD sh -c "make migrate-up && ./go_moviesapp_backend"