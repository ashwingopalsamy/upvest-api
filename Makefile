# Service names
PUBLISHER_NAME=upvest-api-publisher
SUBSCRIBER_NAME=upvest-api-subscriber

# Docker Compose setup
DOCKER_COMPOSE=docker-compose

# Database connection
DB_DSN=postgres://upvest:upvest@localhost:5432/upvest?sslmode=disable
MIGRATIONS_DIR=schema/migrations

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

build-publisher:
	@echo "Building Publisher service..."
	go build -o $(PUBLISHER_NAME) ./cmd/upvest-api-publisher

build-subscriber:
	@echo "Building Subscriber service..."
	go build -o $(SUBSCRIBER_NAME) ./cmd/upvest-api-subscriber

build-all: build-publisher build-subscriber

run-publisher: build-publisher
	@echo "Running Publisher service..."
	./$(PUBLISHER_NAME)

run-subscriber: build-subscriber
	@echo "Running Subscriber service..."
	./$(SUBSCRIBER_NAME)

# Docker Compose commands
up:
	@echo "Starting all services..."
	$(DOCKER_COMPOSE) up --build -d
	@$(DOCKER_COMPOSE) logs -f upvest-api-publisher upvest-api-subscriber

down:
	@echo "Stopping all services..."
	$(DOCKER_COMPOSE) down

clean:
	@echo "Cleaning up..."
	rm -f $(PUBLISHER_NAME) $(SUBSCRIBER_NAME)
	$(DOCKER_COMPOSE) down -v

migrate-create:
	@echo "Creating new migration..."
	goose -dir $(MIGRATIONS_DIR) create $(NAME) sql

migrate-up:
	@echo "Running database migrations..."
	goose -dir $(MIGRATIONS_DIR) postgres "$(DB_DSN)" up

migrate-down:
	@echo "Rolling back the last migration..."
	goose -dir $(MIGRATIONS_DIR) postgres "$(DB_DSN)" down

migrate-status:
	@echo "Checking migration status..."
	goose -dir $(MIGRATIONS_DIR) postgres "$(DB_DSN)" status

test:
	@echo "Running tests..."
	go test -v ./...
