dev:
	go run cmd/server/main.go
	go run cmd/worker/main.go

dev-server:
	go run cmd/server/main.go

dev-worker:
	go run cmd/worker/main.go


.PHONY: dev dev-server dev-worker