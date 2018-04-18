FROM golang:1.10.1
WORKDIR /go/src/github.com/pritesh-patel/redis-proxy
COPY . .
CMD ["go", "run", "cmd/main.go"]