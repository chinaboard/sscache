package sscache

import (
    "sync"
)

func NewCache(name string) *SSCacheMap {

    value := &SSCacheMap{name, sync.Map{}, nil}

    value.timer = time.AfterFunc(0*time.Second, func() {
        go value.expirationCheck()
    })
    return value
}
