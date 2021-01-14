install-tools:
	go get github.com/golang/mock/mockgen@v1.4.4

generate-mocks:
	@rm -r internal/mocks &> /dev/null || true
	mockgen -destination internal/mocks/mock_ensurepkg/mock_ensurepkg.go -source ./ensurepkg/ensurepkg.go T

test:
	go test ./...

test-coverage:
	go test ./... -coverprofile=/tmp/ensure.coverage && go tool cover -html=/tmp/ensure.coverage -o=tests/coverage.html

lint:
	golangci-lint run

generate-toc:
	doctoc README.md
