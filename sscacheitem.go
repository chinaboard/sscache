package sscache

import "time"

type CacheItem struct {
	key interface{}
	value interface{}
	lifeSpan time.Duration
	createdOn time.Time
}

func (item *CacheItem) Key() interface{} {
	return item.key
}

func (item *CacheItem) Value() interface{} {
	return item.value
}

func (item *CacheItem) LifeSpan() time.Duration {
	return item.lifeSpan
}

func (item *CacheItem) CreatedOn() time.Time {
	return item.createdOn
}