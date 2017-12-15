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

func (table *CacheMap) newCacheItem(key interface{}, data interface{}, lifeSpan time.Duration) *CacheItem {
    item := &CacheItem{
        key,
        data,
        lifeSpan,
        time.Now()}

    return item
}

func (table *CacheMap) Range(trans func(key interface{}, value interface{}) bool) {
    table.items.Range(trans)
}

func (table *CacheMap) Set(key interface{}, value interface{}, lifeSpan time.Duration) {
    item :=  func() *CacheItem { return table.newCacheItem(key, value, lifeSpan) }
    table.items.Store(key, item)
}

func (table *CacheMap) set(key interface{}, value *CacheItem) {
    item :=  func() *CacheItem { return value }
    table.items.Store(key, item)
}

func (table *CacheMap) GetOrAdd(key interface{}, value interface{}, lifeSpan time.Duration) interface{} {
    v,_:=table.lazyLoad(key,func()*CacheItem{return table.newCacheItem(key,value,lifeSpan)})
    return v.value
}

func (table *CacheMap) GetOrUpdate(key interface{}, value interface{}, lifeSpan time.Duration) interface{} {
    v,get:=table.lazyLoad(key,func()*CacheItem{return table.newCacheItem(key,value,lifeSpan)})
    if get{
        v.createdOn=time.Now()
        table.set(key,v)
    }
    return v.value
}

func (table *CacheMap) Get(key interface{}) (interface{}, bool) {
    v, ok := table.items.Load(key)
    if ok {
        return v.(func() *CacheItem)().value, ok
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
            item := value.(func() *CacheItem)()
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

func (table *CacheMap) lazyLoad(key interface{}, f func() *CacheItem) (*CacheItem,bool) {
    if g, ok := table.items.Load(key); ok {
        return g.(func() *CacheItem)(),ok
    }

    var (
        once sync.Once
        x *CacheItem
    )
    g, _ := table.items.LoadOrStore(key, func() *CacheItem {
        once.Do(func() {
            x = f()
            table.items.Store(key, func() *CacheItem { return x })
        })
        return x
    })
    return g.(func() *CacheItem)(),false
}
