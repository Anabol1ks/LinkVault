doc:
	docker-compose -f docker-compose.yml up -d

run:
	go run cmd/main.go

swag:
	swag init -g cmd/main.go