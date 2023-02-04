generate-mocks:
	(cd cmd/ensure; go build -o ../../tmp/ensure) && ./tmp/ensure mocks generate

test:
	go test ./...

test-coverage:
	mkdir -p tmp
	go test ./... -coverprofile=tmp/ensure-pkg.coverage && go tool cover -html=tmp/ensure-pkg.coverage -o=tmp/coverage.html

lint:
	staticcheck ./...
	(cd cmd/ensure; staticcheck ./...)
	golangci-lint run
	(cd cmd/ensure; golangci-lint run)
