generate-mocks:
	(cd cmd/ensure; go build -o ../../tmp/ensure) && ./tmp/ensure mocks generate
	(cd cmd/ensure; make generate-mocks)
	(cd experiments/entable; make generate-mocks)

test:
	go test ./...
	(cd cmd/ensure; make test)
	(cd experiments/entable; make test)

test-coverage:
	mkdir -p tmp
	go test ./... -coverprofile=tmp/ensure.coverage && go tool cover -html=tmp/ensure.coverage -o=tmp/coverage.html
	(cd cmd/ensure; make test-coverage)
	(cd experiments/entable; make test-coverage)

lint:
	staticcheck ./...
	golangci-lint run
	(cd cmd/ensure; make lint)
	(cd experiments/entable; make lint)
