package redismodel

import "C"
import (
	"github.com/astaxie/beego"
	"github.com/gomodule/redigo/redis"
	"time"
)

type RedisCLient struct {
	Pool *redis.Pool
}

func GetRedisClient() RedisCLient {
	maxidle, err := beego.AppConfig.Int("redismaxidle")
	if err != nil {
		panic(err)
	}
	maxactive, err := beego.AppConfig.Int("redismaxactive")
	if err != nil {
		panic(err)
	}
	idletimeout, err := beego.AppConfig.Int64("redisidletimeout")
	if err != nil {
		panic(err)
	}
	return RedisCLient{
		Pool: &redis.Pool{
			MaxIdle:     maxidle,
			MaxActive:   maxactive,
			IdleTimeout: time.Duration(idletimeout) * time.Second,
			Dial: func() (conn redis.Conn, e error) {
				return redis.Dial("tcp", beego.AppConfig.String("redishost")+":"+beego.AppConfig.String("redisport"))
			},
		},
	}
}

func (c *RedisCLient) Put(key, value string) (interface{}, error) {
	conn := c.Pool.Get()
	defer func() {
		err := conn.Close()
		if err != nil {
			beego.Error("redis client close failure")
		}
	}()
	return conn.Do("SET", key, value)
}

func (c *RedisCLient) Get(key string) (string, error) {
	conn := c.Pool.Get()
	defer func() {
		err := conn.Close()
		if err != nil {
			beego.Error("redis client close failure")
		}
	}()
	return redis.String(conn.Do("GET", key))
}

func (c *RedisCLient) GetConn() redis.Conn {
	return c.Pool.Get()
}
