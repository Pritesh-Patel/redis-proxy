# Redis Proxy

This is a redis proxy with a configurable LFU internal cache.

This application fetches keyed data from redis and then caches it itself so any further request using that key will be returned directly via the application if it has not been expired in the local cache.

## Architecture
The architecture is very simple:
  - The internal and external cache both have separate interfaces allowing either to be re-written
  - These interfaces are used within a server interface that implements the protobuf spec
  - There is a GRPC server for grpc and a grpc-gateway server to provide http translation

## Proxy Communication

### GRPC 
  - By default available on 0.0.0.0:3000
  - You can use grpcc or another grpc client to "curl" the service, a grpc client can be generated from proto in the api folder

### HTTP 
  - By default available on 0.0.0.0:8080/v1/proxy/
  - An example request `curl 0.0.0.0:8080/v1/proxy/{KEY_NAME}`
  - All data is marshaled to json by the GRPC gateway data types stored in redis will be respected including objects. ie. strings, ints, objects will be represented in there JSON form

## Instructions
  - `make run` will start up the application with redis via docker compose
  - `make test` will start up the application with redis in docker compose run unit and integration tests. For the purpose of testing some test data is inserted into redis on start up of the application. After the tests have been run it will shutdown redis and the application.
  - The application can be configured by environment variables, you can see the environment variables and there defaults in `cmd/main.go`
  - Configuration can be overridden by adding environment variables to the `redis-proxy.env` file, if you change the ports for either grpc / http make sure to update the docker-compose.yml to reflect this.

## Internal caching details
Then internal cache uses a hashmap to store the values and a doubly linked list to track the priority of the item

### Retrieving

- Every time an item is requested it is pushed to the front of the priority queue.
- If the cache has reached the maximum capacity the last item (the back) in the priority queue will be deleted from both the hash map and the priorty queue
- The operation from fetching from the cache has the functional complexity of O(1)

### TTL Eviction

- When the cache is initialized it starts a `deathWatch` method in a go routine, on configured interval this method will scan the hashmap and delete any items that have expired
- Even if the `deathWatch` has not been run in time and a get request requests an expired item it will not be returned as the ttl is checked before returning the item. This item will eventually be cleaned up when the method next runs
- The functional complexity for deleting expired items is O(n) as it iterates through the list to find and delete the element

### Concurrency
- Concurrency in the application is controlled by locks, the critical sections (writing to the cache) lock at the beginning of the operation and unlock at the end.
- The redis connection is pooled with a default of 3 connections, this can be changed with config.



