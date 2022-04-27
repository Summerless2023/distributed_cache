package main

import (
	"main/src/concurrency"
	"main/src/models"
	"sync"
)

type DSCache struct {
	name      string
	mainCache concurrency.ConcurrencyCache
}

var (
	mu          sync.RWMutex
	DSCachesMap = make(map[string]*DSCache)
)

func NewDSCache(name string) *DSCache {
	mu.Lock()
	defer mu.Unlock()
	g := &DSCache{
		name:      name,
		mainCache: *concurrency.NewConcurrencyCache(),
	}
	DSCachesMap[name] = g
	return g
}

func GetDSCache(name string) *DSCache {
	mu.RLock()
	g := DSCachesMap[name]
	mu.RUnlock()
	return g
}

// func (g *DSCache) Get(key string) (string, error) {
// 	if key == "" {
// 		return "", fmt.Errorf("key is required")
// 	}

// 	if v, ok := g.mainCache.Get(models.KeyType(key)); ok {
// 		log.Println("[GeeCache] hit")
// 		return string(v), nil
// 	}

// 	return g.load(key)
// }

// func (g *DSCache) load(key string) (value string, err error) {
// 	return g.getLocally(key)
// }

// func (g *DSCache) getLocally(key string) (string, error) {
// 	value, err := g.getter.Get(key)
// 	if err != nil {
// 		return "", err

// 	}
// 	value := value
// 	g.populateCache(key, value)
// 	return value, nil
// }

func (g *DSCache) populateCache(key string, value string) {
	g.mainCache.Add(models.KeyType(key), models.ValueType(value))
}
