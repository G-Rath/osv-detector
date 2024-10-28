BINARY=osv-detector
VERSION=0.1
OS_ARCH=linux_amd64

.PHONY: ${BINARY}

${BINARY}:
	go build -o ${BINARY}

build: ${BINARY}

build-snapshot:
	goreleaser build --single-target --snapshot --rm-dist -o osv-detector

test:
	go test ./...

test-with-coverage:
	go test ./... -cover

lint:	lint-with-golangci-lint lint-with-go-fmt

lint-with-golangci-lint:
	go run github.com/golangci/golangci-lint/cmd/golangci-lint@v1.59.1 run ./... --max-same-issues 0

lint-with-go-fmt:
	gofmt -s -d */**.go

format-with-prettier:
	npx prettier --prose-wrap always --write .

regenerate-e2e-fixtures: ${BINARY}
	OSV_DETECTOR_CMD=./${BINARY} ./generators/generate-e2e-fixtures.js
