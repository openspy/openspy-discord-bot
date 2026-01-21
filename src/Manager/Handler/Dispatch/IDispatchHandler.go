package Dispatch

import (
	"context"

	"github.com/redis/go-redis/v9"
	"openspy.net/discord-gateway-bot/Message"
)

type IDispatch interface {
	Handle(ctx context.Context, redisClient *redis.Client, msg Message.Event)
}
