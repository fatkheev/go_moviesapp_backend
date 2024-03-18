.PHONY: build run db-up clean

BINARY_NAME=filmoteca

# Сборка приложения
build:
	@echo "Сборка приложения..."
	go build -o $(BINARY_NAME) ./cmd/main.go

# Запуск приложения
run:
	@echo "Запуск приложения..."
	@make build
	./$(BINARY_NAME)

# Запуск всех тестов
test:
	@echo "Запуск тестов..."
	go test ./...

# Запуск сервисов Docker Compose (например, базы данных)
docker-start:
	@echo "Запуск контейнера с базой данных..."
	docker-compose up -d

# Останавливаем контейнер бэкенда
stop-backend:
	@echo "Остановка всех контейнеров бэкенда, начинающихся с go_moviesapp_backend-app..."
	@docker ps --filter "name=go_moviesapp_backend-app" -q | xargs -r docker stop

# Останавливаем все контейнеры
stop-all:
	@echo "Остановка всех контейнеров..."
	docker stop $$(docker ps -aq) || true

# Удаляем контейнеры, связанные с проектом
remove-containers: stop-all
	@echo "Удаление контейнеров проекта..."
	@docker ps -a --filter "name=go_moviesapp_backend-app" --format "{{.ID}}" | xargs -r docker rm || true
	@docker ps -a --filter "name=swagger_ui" --format "{{.ID}}" | xargs -r docker rm || true
	@docker ps -a --filter "ancestor=postgres:13" --format "{{.ID}}" | xargs -r docker rm || true

# Удаляем образы, связанные с проектом
remove-images: stop-all remove-containers
	@echo "Удаление образов проекта..."
	@docker images --filter "reference=swaggerapi/swagger-ui" --format "{{.ID}}" | xargs -r docker rmi || true
	@docker images --filter "reference=postgres:13" --format "{{.ID}}" | xargs -r docker rmi || true
	
# Применение миграций
migrate-up:
	@echo "Применение миграций..."
	@eval $$(cat .env | sed 's/^/export /') && \
	migrate -path docs/migrations -database "postgresql://filmoteca:filmoteca@db:5432/filmoteca?sslmode=disable" up

# Откат миграций
migrate-down:
	@echo "Откат миграций..."
	@eval $$(cat .env | sed 's/^/export /') && \
	migrate -path docs/migrations -database "postgresql://$$DB_USER:$$DB_PASSWORD@localhost:$$DB_PORT/$$DB_NAME?sslmode=disable" down

# Очистка сгенерированных файлов
clean:
	@echo "Очистка..."
	@rm -f $(BINARY_NAME)