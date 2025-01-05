SERVICE_NAME=upvest-api
DOCKER_COMPOSE=docker-compose -f docker-compose.yml -f docker-compose-db.yml

dep:
	@echo "Handling dependencies..."
	go mod tidy
	go mod vendor

lint:
	@echo "Running linter..."
	golangci-lint run

fmt:
	@echo "Formatting Go code..."
	go fmt ./...

build:
	@echo "Building Go service..."
	go build -o $(SERVICE_NAME) ./cmd/$(SERVICE_NAME)

run: build
	@echo "Building and running Go service..."
	./$(SERVICE_NAME)

test:
	@echo "Running tests..."
	go test -v ./...

up:
	@echo "Starting Go service and PostgreSQL database..."
	$(DOCKER_COMPOSE) up --build -d
	@$(DOCKER_COMPOSE) logs -f upvest-api

down:
	@echo "Stopping Go service and PostgreSQL database..."
	$(DOCKER_COMPOSE) down

clean:
	@echo "Cleaning up..."
	rm -f $(SERVICE_NAME)
	$(DOCKER_COMPOSE) down -v

migrate-create:
	@echo "Creating new migration..."
	goose -dir schema/migrations create $(NAME) sql

migrate-up:
	@echo "Running up migrations..."
	goose -dir schema/migrations postgres "$(DB_URL)" up

migrate-down:
	@echo "Rolling back the last migration..."
	goose -dir schema/migrations postgres "$(DB_URL)" down

migrate-status:
	@echo "Checking migration status..."
	goose -dir schema/migrations postgres "$(DB_URL)" status