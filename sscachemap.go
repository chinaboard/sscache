package sscache

import (
	"sync"
	"time"
)

type SSCacheMap struct {
	name  string
	items sync.Map

	timer *time.Timer
}

func (table *SSCacheMap) Range(trans func(key interface{}, value interface{}) bool) {
	table.items.Range(trans)
}

func (table *SSCacheMap) Set(key interface{}, value interface{}, lifeSpan time.Duration) {
	item := NewSSCacheItem(key, value, lifeSpan)
	table.items.Store(key, item)
}

func (table *SSCacheMap) GetOrAdd(key interface{}, value interface{}, lifeSpan time.Duration) interface{} {
	item := NewSSCacheItem(key, value, lifeSpan)

	v, notFound := table.items.LoadOrStore(key, item)

	if !notFound {
		v.(*SSCacheItem).createdOn = time.Now()
	}

	return v.(*SSCacheItem).value
}

func (table *SSCacheMap) Get(key interface{}) (interface{}, bool) {
	v, ok := table.items.Load(key)
	if ok {
		return v.(*SSCacheItem).value, ok
	}
	return nil, ok
}

func (table *SSCacheMap) Delete(key interface{}) {
	table.items.Delete(key)
}

func (table *SSCacheMap) expirationCheck() {

	for {
		now := time.Now()

		keys := make([]interface{}, 0)
		table.items.Range(func(key, value interface{}) bool {

			item := value.(*SSCacheItem)
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
