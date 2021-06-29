  
SYSLFILE = gop.sysl
APPS = gop

.PHONY: run test lint

run:
	go run . config.yaml

test:
	go test ./... -short

lint:
	golangci-lint run ./...
