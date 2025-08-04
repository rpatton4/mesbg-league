# © 2025 Robert Patton robpatton@infiniteskye.com

.DEFAULT_GOAL := build

.PHONY:fmt vet build

cyclo:
	-gocyclo -over 10 -ignore \s*_test.go .

fmt:
	go fmt ./...

vet: fmt
	go vet ./...

lint: cyclo
	go vet -vettool=$(which shadow) ./...
	errcheck ./...

generate:
	go generate ./...

test: generate
	go test ./...

build: vet
	go build -o ./fsa

# run:
#	func start
