package internalcache_test

import (
	"bytes"
	"testing"
	"time"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/ptypes/struct"
	"github.com/pritesh-patel/redis-proxy/pkg/internalcache"
)

func dataGen(dataString string) *structpb.Value {
	testElement := []byte(dataString)
	data := &structpb.Value{}
	jsonpb.Unmarshal(bytes.NewReader(testElement), data)
	return data
}
func TestInsert(t *testing.T) {
	ttl := 1 * time.Minute
	ttlExpirationInterval := 5 * time.Minute
	cache := internalcache.NewSimpleInternalCache(1, ttl, ttlExpirationInterval)

	testElement := dataGen("bob")

	cache.InsertCachedItem("test", testElement)

	element, exists := cache.GetCachedItem("test")

	if exists == false {
		t.Error("Expected element to exist")
	}

	if element != testElement {
		t.Error("Expected test element value instead got", element)
	}

}

func TestTTLCachedEviction(t *testing.T) {
	ttl := 0 * time.Nanosecond
	ttlExpirationInterval := 1 * time.Minute

	cache := internalcache.NewSimpleInternalCache(1, ttl, ttlExpirationInterval)

	testElement := dataGen("bob")

	cache.InsertCachedItem("test", testElement)
	_, exists := cache.GetCachedItem("test")

	if exists == true {
		t.Error("Expected element not to exist.")
	}

}

func TestMaxCapCachedEviction(t *testing.T) {
	ttl := 1 * time.Duration(time.Minute)
	ttlExpirationInterval := 1 * time.Minute

	cache := internalcache.NewSimpleInternalCache(2, ttl, ttlExpirationInterval)

	testValues := []string{"test1", "test2", "test3"}
	for _, e := range testValues {

		cache.InsertCachedItem(e, dataGen(e))
	}

	_, exists := cache.GetCachedItem("test1")

	if exists == true {
		t.Error("Expected element should not to exist.")
	}

}
