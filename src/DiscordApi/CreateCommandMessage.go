package DiscordApi

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type CreateCommandMessage struct {
	Name                     string                `json:"name"`
	Type                     int                   `json:"type"`
	Description              string                `json:"description"`
	DefaultMemberPermissions int                   `json:"default_member_permissions"`
	Options                  []CreateCommandOption `json:"options"`
}

func (msg *CreateCommandMessage) Send(applicationId string, apiKey string) {
	var oauthUrl = API_BASE_URL + "/applications/" + applicationId + "/commands"

	jsnoBytes, _ := json.Marshal(&msg)
	requestBody := bytes.NewBuffer(jsnoBytes)
	req, _ := http.NewRequest("POST", oauthUrl, requestBody)
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Add("Content-Type", "application/json")

	client := GetHttpClient()

	client.Do(req)
}
