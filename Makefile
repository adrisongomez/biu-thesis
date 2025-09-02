docker-dev:
	docker compose -f ./docker/docker-compose.development.yml up

start-worker:
	go run ./cmd/worker/main.go

start-client:
	go run ./cmd/client/main.go

start-server:
	go run ./cmd/server/main.go
