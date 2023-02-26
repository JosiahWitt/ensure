generate-mocks:
	(cd cmd/ensure; go build -o ../../tmp/ensure) && ./tmp/ensure mocks generate
	(cd cmd/ensure; make generate-mocks)
	(cd exp/entable; make generate-mocks)

test:
	go test ./...
	(cd cmd/ensure; make test)
	(cd exp/entable; make test)

test-coverage:
	mkdir -p tmp
	go test ./... -coverprofile=tmp/ensure.coverage && go tool cover -html=tmp/ensure.coverage -o=tmp/coverage.html
	(cd cmd/ensure; make test-coverage)
	(cd exp/entable; make test-coverage)

lint:
	staticcheck ./...
	golangci-lint run
	(cd cmd/ensure; make lint)
	(cd exp/entable; make lint)
