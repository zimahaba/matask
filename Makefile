.DEFAULT_GOAL := build

build:
	go build -o ./ ./...

run:
	go run cmd/matask/matask.go