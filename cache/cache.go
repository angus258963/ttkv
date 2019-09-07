package cache

import (
	"fmt"
)

var (
	ErrNotFound     = fmt.Errorf("cache not found")
	ErrMaxValueSize = fmt.Errorf("max value size")
	ErrNoMoreCap    = fmt.Errorf("no more cap")
)

type StrategyType int

const (
	FIFO StrategyType = iota
)

type Cache interface {
	Get(key string) ([]byte, error)
	Set(key string, val []byte) error
}
