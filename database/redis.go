package database

import (
	"context"
	"github.com/redis/go-redis/v9"
)

func ConnectRedis() (*redis.Client, error){
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	err := client.Ping(context.Background()).Err()

	if err!= nil{
		return nil, err
	}

	return client, nil
}