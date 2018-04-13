package internalcache

import (
	"time"
	
	"github.com/golang/protobuf/ptypes/struct"
)


type InternalCache interface {
	GetCachedItem(key string) (*structpb.Value, bool)
	InsertCachedItem(key string, item *structpb.Value)
}

type CachedItem struct {
	Key string
	Item *structpb.Value
	Timestamp time.Time
}