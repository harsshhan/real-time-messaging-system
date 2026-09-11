package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type RedisService struct {
	Client *redis.Client
}

func NewRedisService(client *redis.Client) *RedisService {
	return &RedisService{
		Client: client,
	}
}

func (r *RedisService) Publish(ctx context.Context, channel string, message string) error {
	return r.Client.Publish(ctx, channel, message).Err()
}

func (r *RedisService) Subscribe(ctx context.Context, channel string) *redis.PubSub {

	pubsub := r.Client.Subscribe(ctx, channel)

	return pubsub
}


func (r *RedisService) SetOnline(ctx context.Context, userID uuid.UUID) error {
	return r.Client.Set(
		ctx,
		"user:"+userID.String(),
		"online",
		0,
	).Err()
}

func (r *RedisService) SetOffline(ctx context.Context, userID uuid.UUID) error {
	return r.Client.Del(
		ctx,
		"user:"+userID.String(),
	).Err()
}