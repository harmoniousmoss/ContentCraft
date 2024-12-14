// Configuration for In-Memory Cache and Synchronization

package config

import (
	"sync"
	"time"

	"github.com/patrickmn/go-cache" //A third-party library for in-memory caching with expiration capabilities.
)

var (
	CacheInstance   = cache.New(5*time.Minute, 10*time.Minute)
	FetchLock       sync.Mutex
	FetchInProgress = new(sync.Map) // Using sync.Map for concurrent map operations
)
