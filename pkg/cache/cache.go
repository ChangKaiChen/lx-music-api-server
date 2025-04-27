package cache

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/ChangKaiChen/lx-music-api-server/global"
	"github.com/ChangKaiChen/lx-music-api-server/pkg/logger"
	"github.com/redis/go-redis/v9"
	"reflect"
	"sync"
	"time"
)

type RedisCache struct {
	rdb *redis.Client
	ctx context.Context
}

var (
	once sync.Once
	rdb  *RedisCache
)

func Init() {
	once.Do(func() {
		ctx := context.Background()
		r := redis.NewClient(&redis.Options{
			Addr:     global.GetConf().Redis.Addr,
			Password: global.GetConf().Redis.Password,
			DB:       1,
		})
		rdb = &RedisCache{rdb: r, ctx: ctx}
	})
}
func GetCache() *RedisCache {
	if rdb == nil {
		Init()
		if rdb == nil {
			panic("cache not initialized")
		}
	}
	return rdb
}
func (c *RedisCache) Set(key string, value any, expiration time.Duration) {
	log := logger.GetLogger()
	data, err := json.Marshal(value)
	if err != nil {
		log.Errorf("", "error marshaling value: %v", err.Error())
		return
	}
	err = c.rdb.Set(c.ctx, key, data, expiration).Err()
	if err != nil {
		log.Errorf("", "error setting value in Redis: %v", err)
		return
	}
}
func (c *RedisCache) Get(key string, result any) error {
	resultType := reflect.TypeOf(result)
	if resultType.Kind() != reflect.Ptr {
		return errors.New("result must be a pointer")
	}
	data, err := c.rdb.Get(c.ctx, key).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return errors.New("key does not exist")
		}
		return err
	}
	if err = json.Unmarshal(data, result); err != nil {
		return err
	}
	return nil
}
