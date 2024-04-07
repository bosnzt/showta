package memcache

import (
	"github.com/patrickmn/go-cache"
	"time"
)

var (
	c *cache.Cache
)

const (
	List = "list:"
	Link = "link:"
)

func init() {
	c = cache.New(5*time.Minute, 10*time.Minute)
}

func Get(prefix string, k string) (interface{}, bool) {
	return c.Get(getKey(prefix, k))
}

func Set(prefix string, k string, x interface{}) {
	c.Set(getKey(prefix, k), x, cache.DefaultExpiration)
}

func Expire(prefix string, k string, x interface{}, d time.Duration) {
	if d == 0 {
		d = cache.DefaultExpiration
	}
	c.Set(getKey(prefix, k), x, d)
}

func getKey(prefix string, k string) string {
	return prefix + k
}
