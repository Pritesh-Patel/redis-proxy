package externalcache
import "github.com/golang/protobuf/ptypes/struct"

type ExternalCache interface {
	GetCachedItem(key string) (*structpb.Value, bool)
}