.PHONY: build run test clean deploy

build:
	@echo "Building application..."
	@cd backend && go build -o ../bin/server ./cmd/server

run:
	@echo "Starting development server..."
	@cd backend && air

test:
	@echo "Running tests..."
	@cd backend && go test ./... -v

clean:
	@echo "Cleaning up..."
	@rm -rf bin/
	@rm -rf backend/tmp/
	@docker-compose down -v

docker-build:
	@echo "Building Docker images..."
	@docker-compose build

docker-up:
	@echo "Starting Docker containers..."
	@docker-compose up -d

docker-down:
	@echo "Stopping Docker containers..."
	@docker-compose down

deploy:
	@echo "Deploying to production..."
	@git pull origin main
	@make docker-build
	@make docker-up
	@echo "Deployment complete!"

migrate:
	@echo "Running database migrations..."
	@cd backend && go run cmd/migrate/main.go

seed:
	@echo "Seeding database..."
	@cd backend && go run cmd/seed/main.go
