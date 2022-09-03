package service

import (
	"fmt"
	"os"

	"github.com/gin-contrib/sessions/redis"
)

type RedisService struct {
	redisStore *redis.Store
}

var redisService *RedisService = nil

func Redis_New() *RedisService {
	if redisService == nil {
		redisStore, err := redis.NewStore(10, "tcp", fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PASSWORD")), os.Getenv("REDIS_PASSWORD"), []byte("secret"))
		if err != nil {
			panic(err.Error())
		}
		redisService = &RedisService{
			redisStore: &redisStore,
		}
	}

	return redisService
}

func (r *RedisService) GetStore() *redis.Store {
	return r.redisStore
}

func (r *RedisService) Set(key string, value string) error {
	err, store := redis.GetRedisStore(*r.redisStore)
	if err != nil {
		return err
	}

	conn := store.Pool.Get()
	if conn.Err() != nil {
		return conn.Err()
	}

	defer conn.Close()
	_, err = conn.Do("SET", key, value)
	return err
}

func (r *RedisService) SetWithExp(key string, value string, expired int) error { //in seconds
	err, store := redis.GetRedisStore(*r.redisStore)
	if err != nil {
		return err
	}

	conn := store.Pool.Get()
	if conn.Err() != nil {
		return conn.Err()
	}

	defer conn.Close()
	_, err = conn.Do("SET", key, value, "EX", expired)
	return err
}

func (r *RedisService) Get(key string) (string, error) {
	err, store := redis.GetRedisStore(*r.redisStore)
	if err != nil {
		return "", err
	}

	conn := store.Pool.Get()
	if conn.Err() != nil {
		return "", conn.Err()
	}

	defer conn.Close()

	data, err := conn.Do("GET", key)
	if err != nil {
		return "", err
	}

	return data.(string), nil
}

func (r *RedisService) Del(keys ...string) error {
	err, store := redis.GetRedisStore(*r.redisStore)
	if err != nil {
		return err
	}

	conn := store.Pool.Get()
	if conn.Err() != nil {
		return conn.Err()
	}

	defer conn.Close()

	_, err = conn.Do("DEL", keys)
	if err != nil {
		return err
	}
	return nil
}
