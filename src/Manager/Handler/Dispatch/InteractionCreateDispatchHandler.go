package Dispatch

import (
	"context"

	"github.com/redis/go-redis/v9"
	"openspy.net/discord-gateway-bot/Manager/Handler/Dispatch/InteractionCommandHandler"
	"openspy.net/discord-gateway-bot/Message"
	"openspy.net/discord-gateway-bot/Message/DispatchEvents"
)

type InteractionCreateDispatchHandler struct {
}

func (h *InteractionCreateDispatchHandler) Handle(ctx context.Context, redisClient *redis.Client, msg Message.Event) {
	var interaction = msg.Data.(*DispatchEvents.InteractionCreateEvent)
	var data = interaction.Data.(*DispatchEvents.InteractionDataApplicationCommand)

	var handler InteractionCommandHandler.ICommandHandler = nil
	if data.Name == "serverinfo" {
		handler = &InteractionCommandHandler.ServerInfoHandler{}
	} else if data.Name == "listservers" {
		handler = &InteractionCommandHandler.ListServersHandler{}
	}

	if handler != nil {
		handler.HandleCommand(ctx, redisClient, interaction)
	}

}
