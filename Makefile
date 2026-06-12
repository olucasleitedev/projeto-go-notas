.PHONY: up down logs build-ms run-gateway run-notes run-audit

up:
	docker compose up --build -d

down:
	docker compose down

logs:
	docker compose logs -f

build-ms:
	go build -o bin/gateway.exe ./cmd/gateway
	go build -o bin/notes-service.exe ./cmd/notes-service
	go build -o bin/audit-service.exe ./cmd/audit-service

run-gateway:
	go run ./cmd/gateway

run-notes:
	go run ./cmd/notes-service

run-audit:
	go run ./cmd/audit-service
