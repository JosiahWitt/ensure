generate-mocks:
	go run main.go mocks generate

test:
	go test ./...

test-coverage:
	mkdir -p tmp
	go test ./... -coverprofile=tmp/ensure-cmd.coverage && go tool cover -html=tmp/ensure-cmd.coverage -o=tmp/coverage.html

lint:
	golangci-lint run
