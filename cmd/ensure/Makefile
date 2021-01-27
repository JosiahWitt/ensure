generate-mocks:
	@rm -r internal/mocks > /dev/null 2>&1 || true
	mockgen -destination internal/mocks/io/mock_fs/mock_fs.go bursavich.dev/fs-shim/io/fs ReadFileFS

test:
	go test ./...

test-coverage:
	go test ./... -coverprofile=/tmp/ensure-cmd.coverage && go tool cover -html=/tmp/ensure-cmd.coverage -o=../../tests/coverage-cmd.html

lint:
	golangci-lint run
