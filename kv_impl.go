package kv

import (
	"sync"

	"hash/fnv"

	"github.com/angus258963/ttkv/cache"
)

const (
	// shard is the number of locks
	shard = 256
)

type store struct {
	get func(key string) []byte
	set func(key string, value []byte)
}

type impl struct {
	cache cache.Cache
	store store
	// each lock is used for each bucket of keys
	locks []sync.RWMutex
}

func New(get func(key string) []byte) KV {
	return &impl{
		cache: cache.NewCache(1*MB, cache.FIFO),
		store: store{
			get: get,
			set: func(string, []byte) {},
		},
		locks: make([]sync.RWMutex, shard),
	}
}

func hash(s string) int {
	h := fnv.New32a()
	h.Write([]byte(s))
	return int(h.Sum32()) % shard
}

func (im *impl) Get(key string) ([]byte, error) {
	shardNum := hash(key)
	im.locks[shardNum].Lock()
	defer im.locks[shardNum].Unlock()

	v, err := im.cache.Get(key)
	if err != nil && err != cache.ErrNotFound {
		return nil, err
	} else if err == cache.ErrNotFound {
		// read through strategy
		v = im.store.get(key)
		im.cache.Set(key, v)
		return v, nil
	}

	return v, nil
}

func (im *impl) Set(key string, value []byte) error {
	shardNum := hash(key)
	im.locks[shardNum].Lock()
	defer im.locks[shardNum].Unlock()

	// write through strategy
	im.store.set(key, value)
	if err := im.cache.Set(key, value); err != nil {
		return err
	}

	return nil
}
