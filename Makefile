PROTO_SOURCES := -I /usr/local/include
PROTO_SOURCES += -I .
PROTO_SOURCES += -I ${GOPATH}/src
PROTO_SOURCES += -I ${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis

gen:
	protoc ${PROTO_SOURCES} --go_out=plugins=grpc:. --grpc-gateway_out=logtostderr=true:. api/redis_proxy.proto

run:
	go run cmd/main.go

run-with-redis:
	docker-compose up

build:
	dep ensure
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/app cmd/main.go

build-container: build
	docker build -t "redis-proxy:latest" .

test: build-container
	go test pkg/internalcache/simple_internal_cache_test.go
	docker.compose up -d
	docker.compose exec redis redis-cli set it-test-1 1234
	go test it/integration_test.go
	docker.compose down




