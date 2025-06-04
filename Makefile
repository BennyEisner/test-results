run-api:
	cd api && go run main.go

run-cli:
	cd cli && go run main.go

run-frontend:
	cd frontend && npm start

lint:
	golangci-lint run ./...

test:
	go test ./...
