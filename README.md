# sscache
golang sync.map cache

# ex
cache:=sscache.NewCache("ex")

cache.Set("test", "123", 0)//lifeSpan 0 means no TTL

cache.GetOrAdd("abc", "123", 5*time.Second)// if "abc" not found, set {"abc","123"}
