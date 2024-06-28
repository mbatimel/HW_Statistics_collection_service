.PHONY: service
service:
	go build -o service cmd/service/main.go  
	./service

.PHONY: migrate
migrate:
	go build -o migrate cmd/migrate/main.go
	./migrate

.PHONY: up
up:
	docker-compose up -d

.PHONY: down
down:
	docker-compose down

.DEFAULT_GOAL := build