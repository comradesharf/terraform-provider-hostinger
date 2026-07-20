default: fmt lint install generate

build:
	go build -v ./...

install: build
	go install -v ./...

lint:
	golangci-lint run

generate:
	cd tools; go generate ./...

fmt:
	gofmt -s -w -e .

test:
	go test -v -cover -timeout=120s -parallel=10 ./...

testacc:
	TF_ACC=1 HOSTINGER_HOST=http://localhost:1234 HOSTINGER_API_TOKEN=123 go test -v -cover -timeout 120m ./...

.PHONY: fmt lint test testacc build install generate
