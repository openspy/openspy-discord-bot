package DiscordApi

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
)

const (
	INTERACTION_CALLBACK_TYPE_PONG                                 int = 1
	INTERACTION_CALLBACK_TYPE_CHANNEL_MESSAGE_WITH_SOURCE              = 4
	INTERACTION_CALLBACK_TYPE_DEFERRED_CHANNEL_MESSAGE_WITH_SOURCE     = 5
	INTERACTION_CALLBACK_TYPE_DEFERRED_UPDATE_MESSAGE                  = 6
	INTERACTION_CALLBACK_TYPE_UPDATE_MESSAGE                           = 7
	INTERACTION_CALLBACK_TYPE_MODAL                                    = 9
	INTERACTION_CALLBACK_TYPE_PREMIUM_REQUIRED                         = 10
	INTERACTION_CALLBACK_TYPE_LAUNCH_ACTIVITY                          = 11
)

type InteractionCallbackData struct {
	Content string  `json:"content,omitempty"`
	Embeds  []Embed `json:"embeds,omitempty"`
}
type InteractionCallback struct {
	Type int                      `json:"type"`
	Data *InteractionCallbackData `json:"data,omitempty"`
}

func SendInteractionCallback(interactionId string, interactionToken string, msg InteractionCallback) {
	var url = API_BASE_URL + "/interactions/" + interactionId + "/" + interactionToken + "/callback"
	jsonBytes, _ := json.Marshal(&msg)
	requestBody := bytes.NewBuffer(jsonBytes)
	req, _ := http.NewRequest("POST", url, requestBody)
	req.Header.Add("Content-Type", "application/json")

	client := GetHttpClient()

	client.Do(req)
}
func EditOriginalInteractionResponse(interactionId string, interactionToken string, msg InteractionCallback) {
	var appId = os.Getenv("CLIENT_ID")
	var url = API_BASE_URL + "/webhooks/" + appId + "/" + interactionToken + "/messages/@original"
	jsonBytes, _ := json.Marshal(&msg.Data)
	requestBody := bytes.NewBuffer(jsonBytes)
	req, _ := http.NewRequest("PATCH", url, requestBody)
	req.Header.Add("Content-Type", "application/json")

	client := GetHttpClient()

	client.Do(req)
}
