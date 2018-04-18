PROTO_SOURCES := -I /usr/local/include
PROTO_SOURCES += -I .
PROTO_SOURCES += -I ${GOPATH}/src
PROTO_SOURCES += -I ${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis

gen:
	protoc ${PROTO_SOURCES} --go_out=plugins=grpc:. --grpc-gateway_out=logtostderr=true:. api/redis_proxy.proto

run:
	docker-compose build
	docker-compose up

test:
	docker-compose build
	docker-compose up -d
	docker-compose exec redis-proxy go test pkg/internalcache/simple_internal_cache_test.go
	@echo Waiting 5 seconds for redis and the redis proxy to come up...
	@sleep 5
	docker-compose exec redis redis-cli set it-test-1 1234
	docker-compose exec redis-proxy go test it/integration_test.go
	docker-compose down




