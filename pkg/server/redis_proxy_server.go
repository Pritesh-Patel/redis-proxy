package server

import (
	"github.com/pritesh-patel/redis-proxy/api"
	"github.com/pritesh-patel/redis-proxy/pkg/externalcache"
	"github.com/pritesh-patel/redis-proxy/pkg/internalcache"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

type RedisProxyServer struct {
	internalCache internalcache.InternalCache
	externalcache externalcache.ExternalCache
}

func NewRedisProxyServer(internalCache internalcache.InternalCache, externalCache externalcache.ExternalCache) *RedisProxyServer {
	return &RedisProxyServer{
		internalCache: internalCache,
		externalcache: externalCache,
	}
}

func (s *RedisProxyServer) GetItem(ctx context.Context, req *redisproxy.GetItemRequest) (*redisproxy.GetItemResponse, error) {
	itemLocal, existsLocal := s.internalCache.GetCachedItem(req.Id)

	logger := log.WithField("item_key", req.Id)

	if existsLocal {
		logger.Info("Item found in local cache.")
		return &redisproxy.GetItemResponse{itemLocal}, nil
	}

	itemExternal, existExternal := s.externalcache.GetCachedItem(req.Id)

	if existExternal {
		logger.Info("Item found in external cache.")
		go s.internalCache.InsertCachedItem(req.Id, itemExternal)
		return &redisproxy.GetItemResponse{itemExternal}, nil
	}

	return nil, grpc.Errorf(codes.NotFound, "Item with that key was not found")
}
