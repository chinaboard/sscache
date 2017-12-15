package sscache

import "time"

type SSCacheItem struct {
	key interface{}
	value interface{}
	lifeSpan time.Duration
	createdOn time.Time
}

func (item *SSCacheItem) Key() interface{} {
	return item.key
}

func (item *SSCacheItem) Value() interface{} {
	return item.value
}

func (item *SSCacheItem) LifeSpan() time.Duration {
	return item.lifeSpan
}

func (item *SSCacheItem) CreatedOn() time.Time {
	return item.createdOn
}