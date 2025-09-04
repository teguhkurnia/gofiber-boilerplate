default: dev

dev:
	go run cmd/web/main.go

.PHONY: test
test:
	go test -v ./test
