package cache

import (
	"log"

	"github.com/gomodule/redigo/redis"
)

var (
	pool *redis.Pool
)

func init() {
	pool = newPool()
}

func newPool() *redis.Pool {
	return &redis.Pool{
		MaxIdle:   80,
		MaxActive: 12000,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", ":6379")
			if err != nil {
				log.Printf("nasapic: failed to dial redis: %v\n", err)
			}
			return c, err
		},
	}
}

func ping(c redis.Conn) error {
	_, err := redis.String(c.Do("PING"))
	if err != nil {
		return err
	}

	return nil
}
