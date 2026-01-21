package InteractionCommandHandler

import (
	"context"

	"github.com/redis/go-redis/v9"
	"openspy.net/discord-gateway-bot/Message/DispatchEvents"
)

type ICommandHandler interface {
	HandleCommand(ctx context.Context, redisClient *redis.Client, interaction *DispatchEvents.InteractionCreateEvent)
}
