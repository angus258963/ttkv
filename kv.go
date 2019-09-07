package kv

const (
	MB = 1024 * 1024
)

type KV interface {
	Get(key string) ([]byte, error)
	Set(key string, value []byte) error
}
