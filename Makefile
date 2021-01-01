install-tools:
	go get github.com/golang/mock/mockgen@v1.4.4

generate-mocks:
	@rm -r tests/mocks &> /dev/null || true
	mockgen -destination tests/mocks/mock_ensurepkg/mock_ensurepkg.go -source ./ensurepkg/ensurepkg.go T

test:
	go test ./...

test-coverage:
	go test ./... -coverprofile=/tmp/ensure.coverage  && go tool cover -html=/tmp/ensure.coverage

lint:
	golangci-lint run

generate-toc:
	doctoc README.md
