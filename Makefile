.PHONY: build up down logs test clean

# Собрать контейнеры
build:
	docker compose build

# Запустить контейнеры в фоне
up:
	docker compose up -d

# Остановить и удалить контейнеры
down:
	docker compose down

# Просмотр логов приложения
logs:
	docker compose logs -f app

# Запустить тесты (если есть)
test:
	go test ./...

# Очистить неиспользуемые ресурсы Docker (контейнеры, образы)
clean:
	docker system prune -f
