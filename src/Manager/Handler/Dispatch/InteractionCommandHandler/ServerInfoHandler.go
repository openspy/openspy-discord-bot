package InteractionCommandHandler

import (
	"context"
	"strings"

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
	cb.Type = 4
	cb.Data = &DiscordApi.InteractionCallbackData{}
	cb.Data.Embeds = make([]DiscordApi.Embed, 0)

	var data = interaction.Data.(*DispatchEvents.InteractionDataApplicationCommand)
	var serverId = data.Options[0].Value //validate, error out if no such game exists, etc
	if exists, _ := redisClient.HExists(ctx, serverId, "id").Result(); !exists {
		h.SendNoSuchServer(interaction)
		return
	}

	deleted, _ := redisClient.HGet(ctx, serverId, "deleted").Result()
	if len(deleted) > 0 && deleted != "0" {
		h.SendNoSuchServer(interaction)
		return
	}
	var embed *DiscordApi.Embed = &DiscordApi.Embed{}
	var col int = 0x0000ff
	embed.Color = &col
	embed.Title = "Server Info " + serverId
	embed.Description = "Server Info for: " + serverId

	keys, _ := redisClient.HGetAll(ctx, serverId+"custkeys").Result()

	var True bool = true
	var isFieldClean bool = true

	var total_length int = 0

	const MAX_TOTAL_LEN int = 5500

	basicInfo, _ := redisClient.HMGet(ctx, serverId, "allow_unsolicited_udp", "country").Result()

	var portsForwarded bool = basicInfo[0] != nil && basicInfo[0].(string) == "1"
	var country = h.simpleStringConvert(basicInfo[1])

	var field *DiscordApi.EmbedField = &DiscordApi.EmbedField{}
	field.Inline = &True
	field.Name = "Server ID"
	field.Value = serverId
	total_length = total_length + len(field.Name) + len(field.Value)
	embed.Fields = append(embed.Fields, *field)

	field.Inline = &True
	field.Name = "Ports Forwarded"
	if portsForwarded {
		field.Value = "True"
	} else {
		field.Value = "False"
	}
	total_length = total_length + len(field.Name) + len(field.Value)
	embed.Fields = append(embed.Fields, *field)

	field.Inline = &True
	field.Name = "Country"
	if country != "Unknown" {
		field.Value = country + " - :flag_" + strings.ToLower(country) + ":"
	} else {
		field.Value = "Unknown"
	}
	total_length = total_length + len(field.Name) + len(field.Value)
	embed.Fields = append(embed.Fields, *field)

	for key, value := range keys {
		if h.isSkipKey(key) {
			continue
		}
		field = &DiscordApi.EmbedField{}
		field.Inline = &True
		field.Name = key
		field.Value = value
		total_length = total_length + len(field.Name) + len(field.Value)
		embed.Fields = append(embed.Fields, *field)
		isFieldClean = false

		if len(embed.Fields) >= 24 {
			isFieldClean = true
			total_length = total_length + len(field.Name) + len(field.Value)
			total_length = total_length + len(embed.Title) + len(embed.Description) + len(embed.Url)
			cb.Data.Embeds = append(cb.Data.Embeds, *embed)
			embed.Fields = make([]DiscordApi.EmbedField, 0)
		}
		if len(cb.Data.Embeds) >= 10 || total_length > MAX_TOTAL_LEN {
			break
		}
	}
	if !isFieldClean {
		cb.Data.Embeds = append(cb.Data.Embeds, *embed)
	}
	DiscordApi.SendInteractionCallback(interaction.Id, interaction.Token, cb)
}
func (h *ServerInfoHandler) SendNoSuchServer(interaction *DispatchEvents.InteractionCreateEvent) {
	var cb DiscordApi.InteractionCallback
	cb.Data = &DiscordApi.InteractionCallbackData{}
	cb.Data.Content = "Server id not found"
	cb.Type = 4
	DiscordApi.SendInteractionCallback(interaction.Id, interaction.Token, cb)
}

func (h *ServerInfoHandler) simpleStringConvert(obj interface{}) string {
	if obj != nil {
		var s string = obj.(string)
		if len(s) > 0 {
			return s
		}
	}
	return "Unknown"
}

func (h *ServerInfoHandler) isSkipKey(key string) bool {
	switch key {
	case "publicip":
		fallthrough
	case "publicport":
		fallthrough
	case "statechanged":
		fallthrough
	case "hostport":
		fallthrough
	case "localport":
		fallthrough
	case "localip0":
		fallthrough
	case "localip1":
		fallthrough
	case "localip2":
		fallthrough
	case "localip3":
		fallthrough
	case "localip4":
		fallthrough
	case "localip5":
		fallthrough
	case "localip6":
		fallthrough
	case "localip7":
		fallthrough
	case "queryid":
		fallthrough
	case "final":
		return true
	default:
		return false
	}
}
