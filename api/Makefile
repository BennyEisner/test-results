# Default is unit tests only
test:
	go test ./...

build:
	go build cmd/server/main.go

run:
	go run cmd/server/main.go

test-integration:
	go test -tags=integration ./...

test-e2e:
	go test ./tests/e2e
