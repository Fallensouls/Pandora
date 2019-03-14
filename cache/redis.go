package cache

import (
	. "github.com/go-pandora/core/conf"
	"github.com/go-redis/redis"
	"log"
	"net"
)

var client *redis.Client

func init() {
	client = redis.NewClient(&redis.Options{
		Addr:     net.JoinHostPort(Config.RedisHost, Config.RedisPort),
		Password: Config.RedisPassword,
		DB:       0,
	})

	_, err := client.Ping().Result()
	if err != nil {
		log.Panicln("failed to connect to Redis:" + err.Error())
	}
}
