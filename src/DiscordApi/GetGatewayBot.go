package DiscordApi

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
)

type GetGatewayBotResponse struct {
	Url string `json:"url"`
}

func GetGatewayBot() GetGatewayBotResponse {
	var botToken = os.Getenv("BOT_TOKEN")
	var gatewayUrl = "https://discord.com/api/v10/gateway/bot"
	req, _ := http.NewRequest("GET", gatewayUrl, nil)
	req.Header.Set("Authorization", "Bot "+botToken)
	client := GetHttpClient()

	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	respData, _ := io.ReadAll(res.Body)

	var resp GetGatewayBotResponse
	json.Unmarshal(respData, &resp)

	return resp
}
