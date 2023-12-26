package redis

import (
	"github.com/gomodule/redigo/redis"
	"live_recording_tools/config"
	"time"
)

type RedisCli struct {
	Pool *redis.Pool
}

var RedisClient *RedisCli

func Default() {
	rediscli := new(RedisCli)
	rediscli.Pool = &redis.Pool{
		MaxIdle:     256,
		MaxActive:   1024,
		IdleTimeout: 60 * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.Dial(
				"tcp",
				config.Configs.RedisHost,
				redis.DialReadTimeout(time.Duration(1000)*time.Millisecond),
				redis.DialWriteTimeout(time.Duration(1000)*time.Millisecond),
				redis.DialConnectTimeout(time.Duration(1000)*time.Millisecond),
				//redis.DialDatabase(config.Configs.RedisDb),
				redis.DialPassword(config.Configs.RedisPwd),
			)
		},
	}
	RedisClient = rediscli
}

func (r *RedisCli) Set(key string, val interface{}) {
	redisConn := r.Pool.Get()
	defer redisConn.Close()
	redisConn.Do("set", key, val)
}

func (r *RedisCli) SetEx(key string, val interface{}, expire int) {
	redisConn := r.Pool.Get()
	defer redisConn.Close()
	redisConn.Do("set", key, val, "EX", expire)
}

func (r *RedisCli) Get(key string) string {
	redisConn := r.Pool.Get()
	defer redisConn.Close()
	res, err := redis.String(redisConn.Do("GET", key))
	if err != nil {
		return ""
	}
	return res
}

func (r *RedisCli) Del(key string) {
	redisConn := r.Pool.Get()
	defer redisConn.Close()
	redisConn.Do("del", key)
}
