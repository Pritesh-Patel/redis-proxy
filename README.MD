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
  - By default availble on localhost:3000
  - You can use grpcc or another grpc client to "curl" the service, a grpc client can be generated from proto in the api folder

### HTTP 
  - By default availble on localhost:8080/v1/proxy/
  - An example request `curl localhost:8080/v1/proxy/{KEY_NAME}`
  - All data is marshalled to json by the GRPC gateway data types stored in redis will be respected including objects. ie. strings, ints, objects will be represented in there JSON form

## Instructions
  - `make run-with-redis` will run the application with redis in docker compose
  - `make run` will build and run the main.go file
  - `make test` will start up the application with redis in docker compose run unit and integration tests. For the purpose of testing some test data is inserted into redis on start up of the application. It will the clean up by shutting down the docker containers.
  - The application can be configured by environment variables, you can see the environment variables and there defaults in `cmd/main.go`
  - Config can be override by adding environment variables to the `docker-compose.yml` file

## Internal caching details
Then internal cache uses a hashmap to store the values and a doubly linked list to track the priority of the item

### Retrieving

- Everytime an item is requested it is pushed to the front of the priority queue.
- If the cache has reached the maximum capacity the last item (the back) in the priority queue will be deleted from both the hash map and the priorty queue
- The operation from fetching from the cache has the functional complexity is O(1)

### TTL Eviction

- When the cache is initialized it starts a `deathWatch` method in a go routine, on configured interval this method will scan the hashmap and delete any items that have expired
- Even if the `deathWatch` has not been run in time and a get request requests an expired item it will not be returned as the ttl is checked before returning the item. This item will eventually be cleaned up when the method next runs
- The functional complexity for deleting expired items is O(n) as it iterates through the list to find and delete the element

### Concurrency
- Concurrency in the application is controlled by locks, the critical sections (writing to the cache) lock at the beginning of the operation and unlock at the end.
- The redis connection is pooled with a default of 3 connections, this can be changed with config.

## Time spent (roughly)
- Research (1 hour)
- GRPC & HTTP setup (5 mins)
- Local Cache (25 mins)
- Redis Cache (10 mins)
- Marshalling from redis (40 mins)




