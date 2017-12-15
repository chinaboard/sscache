package sscache

import (
    "sync"
    "time"
)

func NewCache(name string) *CacheMap {

    value := &CacheMap{name, sync.Map{}, nil}

    value.timer=time.AfterFunc(0*time.Second, func() {
        go value.expirationCheck()
    })

    return value
}
