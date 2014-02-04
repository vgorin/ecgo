// Copyright 2013-2014 Vasiliy Gorin.
// Use of this source code is governed by a GNU-style
// license that can be found in the LICENSE file.

// Original Jerasure C/C++ code –
// Copyright 2007 James S. Plank
// See copyright notice inside *.c, *.h files

/*
 * LRU-Caches required for encoding/decoding routines
 */

package cauchy

import "sync"

import "github.com/vgorin/ecgo/util/lru"

// CauchyEncoderCache cache for CauchyEncoder structures
type CauchyEncoderCache struct {
	mutex sync.RWMutex
	cache *lru.Cache
}

// NewCauchyEncoderCache creates new CauchyEncoderCache with capacity specified
func NewCauchyEncoderCache(capacity int) *CauchyEncoderCache {
	return &CauchyEncoderCache{
		cache: lru.New(capacity),
	}
}

// DefaultCauchyEncoderCache default cache
var DefaultCauchyEncoderCache *CauchyEncoderCache = NewCauchyEncoderCache(64)

// Get gets the CauchyEncoder from the cache for the parameters given
// if there is no such encoder, creates new one and caches it
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
