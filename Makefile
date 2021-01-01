install-tools:
	go get github.com/golang/mock/mockgen@v1.4.4
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.34.1

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
