package cache

import (
	"encoding/json"
	"fmt"
	"github.com/bsm/redislock"
	"github.com/go-redis/redis"
	"go.uber.org/zap"
	"time"
)

type RedisCache struct {
	host string
	port int

	Client       *redis.Client
	ClientLocker *redislock.Client
	log          *zap.SugaredLogger
}

func NewRedisCache(logger *zap.SugaredLogger, host string, port int) Cache {
	return &RedisCache{host: host, port: port, log: logger}
}

func (r *RedisCache) Connect() error {
	r.log.Debug(fmt.Sprintf("%v:%v", r.host, r.port))
	r.Client = redis.NewClient(&redis.Options{
		Addr:               fmt.Sprintf("%v:%v", r.host, r.port),
		Password:           "",
		DB:                 0,
		PoolSize:           10,
		IdleTimeout:        time.Second * 10,
		IdleCheckFrequency: time.Second * 5,
	})
	_, err := r.Client.Ping().Result()
	r.ClientLocker = redislock.New(r.Client)
	go r.DeletePreviousOCFSData(r.Client)
	if err != nil {
		r.log.Errorf("Unable to connect to redis %v", err)
		return err
	}
	return nil
}

func (r *RedisCache) DeletePreviousOCFSData(Client *redis.Client) {
	iter := Client.Scan(0, "OCFS_*", 0).Iterator()
	for iter.Next() {
		err := Client.Del(iter.Val()).Err()
		if err != nil {
			panic(err)
		}
	}
	if err := iter.Err(); err != nil {
		r.log.Error("Error in Delete Previous OCFS Data : %v", err)
	}
}

func (r *RedisCache) ObtainLock(key string, duration time.Duration) (*redislock.Lock, string) {
	lock, err := r.ClientLocker.Obtain(key, duration, nil)
	var errString = ""
	if err == redislock.ErrNotObtained {
		r.log.Error("Could not obtain lock!, try again")
		errString = "Could not obtain lock!, try again"
	} else if err != nil {
		r.log.Error(err)
		errString = "redis lock cannot obtain"
	}
	return lock, errString
}

func (r *RedisCache) ReleaseLock(lock *redislock.Lock) {
	if lock != nil {
		lock.Release()
	} else {
		r.log.Error("lock not found to release")
	}
}

func (r *RedisCache) Disconnect() error {
	return r.Client.Close()
}

func (r *RedisCache) SAdd(key string, members ...interface{}) error {
	cmd := r.Client.SAdd(key, members...)
	return cmd.Err()
}

func (r *RedisCache) SMember(key string) []string {
	return r.Client.SMembers(key).Val()
}

func (r *RedisCache) Get(key string) string {
	return r.Client.Get(key).Val()
}

func (r *RedisCache) HGet(hash, key string) string {
	return r.Client.HGet(hash, key).Val()
}

func (r *RedisCache) Set(key string, value interface{}, duration time.Duration) (string, error) {
	marshValue, err := json.Marshal(value)
	if err != nil {
		r.log.Errorf("couldn't update in Redis %v , value : %v", err, value)
		return "", err
	}
	cmd := r.Client.Set(key, marshValue, duration)
	return cmd.Result()
}

func (r *RedisCache) HSet(hash, key string, value string) {
	r.Client.HSet(hash, key, value)
}

func (r *RedisCache) IncrByFloat(key string, value float64) float64 {
	return r.Client.IncrByFloat(key, value).Val()
}
func (r *RedisCache) IncrByInt(key string, value int64) int64 {
	return r.Client.IncrBy(key, value).Val()
}

func (r *RedisCache) Delete(key string) {
	r.Client.Del(key)
}

func (r *RedisCache) Update(key, value string, duration time.Duration) {
	r.Set(key, value, duration)
}

func (r *RedisCache) Exists(key string) int {
	return int(r.Client.Exists(key).Val())
}

func (r *RedisCache) Keys(pattern string) []string {
	return r.Client.Keys(pattern).Val()
}
