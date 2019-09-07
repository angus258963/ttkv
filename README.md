# ttkv a thread safe key value store

### Install:
```go get github.com/angus258963/ttkv```

### Hot to use: 
```
import github.com/angus258963/ttkv

// max capacity is 1 MB
// default strategy is FIFO
kv := ttkv.New(func(key string) []byte {
		return []byte(key)
})
  
v := kv.Get("123") // cache miss
v = kv.Get("123") // cache hit
kv.Set("456", []byte("456"))
v = kv.Get("456") // cache hit
```

### Test: 
```go test```
