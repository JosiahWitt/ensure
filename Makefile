install-tools:
	go get github.com/golang/mock/mockgen

generate-mocks:
	@rm -r internal/mocks > /dev/null 2>&1 || true
	mockgen -destination internal/mocks/mock_ensurepkg/mock_ensurepkg.go -source ./ensurepkg/ensurepkg.go T
	@(cd cmd/ensure; make generate-mocks)

test:
	go test ./...
	@(cd cmd/ensure; make test)

test-coverage:
	go test ./... -coverprofile=/tmp/ensure-pkg.coverage && go tool cover -html=/tmp/ensure-pkg.coverage -o=tests/coverage-pkg.html
	@(cd cmd/ensure; make test-coverage)

lint:
	golangci-lint run
	@(cd cmd/ensure; make lint)

generate-toc:
	doctoc README.md
