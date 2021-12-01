package caches

import (
	"github.com/hashicorp/golang-lru"
)

var (
	cacheManagerInstance CacheManager
)

func GetCacheManager() (instance CacheManager, err error) {
	if cacheManagerInstance == nil {
		instance, err = newCacheManager()
		if err != nil {
			return
		}
		cacheManagerInstance = instance
	}
	return cacheManagerInstance, nil
}

func newCacheManager() (instance CacheManager, err error) {
	return &cacheManager{}, nil
}

type CacheManager interface {
	GetCacheNames() ([]string, error)
	GetCache(cacheName string) (*lru.Cache, error)
}

type cacheManager struct {
	cacheMap map[string]*lru.Cache
}

func (manager *cacheManager) GetCacheNames() (names []string, err error) {
	for key := range manager.cacheMap {
		names = append(names, key)
	}
	return names, nil
}

func (manager *cacheManager) GetCache(cacheName string) (cache *lru.Cache, err error) {
	cache = manager.cacheMap[cacheName]
	if cache == nil {
		cache, err = lru.New(128)
		if err != nil {
			return nil, err
		}
	}
	return cache, nil
}
