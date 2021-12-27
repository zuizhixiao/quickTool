package initialize

import (
	"TEMPLATE/config"
	"fmt"
	"github.com/go-redis/redis"
)

func Redis() {
	conf := config.GVA_CONFIG.Redis
	fmt.Println("redis")
	client := redis.NewClient(&redis.Options{
		Addr: conf.Addr,
		Password: conf.Password,
		DB: conf.DB,
	})
	pong, err := client.Ping().Result()
	if err != nil {
		fmt.Println("redis连接失败")
	} else{
		fmt.Println("redis connect result is:", pong)
		config.GVA_REDIS = client
	}
}
