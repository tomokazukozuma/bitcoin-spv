export GO111MODULE=on

COMMAND=balance

build:
	go build -o cmd/exe ./cmd

run:
	go run cmd/main.go ${COMMAND}
