.PHONY: dev
dev:
	go run cmd/app/main.go

.PHONY: docker-setup
docker-setup:
	docker-compose --file='./docker/api/docker-compose.yaml' up --remove-orphans --build -d

.PHONY: swagger-gen
swagger-gen:
	swagger generate spec -o ./docs/swagger.yaml --scan-models
	
.PHONY: swagger-serve
swagger-serve:
	swagger serve --port=8080 --no-open -F=swagger ./docs/swagger.yaml