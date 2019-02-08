package redis

import (
	. "github.com/Fallensouls/Pandora/setting"
	"github.com/go-redis/redis"
	"log"
)

var client *redis.Client

func init() {
	client = redis.NewClient(&redis.Options{
		Addr:     Config.RedisHost + ":" + Config.RedisPort,
		Password: Config.RedisPassword,
		DB:       0,
	})

	_, err := client.Ping().Result()
	if err != nil {
		log.Panicln("failed to connect to Redis:" + err.Error())
	}
}
