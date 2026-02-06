.PHONY: all proto server agent docker clean

all: proto server agent

proto:
	protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    pkg/pb/c2.proto

server:
	go build -o bin/server cmd/server/main.go

agent:
	go build -o bin/agent cmd/agent/main.go

docker:
	docker build -t titan-c2 .

clean:
	rm -rf bin/
