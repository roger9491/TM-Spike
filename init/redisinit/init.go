package redisinit

import (
	"log"
	"strconv"

	"github.com/redis/go-redis/v9"
)

func InitRedis(ip, port, pass, db string) *redis.Client {

	dbint, err := strconv.Atoi(db)
	if err != nil {
		log.Fatal(err)
	}
	addr := ip + ":" + port
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pass,
		DB:       dbint,
	})

	return rdb
}
