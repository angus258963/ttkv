package kv

import (
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	"github.com/angus258963/ttkv/cache"
)

type KVSuite struct {
	suite.Suite
	kv KV
}

func (s *KVSuite) SetupSuite() {
}

func (s *KVSuite) SetupTest() {
}

func (s *KVSuite) TestGet() {
	store := Store{
		get: func(key string) []byte {
			time.Sleep(time.Millisecond * 100)
			return []byte(key)
		},
		set: func(key string, value []byte) {},
	}
	s.kv = New(store, cache.NewCache(1*MB, cache.FIFO))

	// case: no cache
	beforeTime := time.Now()
	for i := 0; i < 10; i++ {
		key := strconv.Itoa(i)
		b, err := s.kv.Get(key)
		s.Require().NoError(err)
		s.Require().Equal(key, string(b))
	}
	afterTime := time.Now()
	s.Require().Equal(true, beforeTime.Add(time.Second).Before(afterTime))

	// case: cached, use timer to check
	beforeTime = time.Now()
	for i := 0; i < 10; i++ {
		key := strconv.Itoa(i)
		b, err := s.kv.Get(key)
		s.Require().NoError(err)
		s.Require().Equal(key, string(b))
	}
	afterTime = time.Now()
	s.Require().Equal(false, beforeTime.Add(time.Second).Before(afterTime))
}

func (s *KVSuite) TestFIFO() {
	store := Store{
		get: func(key string) []byte {
			time.Sleep(time.Millisecond * 100)
			b := make([]byte, 256*1024-len(key))
			copy(b, key)
			return b
		},
		set: func(key string, value []byte) {},
	}
	s.kv = New(store, cache.NewCache(1*MB, cache.FIFO))

	// init cache
	for i := 0; i < 4; i++ {
		key := strconv.Itoa(i)
		s.kv.Get(key)
	}
	// case: cached, use time to check
	beforeTime := time.Now()
	for i := 0; i < 4; i++ {
		key := strconv.Itoa(i)
		b, err := s.kv.Get(key)
		s.Require().NoError(err)
		exp := make([]byte, 256*1024-len(key))
		copy(exp, key)
		s.Require().Equal(exp, b)
	}
	afterTime := time.Now()
	s.Require().Equal(false, beforeTime.Add(time.Millisecond*100).Before(afterTime))

	// case: cache miss, [1,2,3,4] to get 0, 1, 2, 3
	// set 4
	key := strconv.Itoa(4)
	s.kv.Get(key)

	beforeTime = time.Now()
	for i := 0; i < 4; i++ {
		key = strconv.Itoa(i)
		b, err := s.kv.Get(key)
		s.Require().NoError(err)
		exp := make([]byte, 256*1024-len(key))
		copy(exp, key)
		s.Require().Equal(exp, b)
	}
	afterTime = time.Now()
	s.Require().Equal(true, beforeTime.Add(time.Millisecond*400).Before(afterTime))
}

func (s *KVSuite) TestMultiThreadsGet() {
	var miss int
	store := Store{
		get: func(key string) []byte {
			b := make([]byte, 1*1024-len(key))
			copy(b, key)
			miss++
			return b
		},
		set: func(key string, value []byte) {},
	}
	s.kv = New(store, cache.NewCache(1*MB, cache.FIFO))

	wg := sync.WaitGroup{}
	for t := 0; t < 100; t++ {
		wg.Add(1)
		go func() {
			for i := 0; i < 100; i++ {
				key := strconv.Itoa(i)
				s.kv.Get(key)
			}
			wg.Done()
		}()
	}
	wg.Wait()
	s.Require().Equal(100, miss)
}

func (s *KVSuite) TestMultiThreadsSet() {
	store := Store{
		get: func(key string) []byte {
			return []byte("test")

		},
		set: func(key string, value []byte) {},
	}
	s.kv = New(store, cache.NewCache(1*MB, cache.FIFO))

	wg := sync.WaitGroup{}
	for t := 0; t < 100; t++ {
		wg.Add(1)
		go func() {
			for i := 0; i < 100; i++ {
				key := strconv.Itoa(i)
				s.kv.Set(key, make([]byte, 1*1024-len(key)))
			}
			wg.Done()
		}()
	}
	wg.Wait()

	// cache hit
	for i := 0; i < 100; i++ {
		key := strconv.Itoa(i)
		b, err := s.kv.Get(key)
		s.Require().NoError(err)
		exp := make([]byte, 1*1024-len(key))
		s.Require().Equal(exp, b)
	}
}

func TestKVSuite(t *testing.T) {
	suite.Run(t, new(KVSuite))
}
