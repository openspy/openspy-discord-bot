package Dispatch

import (
	"context"

	"github.com/redis/go-redis/v9"
	"openspy.net/discord-gateway-bot/DiscordApi"
	"openspy.net/discord-gateway-bot/Manager/Handler/Dispatch/InteractionCommandHandler"
	"openspy.net/discord-gateway-bot/Message"
	"openspy.net/discord-gateway-bot/Message/DispatchEvents"
)

type InteractionCreateDispatchHandler struct {
}

func (h *InteractionCreateDispatchHandler) Handle(ctx context.Context, redisClient *redis.Client, msg Message.Event) {
	var interaction = msg.Data.(*DispatchEvents.InteractionCreateEvent)
	var data = interaction.Data.(*DispatchEvents.InteractionDataApplicationCommand)

	h.sendDeferredMessage(interaction)

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

func (h *InteractionCreateDispatchHandler) sendDeferredMessage(interaction *DispatchEvents.InteractionCreateEvent) {
	var cb DiscordApi.InteractionCallback
	cb.Type = 5 //DEFERRED_CHANNEL_MESSAGE_WITH_SOURCE
	DiscordApi.SendInteractionCallback(interaction.Id, interaction.Token, cb)
}
