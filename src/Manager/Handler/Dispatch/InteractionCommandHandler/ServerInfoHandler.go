package InteractionCommandHandler

import (
	"context"
	"log"
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
	embed.Title = "Server Info"
	embed.Description = "Server Info for: " + serverId

	fetchPipeline := redisClient.Pipeline()

	getAllReq := fetchPipeline.HGetAll(ctx, serverId+"custkeys")
	basicInfoReq := fetchPipeline.HMGet(ctx, serverId, "allow_unsolicited_udp", "country")
	_, redisErr := fetchPipeline.Exec(ctx)
	if redisErr != nil {
		h.SendError(interaction)
		log.Println(redisErr.Error())
		return
	}

	keys, _ := getAllReq.Result()

	var True bool = true
	var isFieldClean bool = true

	var total_length int = 0

	basicInfo, _ := basicInfoReq.Result()

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
		if len(cb.Data.Embeds) >= 10 || total_length > DiscordApi.MAX_TOTAL_LEN {
			break
		}
	}
	if !isFieldClean {
		cb.Data.Embeds = append(cb.Data.Embeds, *embed)
	}

	total_length = h.appendTeamKeys(ctx, redisClient, serverId, &cb, total_length)
	total_length = h.appendPlayerKeys(ctx, redisClient, serverId, &cb, total_length)

	DiscordApi.EditOriginalInteractionResponse(interaction.Id, interaction.Token, cb)
}
func (h *ServerInfoHandler) appendPlayerKeys(ctx context.Context, redisClient *redis.Client, serverId string, cb *DiscordApi.InteractionCallback, total_length int) int {
	var needExit bool = false
	var isFieldClean bool = true

	var embed *DiscordApi.Embed = &DiscordApi.Embed{}
	var col int = 0x0000ff
	embed.Color = &col
	embed.Title = "Server Info - Players"
	embed.Description = "Player Values"

	var cursor uint64 = 0

	var True bool = true
	var False bool = false

	for {
		playerKeys, nextCursor, _ := redisClient.Scan(ctx, cursor, serverId+"custkeys_player_*", 50).Result()
		const TEAM_STRING_OFFSET int = 16 //strlen of custkeys_player_
		cursor = nextCursor

		fetchPipeline := redisClient.Pipeline()
		redisCmds := []*redis.MapStringStringCmd{}
		for _, playerKey := range playerKeys {
			redisCmds = append(redisCmds, fetchPipeline.HGetAll(ctx, playerKey))
		}
		_, fetchErr := fetchPipeline.Exec(ctx)
		if fetchErr != nil {
			log.Println("player key error: " + fetchErr.Error())
			break
		}
		for idx, playerKey := range playerKeys {
			var indexString = playerKey[len(serverId)+TEAM_STRING_OFFSET:]
			var field *DiscordApi.EmbedField = &DiscordApi.EmbedField{}
			field.Inline = &False

			field.Name = "Player Index"
			field.Value = indexString
			total_length = total_length + len(field.Name) + len(field.Value)
			embed.Fields = append(embed.Fields, *field)

			allKeys, allKeyErr := redisCmds[idx].Result()
			if allKeyErr != nil {
				log.Println("team key error: " + allKeyErr.Error())
				break
			}
			for key, value := range allKeys {
				var field *DiscordApi.EmbedField = &DiscordApi.EmbedField{}
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
				if len(cb.Data.Embeds) >= 10 || total_length > DiscordApi.MAX_TOTAL_LEN {
					needExit = true
					break
				}
			}
			if needExit {
				break
			}
		}
		if cursor == 0 || needExit {
			break
		}
	}

	if !isFieldClean {
		cb.Data.Embeds = append(cb.Data.Embeds, *embed)
		isFieldClean = true
	}

	return total_length
}
func (h *ServerInfoHandler) appendTeamKeys(ctx context.Context, redisClient *redis.Client, serverId string, cb *DiscordApi.InteractionCallback, total_length int) int {
	var needExit bool = false
	var isFieldClean bool = true

	var embed *DiscordApi.Embed = &DiscordApi.Embed{}
	var col int = 0x0000ff
	embed.Color = &col
	embed.Title = "Server Info - Teams"
	embed.Description = "Team Values"

	var cursor uint64 = 0

	var True bool = true
	var False bool = false

	for {
		playerKeys, nextCursor, _ := redisClient.Scan(ctx, cursor, serverId+"custkeys_team_*", 50).Result()
		const TEAM_STRING_OFFSET int = 14 //strlen of custkeys_team_
		cursor = nextCursor

		fetchPipeline := redisClient.Pipeline()
		redisCmds := []*redis.MapStringStringCmd{}
		for _, playerKey := range playerKeys {
			redisCmds = append(redisCmds, fetchPipeline.HGetAll(ctx, playerKey))
		}
		_, fetchErr := fetchPipeline.Exec(ctx)
		if fetchErr != nil {
			log.Println("team key error: " + fetchErr.Error())
			break
		}
		for idx, playerKey := range playerKeys {
			var indexString = playerKey[len(serverId)+TEAM_STRING_OFFSET:]
			var field *DiscordApi.EmbedField = &DiscordApi.EmbedField{}
			field.Inline = &False

			field.Name = "Team Index"
			field.Value = indexString
			total_length = total_length + len(field.Name) + len(field.Value)
			embed.Fields = append(embed.Fields, *field)

			allKeys, allKeyErr := redisCmds[idx].Result()
			if allKeyErr != nil {
				log.Println("team key error: " + allKeyErr.Error())
				break
			}
			for key, value := range allKeys {
				var field *DiscordApi.EmbedField = &DiscordApi.EmbedField{}
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
				if len(cb.Data.Embeds) >= 10 || total_length > DiscordApi.MAX_TOTAL_LEN {
					needExit = true
					break
				}
			}
			if needExit {
				break
			}
		}
		if cursor == 0 || needExit {
			break
		}
	}

	if !isFieldClean {
		cb.Data.Embeds = append(cb.Data.Embeds, *embed)
		isFieldClean = true
	}

	return total_length
}
func (h *ServerInfoHandler) SendNoSuchServer(interaction *DispatchEvents.InteractionCreateEvent) {
	var cb DiscordApi.InteractionCallback
	cb.Data = &DiscordApi.InteractionCallbackData{}
	cb.Data.Content = "Server id not found"
	cb.Type = 4
	DiscordApi.EditOriginalInteractionResponse(interaction.Id, interaction.Token, cb)
}

func (h *ServerInfoHandler) SendError(interaction *DispatchEvents.InteractionCreateEvent) {
	var cb DiscordApi.InteractionCallback
	cb.Data = &DiscordApi.InteractionCallbackData{}
	cb.Data.Content = "Request failed, error occured"
	cb.Type = 4
	DiscordApi.EditOriginalInteractionResponse(interaction.Id, interaction.Token, cb)
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
