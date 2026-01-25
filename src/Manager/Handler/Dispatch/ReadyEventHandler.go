package Dispatch

import (
	"context"

	"github.com/redis/go-redis/v9"
	"openspy.net/discord-gateway-bot/DiscordApi"
	"openspy.net/discord-gateway-bot/Message"
	"openspy.net/discord-gateway-bot/Message/DispatchEvents"
)

type ReadyEventDispatchHandler struct {
}

func (h *ReadyEventDispatchHandler) Handle(ctx context.Context, redisClient *redis.Client, msg Message.Event) {
	var event = msg.Data.(*DispatchEvents.ReadyEvent)
	DiscordApi.ResumeGatewayUrl = event.ResumeGatewayUrl
	DiscordApi.SessionId = event.SessionId
}
