package internalcache

import (
	"container/list"
	"sync"
	"time"

	"github.com/golang/protobuf/ptypes/struct"
	log "github.com/sirupsen/logrus"
)

type SimpleInternalCache struct {
	maxSize int
	cache   map[string]*CachedItem
	tracker *list.List
	ttl     time.Duration
	lock    *sync.Mutex
}

func NewSimpleInternalCache(size int, ttl time.Duration, ttlCleanUpInterval time.Duration) *SimpleInternalCache {
	cache := make(map[string]*CachedItem)
	evictionList := list.New()
	lock := &sync.Mutex{}

	simpleCache := &SimpleInternalCache{size, cache, evictionList, ttl, lock}
	
	go simpleCache.deathWatch(ttlCleanUpInterval)

	return simpleCache
}

func (c *SimpleInternalCache) InsertCachedItem(key string, item *structpb.Value) {

	cachedItem := &CachedItem{key, item, time.Now()}
	c.cache[key] = cachedItem
	c.tracker.PushFront(cachedItem)

	if c.tracker.Len() > c.maxSize {
		c.removeLFUItem()
	}
}

func (c *SimpleInternalCache) GetCachedItem(key string) (*structpb.Value, bool) {
	cachedItem := c.cache[key]

	if cachedItem != nil && !c.isExpired(cachedItem) {
		defer c.lock.Unlock()

		c.lock.Lock()
		c.tracker.PushFront(cachedItem)

		return cachedItem.Item, true
	}

	return nil, false
}

func (c *SimpleInternalCache) deathWatch(cleaningInterval time.Duration) {
	ticker := time.NewTicker(cleaningInterval)

	for t := range ticker.C {
		log.Debug(t)
		c.evictExpiredItems()
	}

}

func (c *SimpleInternalCache) isExpired(item *CachedItem) bool {

	return item.Timestamp.Add(c.ttl).Before(time.Now())
}

func (c *SimpleInternalCache) evictExpiredItems() {

	for e := c.tracker.Front(); e != nil; e = e.Next() {
		item := e.Value.(*CachedItem)

		if c.isExpired(item) {

			c.lock.Lock()

			delete(c.cache, item.Key)
			c.tracker.Remove(e)

			c.lock.Unlock()
		}
	}
}

func (c *SimpleInternalCache) removeLFUItem() {

	defer c.lock.Unlock()

	lastElement := c.tracker.Back()
	lastItem := lastElement.Value.(*CachedItem)

	c.lock.Lock()

	delete(c.cache, lastItem.Key)
	c.tracker.Remove(lastElement)
}
