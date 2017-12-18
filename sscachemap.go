package sscache

import (
	"sync"
	"time"
)

type CacheMap struct {
	name  string
	items sync.Map
	timer *time.Timer
}

func newCacheItem(key interface{}, data interface{}, lifeSpan time.Duration) *CacheItem {
	return &CacheItem{
		key,
		data,
		lifeSpan,
		time.Now()}
}

func (table *CacheMap) Range(trans func(key interface{}, value interface{}) bool) {
	table.items.Range(trans)
}

func (table *CacheMap) Set(key interface{}, value interface{}, lifeSpan time.Duration) {
	table.items.Store(key, newCacheItem(key, value, lifeSpan))
}

func (table *CacheMap) Get(key interface{}) (interface{}, bool) {
	v, ok := table.items.Load(key)
	if ok {
		return v.(*CacheItem).value, ok
	}
	return nil, ok
}

func (table *CacheMap) Delete(key interface{}) {
	table.items.Delete(key)
}

func (table *CacheMap) expirationCheck() {
	for {
		now := time.Now()
		keys := make([]interface{}, 0)
		table.items.Range(func(key, value interface{}) bool {
			item := value.(*CacheItem)
			if item.lifeSpan > 0 && now.Sub(item.createdOn) > item.lifeSpan {
				keys = append(keys, key)
			}
			return true
		})
		for _, key := range keys {
			table.items.Delete(key)
		}
		time.Sleep(100 * time.Millisecond)
	}
}
