package sscache

import (
    "sync"
    "time"
)

type SSCacheMap struct {
    name  string
    items sync.Map
}

func (table *SSCacheMap) newSSCacheItem(key interface{}, data interface{}, lifeSpan time.Duration) *SSCacheItem {
    item := &SSCacheItem{
        key,
        data,
        lifeSpan,
        time.Now(), nil}

    item.timer = time.AfterFunc(item.lifeSpan, func() {
        if item.lifeSpan == 0 {
            return
        }
        for {
            if item.lifeSpan > 0 && time.Now().Sub(item.createdOn) > item.lifeSpan {
                table.Delete(item.key)
            }
            time.Sleep(100 * time.Millisecond)
        }
    })
    return item
}

func (table *SSCacheMap) Range(trans func(key interface{}, value interface{}) bool) {
    table.items.Range(trans)
}

func (table *SSCacheMap) Set(key interface{}, value interface{}, lifeSpan time.Duration) {
    item := table.newSSCacheItem(key, value, lifeSpan)
    table.items.Store(key, item)
}

func (table *SSCacheMap) GetOrAdd(key interface{}, value interface{}, lifeSpan time.Duration) interface{} {
    item := table.newSSCacheItem(key, value, lifeSpan)
    v, loaded := table.items.LoadOrStore(key, item)
    if loaded {
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
