.PHONY: dev
dev:
	go run cmd/app/main.go

.PHONY: docker-setup
docker-setup:
	docker-compose --file='./docker/api/docker-compose.yaml' up --remove-orphans --build -d