generate-mocks:
	(cd ../../cmd/ensure; go build -o ../../tmp/ensure) && ../../tmp/ensure mocks generate

test:
	go test ./...

test-coverage:
	mkdir -p tmp
	go test ./... -coverprofile=tmp/ensure.coverage && go tool cover -html=tmp/ensure.coverage -o=tmp/coverage.html

lint:
	golangci-lint run
	staticcheck ./...
