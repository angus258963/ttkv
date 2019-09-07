# ttkv a thread safe key value store

### Install:
```go get github.com/angus258963/ttkv```

### Hot to use: 
```
import github.com/angus258963/ttkv
import github.com/angus258963/ttkv/cache

// create your store
store := ttkv.Store{
    get: func(key string) []byte {
        return []byte(key)
    },
    set: func(key string, value []byte) {},
}
// New KV with store and cache
kv := New(store, cache.NewCache(1*ttkv.MB, cache.FIFO))
  
v, err := kv.Get("123") // cache miss
v, err = kv.Get("123")  // cache hit

err = kv.Set("456", []byte("456"))
v, err = kv.Get("456")  // cache hit

```

### Test: 
```go test```
