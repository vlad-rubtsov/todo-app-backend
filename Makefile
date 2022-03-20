.PHONY: build
build:
	go build -v cmd/main.go

.PHONY: test
test:
	go test -v ./...
