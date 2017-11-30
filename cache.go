package sscache

import (
	"sync"
)

func NewCache(name string) *SSCacheMap {

	value := &SSCacheMap{name, sync.Map{}}

	return value
}
