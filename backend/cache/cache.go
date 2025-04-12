package cache

import (
	"errors"
	"sync"
)

var cacheMap sync.Map

func Get(key string) (any, error) {
	value, ok := cacheMap.Load(key)
	if !ok {
		return nil, errors.New("key not found")
	}
	return value, nil
}

func Set(key string, value any) error {
	cacheMap.Store(key, value)
	return nil
}

func Delete(key string) error {
	cacheMap.Delete(key)
	return nil
}
