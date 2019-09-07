package cache

import (
	"sync"
)

const (
	// maxValueSize maximum size of item(key+value)
	maxValueSize = 1 * 1024 * 1024 // 1MB
)

type cacheImpl struct {
	strategy strategy
	buffer   map[string][]byte
	// lock protects concurrent map writes
	lock sync.RWMutex
}

type strategy interface {
	get(buffer map[string][]byte, key string) ([]byte, error)
	set(buffer map[string][]byte, key string, value []byte)
}

// NewCache new a cache with maxCapcity and strategy
func NewCache(maxCapcity int, strategy StrategyType) Cache {
	b := make(map[string][]byte)

	return &cacheImpl{
		buffer:   b,
		strategy: newStrategy(maxCapcity, strategy),
	}
}

func (im *cacheImpl) Get(key string) ([]byte, error) {
	im.lock.RLock()
	defer im.lock.RUnlock()

	v, err := im.strategy.get(im.buffer, key)
	if err != nil {
		return nil, ErrNotFound
	}

	return v, nil
}

func (im *cacheImpl) Set(key string, value []byte) error {
	im.lock.Lock()
	defer im.lock.Unlock()

	// check input item is no more than maxValueSize
	if len(key)+len(value) > maxValueSize {
		return ErrMaxValueSize
	}

	im.strategy.set(im.buffer, key, value)

	return nil
}

func newStrategy(maxCapcity int, strategy StrategyType) strategy {
	switch strategy {
	case FIFO:
		return newFIFO(maxCapcity)
	}

	panic("no mapping strategy")
}

type fifoStrategy struct {
	queue      []string
	size       int
	maxCapcity int
}

func newFIFO(maxCapcity int) *fifoStrategy {
	return &fifoStrategy{
		queue:      []string{},
		size:       0,
		maxCapcity: maxCapcity,
	}
}

func (f *fifoStrategy) get(buffer map[string][]byte, key string) ([]byte, error) {
	v, ok := buffer[key]
	if !ok {
		return nil, ErrNotFound
	}

	return v, nil
}

func (f *fifoStrategy) set(buffer map[string][]byte, key string, itemValue []byte) {
	v, ok := buffer[key]
	if ok {
		// remove old item
		delete(buffer, key)
		f.size -= (len(v) + len(key))
	}

	sizeOfItem := len(key) + len(itemValue)
	if sizeOfItem+f.size > f.maxCapcity {
		// clear item from queue
		for i, k := range f.queue {
			v := buffer[k]
			delete(buffer, k)
			f.size -= (len(v) + len(k))

			// check deleted items are enough
			if sizeOfItem+f.size <= f.maxCapcity && k != key {
				f.queue = f.queue[i+1:]
				break
			}
		}
	}

	// make sure no duplicated key in queue
	if !ok {
		f.queue = append(f.queue, key)
	}
	buffer[key] = itemValue
	f.size += sizeOfItem
}
