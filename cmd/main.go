package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/pritesh-patel/redis-proxy/api"
	"github.com/pritesh-patel/redis-proxy/pkg/server"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/pritesh-patel/redis-proxy/pkg/externalcache"
	"github.com/pritesh-patel/redis-proxy/pkg/internalcache"

	"github.com/spf13/viper"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

var sigwatcher = make(chan os.Signal)

func main() {
	log.Info("Starting up redis proxy...")

	
	// default config
	viper.SetDefault("GRPC_HOST","localhost")
	viper.SetDefault("GRPC_PORT", 3000)
	viper.SetDefault("HTTP_HOST","localhost")
	viper.SetDefault("HTTP_PORT", 8080)
	viper.SetDefault("REDIS_HOST","localhost")
	viper.SetDefault("REDIS_PORT", 6379)
	viper.SetDefault("LOCAL_CACHE_MAX_SIZE", 20)
	viper.SetDefault("LOCAL_CACHE_TTL_SECONDS", "1m")
	viper.SetDefault("LOCAL_CACHE_TTL_CLEAN_INTERVAL_SECONDS", "1s")
	viper.SetDefault("REDIS_MAX_IDLE_CONNECTIONS", 3)

	viper.AutomaticEnv()

	grpcHost := viper.GetString("GRPC_HOST")
	grpcPort := viper.GetInt("GRPC_PORT")
	httpHost := viper.GetString("HTTP_HOST")
	httpPort := viper.GetInt("HTTP_PORT")
	redisHost := viper.GetString("REDIS_HOST")
	redisPort := viper.GetInt("REDIS_PORT")

	localCahceMaxSize := viper.GetInt("LOCAL_CACHE_MAX_SIZE")
	localCahceTTL := viper.GetDuration("LOCAL_CACHE_TTL_SECONDS")
	localCacheCleanInterval := viper.GetDuration("LOCAL_CACHE_TTL_CLEAN_INTERVAL_SECONDS")

	redisIdleConnections := viper.GetInt("REDIS_MAX_IDLE_CONNECTIONS")

	grpcAddr := fmt.Sprintf("%s:%d", grpcHost, grpcPort)
	httpAddr := fmt.Sprintf("%s:%d", httpHost, httpPort)
	redisAddr := fmt.Sprintf("%s:%d", redisHost, redisPort)

	ic := internalcache.NewSimpleInternalCache(localCahceMaxSize, localCahceTTL, localCacheCleanInterval)
	ec := externalcache.NewRedisCache(redisAddr, redisIdleConnections)

	proxyServer := server.NewRedisProxyServer(ic, ec)

	go runGRPC(grpcAddr, proxyServer)
	go runHTTP(grpcAddr, httpAddr)

	log.Infof("HTTP server running on %v", httpAddr)
	log.Infof("GRPC server running on %v", grpcAddr)

	signal.Notify(sigwatcher, syscall.SIGTERM)
	signal.Notify(sigwatcher, syscall.SIGINT)

	shutdownHook()

}

func runGRPC(grpcAddr string, redisProxyServer *server.RedisProxyServer) {
	lis, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		log.Errorf("Failed to bind: %v", err)
	}
	grpcServer := grpc.NewServer()
	redisproxy.RegisterRedisProxyServiceServer(grpcServer, redisProxyServer)
	grpcServer.Serve(lis)
}

func runHTTP(grpcAddr string, httpAddr string) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	opts := []grpc.DialOption{grpc.WithInsecure()}
	mux := runtime.NewServeMux()

	err := redisproxy.RegisterRedisProxyServiceHandlerFromEndpoint(ctx, mux, grpcAddr, opts)

	if err != nil {
		log.Fatalf("Failed to start HTTP server: %v", err)
	}
	http.ListenAndServe(httpAddr, mux)
}

func getEnvOr(envKey string, defaultValue string) string {
	value, exists := os.LookupEnv(envKey)
	if exists {
		return value
	}
	return defaultValue
}

func shutdownHook() {
	sig := <-sigwatcher
	fmt.Printf("Shutting down. SIGTERM Recieved: %v", sig)
	os.Exit(0)
}
