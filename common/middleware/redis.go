package middleware

import (
	"fmt"
	"github.com/go-redis/redis"
	"github.com/project-template/common/config"
	"os"
)

const RedisDefaultValue = "NONE"

func OpenRedisConnect(config *config.Config) {
	client := redis.NewClient(&redis.Options{
		DB:       config.RedisConfig.DB,
		Addr:     fmt.Sprintf("%s:%d", config.RedisConfig.Host, config.RedisConfig.Port),
		Password: config.RedisConfig.Password,
	})
	_, err := client.Ping().Result()
	if err != nil {
		fmt.Printf("初始化redis链接---->>>>> ping err:%v\n ", err)
		os.Exit(-1)
	}
	config.RedisClient = client
}

func CloseRedisConnect() {
	err := config.Info().RedisClient.Close()
	if err != nil {
		fmt.Printf("关闭redis链接---->>>>> err:%v\n ", err)
	}
}
