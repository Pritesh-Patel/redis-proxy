package externalcache

import (
	"time"
	"bytes"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/ptypes/struct"
	"github.com/gomodule/redigo/redis"
	log "github.com/sirupsen/logrus"
)

type RedisCache struct {
	pool *redis.Pool
}

func initRedisPool(redisAddr string, maxIdleConns int) *redis.Pool {
	return &redis.Pool{
		MaxIdle:      maxIdleConns,
		IdleTimeout:  240 * time.Second,
		Dial:         func() (redis.Conn, error) { return dialConn(redisAddr) },
		TestOnBorrow: testConn,
	}
}

func dialConn(redisAddr string) (redis.Conn, error) {
	con, err := redis.Dial("tcp", redisAddr)
	if err != nil {
		return nil, err
	}
	return con, err
}

func testConn(con redis.Conn, t time.Time) error {
	if time.Since(t) < time.Minute {
		return nil
	}
	_, err := con.Do("PING")
	return err
}

func NewRedisCache(redisAddr string, maxIdleConns int) *RedisCache {
	pool := initRedisPool(redisAddr, maxIdleConns)

	return &RedisCache{pool}
}

func (c *RedisCache) GetCachedItem(key string) (*structpb.Value, bool) {

	con := c.pool.Get()
	result, err := redis.Bytes(con.Do("GET", key))

	if err != nil {
		log.WithError(err).Error("Failed to get item from redis.")
	}

	if result == nil {
		return nil, false
	}

	data := &structpb.Value{}
	jsonpb.Unmarshal(bytes.NewReader(result), data)

	return data, true
}