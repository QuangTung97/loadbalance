.PHONY: all generate

all:
	go build -o server cmd/server/main.go
	go build -o client cmd/client/main.go

generate:
	protoc -I . --go_out=paths=source_relative:./balancepb --go-grpc_out=paths=source_relative:./balancepb balance.proto
