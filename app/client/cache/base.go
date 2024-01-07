package cache

import (
	"time"

	"github.com/bsm/redislock"
)

type Cache interface {
	ObtainLock(key string, duration time.Duration) (*redislock.Lock, string)
	ReleaseLock(lock *redislock.Lock)
	Connect() error
	Disconnect() error
	Get(key string) string
	HGet(hash, key string) string
	Set(key string, value interface{}, duration time.Duration) (string, error)
	HSet(hash, key, value string)
	IncrByFloat(key string, value float64) float64
	IncrByInt(key string, value int64) int64
	Delete(key string)
	Update(key, value string, duration time.Duration)
	SAdd(key string, members ...interface{}) error
	SMember(key string) []string
	Exists(key string) int
	Keys(pattern string) []string
}
