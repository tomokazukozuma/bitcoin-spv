export GO111MODULE=on

build:
	go build -o cmd/exe ./cmd

run:
	go run cmd/main.go
