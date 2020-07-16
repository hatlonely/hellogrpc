export GOPROXY=https://goproxy.io

build: cmd/*/main.go
	mkdir -p build/bin
	go build cmd/server/main.go && mv server build/bin
	go build cmd/client/main.go && mv client build/bin

codegen:
	protoc -I. --go_out=plugins=grpc,paths=source_relative:api/gen/go/ api/echo.proto
	protoc -I. --grpc-gateway_out=logtostderr=true,paths=source_relative:api/gen/go api/echo.proto
