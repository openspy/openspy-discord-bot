package InteractionCommandHandler

import (
	"context"

	"github.com/redis/go-redis/v9"
	"openspy.net/discord-gateway-bot/DiscordApi"
	"openspy.net/discord-gateway-bot/Message/DispatchEvents"
)

type ServerInfoHandler struct {
}

func (h *ServerInfoHandler) HandleCommand(ctx context.Context, redisClient *redis.Client, interaction *DispatchEvents.InteractionCreateEvent) {
	//handle interaction create event
	var cb DiscordApi.InteractionCallback
	cb.Type = 4
	cb.Data = &DiscordApi.InteractionCallbackData{}
	cb.Data.Content = "server info handler"
	DiscordApi.SendInteractionCallback(interaction.Id, interaction.Token, cb)
}
