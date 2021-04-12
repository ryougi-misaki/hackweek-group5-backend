package redis

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"hackathon/config"
	"time"
)

var Client *redis.Client

func Init() (err error) {
	cfg := config.Conf.RedisConfig
	Client = redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password:     cfg.Password,
		DB:           cfg.DB,
		PoolSize:     cfg.PoolSize,
		MinIdleConns: cfg.MinIdleConns,
	})

	_, err = Client.Ping().Result()
	if err != nil {
		return err
	}
	return nil
}

func Close() {
	_ = Client.Close()
}
func CacheSetData(prefix string, id string, data interface{}, expiration time.Duration) error {
	key := prefix + id
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return err
	}

	err = Client.Set(key, string(dataBytes), expiration).Err()
	if err != nil {
		return err
	}
	return nil
}

//根据id的描述prefix及id查data
//date必须传入引用，用于接收数据
func CacheGetData(prefix string, id string, data interface{}) error {
	key := prefix + id

	value, err := Client.Get(key).Result()
	if err != nil {
		return err
	}

	err = json.Unmarshal([]byte(value), data)
	if err != nil {
		return err
	}

	return nil
}

//根据id的描述prefix及id查data，若命中缓存，则按expiration更新有效期
//date必须传入引用，用于接收数据
func CacheGetDataEpr(prefix string, id string, data interface{}, expiration time.Duration) error {
	key := prefix + id

	value, err := Client.Get(key).Result()
	if err != nil {
		return err
	}

	Client.PExpire(key, expiration)

	err = json.Unmarshal([]byte(value), data)
	if err != nil {
		return err
	}

	return nil
}

func CacheDelData(prefix string, id string) error {
	key := prefix + id

	err := Client.Del(key).Err()
	if err != nil && err != redis.Nil {
		return err
	}

	return nil
}
