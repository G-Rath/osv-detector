BINARY=osv-detector
VERSION=0.1
OS_ARCH=linux_amd64

build:
	go build -o ${BINARY}

test:
	go test ./... -parallel=4

lint:	lint-with-golangci-lint lint-with-go-fmt

lint-with-golangci-lint:
	go run github.com/golangci/golangci-lint/cmd/golangci-lint@v1.44.2 run ./... --max-same-issues 0

lint-with-go-fmt:
	gofmt -s -d */**.go
