package InteractionCommandHandler

import (
	"context"
	"log"
	"strconv"
	"strings"

	"github.com/redis/go-redis/v9"
	"openspy.net/discord-gateway-bot/DiscordApi"
	"openspy.net/discord-gateway-bot/Message/DispatchEvents"
)

type ListServersHandler struct {
}

func (h *ListServersHandler) HandleCommand(ctx context.Context, redisClient *redis.Client, interaction *DispatchEvents.InteractionCreateEvent) {

	var data = interaction.Data.(*DispatchEvents.InteractionDataApplicationCommand)
	var gamename = data.Options[0].Value //validate, error out if no such game exists, etc
	if exists, _ := redisClient.Exists(ctx, gamename).Result(); exists == 0 {
		h.SendNoSuchGamename(interaction)
		return
	}

	var cb DiscordApi.InteractionCallback
	cb.Type = 4
	cb.Data = &DiscordApi.InteractionCallbackData{}
	cb.Data.Embeds = make([]DiscordApi.Embed, 0)

	redisCmds := []*redis.SliceCmd{}

	fetchPipeline := redisClient.Pipeline()
	var cursor uint64 = 0
	const SERVER_BATCH_COUNT int = 4
	for {
		keys, nextCursor, err := redisClient.ZScan(ctx, gamename, cursor, "*", int64(SERVER_BATCH_COUNT)).Result()
		if err != nil { //XXX: handle and send error back
			log.Panic(err)
		}

		for _, key := range keys {
			redisCmds = append(redisCmds, fetchPipeline.HMGet(ctx, key, "id", "country"))
			redisCmds = append(redisCmds, fetchPipeline.HMGet(ctx, key+"custkeys", "hostname", "numplayers", "maxplayers", "mapname", "gamemode"))

		}
		cursor = nextCursor
		if cursor == 0 {
			break
		}
	}

	var embed *DiscordApi.Embed = &DiscordApi.Embed{}
	var col int = 0x0000ff
	embed.Color = &col
	var page int = 1
	embed.Title = "Server List - Page " + strconv.Itoa(page)
	embed.Description = "Server List - " + gamename + " - Page " + strconv.Itoa(page)
	embed.Url = "https://openspy.net/server-list/" + gamename
	embed.Footer = &DiscordApi.EmbedFooter{}
	embed.Footer.Text = "Use /serverinfo [server key] for detailed server information"

	_, redisErr := fetchPipeline.Exec(ctx)
	if redisErr != nil {
		log.Panic(redisErr)
	}
	var True bool = true
	var isFieldClean bool = true

	var total_length int = 0

	for i := 0; i < len(redisCmds); i += 2 {
		keys, _ := redisCmds[i].Result()
		custkeys, _ := redisCmds[i+1].Result()

		if total_length > DiscordApi.MAX_TOTAL_LEN {
			break
		}
		if keys[0] == nil {
			continue
		}

		var backendId string = h.simpleStringConvert(keys[0])
		var country string = h.simpleStringConvert(keys[1])

		var hostname string = h.simpleStringConvert(custkeys[0])
		var numplayers string = h.simpleStringConvert(custkeys[1])
		var maxplayers string = h.simpleStringConvert(custkeys[2])
		var mapname string = h.simpleStringConvert(custkeys[3])
		var gamemode string = h.simpleStringConvert(custkeys[4])

		isFieldClean = false
		var field *DiscordApi.EmbedField = &DiscordApi.EmbedField{}
		field.Inline = &True
		field.Name = "Server Key"
		field.Value = gamename + ":" + backendId + ":"
		if strings.Contains(backendId, "ut2004") { //temporarily skip injected servers until key bug is fixed
			field.Value = backendId
		}
		total_length = total_length + len(field.Name) + len(field.Value)
		embed.Fields = append(embed.Fields, *field)

		field.Inline = &True
		field.Name = "Hostname"
		field.Value = hostname
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

		field.Inline = &True
		field.Name = "Players"
		field.Value = numplayers + "/" + maxplayers
		total_length = total_length + len(field.Name) + len(field.Value)
		embed.Fields = append(embed.Fields, *field)

		field.Inline = &True
		field.Name = "Mapname"
		field.Value = mapname
		total_length = total_length + len(field.Name) + len(field.Value)
		embed.Fields = append(embed.Fields, *field)

		field.Inline = &True
		field.Name = "Gamemode"
		field.Value = gamemode
		total_length = total_length + len(field.Name) + len(field.Value)
		embed.Fields = append(embed.Fields, *field)
		if len(embed.Fields) >= 24 {
			isFieldClean = true
			total_length = total_length + len(field.Name) + len(field.Value)
			total_length = total_length + len(embed.Title) + len(embed.Description) + len(embed.Url)
			cb.Data.Embeds = append(cb.Data.Embeds, *embed)
			embed.Fields = make([]DiscordApi.EmbedField, 0)

			page = page + 1
			embed.Url = ""
			embed.Title = "Server List - Page " + strconv.Itoa(page)
			embed.Description = "Server List - " + gamename + " - Page " + strconv.Itoa(page)
		}
		if len(cb.Data.Embeds) >= 10 || total_length > DiscordApi.MAX_TOTAL_LEN {
			break
		}
	}

	if !isFieldClean {
		cb.Data.Embeds = append(cb.Data.Embeds, *embed)
	}
	DiscordApi.EditOriginalInteractionResponse(interaction.Id, interaction.Token, cb)

}

func (h *ListServersHandler) SendNoSuchGamename(interaction *DispatchEvents.InteractionCreateEvent) {
	var cb DiscordApi.InteractionCallback
	cb.Data = &DiscordApi.InteractionCallbackData{}
	cb.Data.Content = "No servers found, or game doesn't exist."
	cb.Type = 4
	DiscordApi.EditOriginalInteractionResponse(interaction.Id, interaction.Token, cb)
}

func (h *ListServersHandler) simpleStringConvert(obj interface{}) string {
	if obj != nil {
		var s string = obj.(string)
		if len(s) > 0 {
			return s
		}
	}
	return "Unknown"
}
