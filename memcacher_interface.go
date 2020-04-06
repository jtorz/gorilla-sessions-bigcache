package gsb

import (
	"github.com/allegro/bigcache"
)

// BigCacher is the interface gsb uses to interact with the bigcache client
type BigCacher interface {
	Get(key string) (string, error)
	Set(key, val string, exp uint32, ocas uint64) (cas uint64, err error)
}

// GoBigcacher is a wrapper to the gobigcache client that implements the
// BigCacher interface
type GoBigcacher struct {
	client *bigcache.BigCache
}

// NewGoBigcacher returns a wrapped gobigcache client that implements the
// BigCacher interface
func NewGoBigcacher(c *bigcache.BigCache) *GoBigcacher {
	if c == nil {
		panic("Cannot have nil bigcache client")
	}
	return &GoBigcacher{client: c}
}

func (gm *GoBigcacher) Get(key string) (string, error) {
	if v, err := gm.client.Get(key); err == nil {
		return string(v), err
	} else {
		return "", err
	}
}

func (gm *GoBigcacher) Set(key, val string, exp uint32, ocas uint64) (cas uint64, err error) {
	err = gm.client.Set(key, []byte(val))
	return ocas, err
}
