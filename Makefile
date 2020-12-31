generate-mocks:
	rm -r tests/mocks
	mockgen -destination tests/mocks/mock_ensurepkg/mock_ensurepkg.go -source ensurepkg/ensurepkg.go T

test:
	go test ./...

test-coverage:
	go test ./... -coverprofile=/tmp/ensure.coverage  && go tool cover -html=/tmp/ensure.coverage

lint:
	golangci-lint run
