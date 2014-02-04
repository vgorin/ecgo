package jerasurego

import "sync"

import "github.com/golang/groupcache/lru"

type CauchyEncoderCache struct {
	mutex sync.RWMutex
	cache *lru.Cache
}

func NewCauchyEncoderCache(capacity int) *CauchyEncoderCache {
	return &CauchyEncoderCache{
		cache: lru.New(capacity),
	}
}

var DefaultCauchyEncoderCache *CauchyEncoderCache = NewCauchyEncoderCache(64)

func (c *CauchyEncoderCache) Get(p encoder_params) *CauchyEncoder {
	if encoder, ok := c.get(p); ok {
		return encoder
	}
	encoder := NewCauchyEncoder(p)

	c.put(p, encoder)

	return encoder
}

// get is thread-safe function
func (c *CauchyEncoderCache) get(p encoder_params) (*CauchyEncoder, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	if encoder, ok := c.cache.Get(p); ok {
		return encoder.(*CauchyEncoder), ok
	}
	return nil, false
}

// put is thread-safe function
func (c *CauchyEncoderCache) put(p encoder_params, encoder *CauchyEncoder) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.cache.Add(p, encoder)
}
