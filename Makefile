run:
	go run cmd/api/main.go

tidy:
	go mod tidy

build:
	go build -o bin/api cmd/api/main.go

test:
	go test -v ./...

test-cover:
	go test -cover ./...

# Docker commands
docker-build:
	docker compose build

docker-up:
	docker compose up -d

docker-down:
	docker compose down

docker-logs:
	docker compose logs -f

docker-restart:
	docker compose restart

docker-clean:
	docker compose down -v --rmi local