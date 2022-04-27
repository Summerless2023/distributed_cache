package main

import (
	"fmt"
	"main/src/concurrency"
	"main/src/models"
	"sync"

	"github.com/sirupsen/logrus"
)

type DSCache struct {
	name      string                       //缓存的名称
	mainCache concurrency.ConcurrencyCache //主缓存
	getter    models.Getter                //回调函数
}

var (
	mu          sync.RWMutex
	DSCachesMap = make(map[string]*DSCache) //保存所有缓存的map，key是缓存名称，value是DSCache指针
)

func NewDSCache(name string, getter models.Getter) *DSCache {
	if getter == nil {
		panic("nil Getter")
	}
	mu.Lock()
	defer mu.Unlock()
	g := &DSCache{
		name:      name,
		mainCache: *concurrency.NewConcurrencyCache(),
		getter:    getter,
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

func (g *DSCache) Get(key string) (string, error) {
	if key == "" {
		return "", fmt.Errorf("key is required")
	}

	if v, ok := g.mainCache.Get(models.KeyType(key)); ok {
		logrus.Debug("[GeeCache] hit")
		return string(v), nil
	}

	return g.load(key)
}

func (g *DSCache) load(key string) (value string, err error) {
	return g.getLocally(key)
}

func (g *DSCache) getLocally(key string) (string, error) {
	value, err := g.getter.Get(models.KeyType(key))
	if err != nil {
		return "", err

	}
	g.populateCache(key, string(value))
	return string(value), nil
}

func (g *DSCache) populateCache(key string, value string) {
	g.mainCache.Add(models.KeyType(key), models.ValueType(value))
}
