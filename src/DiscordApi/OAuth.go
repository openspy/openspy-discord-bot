package DiscordApi

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type OAuthResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	Scope       string `json:"scope"`
}

func GetOAuthToken() string {
	var clientId = os.Getenv("CLIENT_ID")
	var clientSecret = os.Getenv("CLIENT_SECRET")
	var oauthUrl = "https://discord.com/api/v10/oauth2/token"
	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("scope", "identify connections applications.commands.update")
	req, _ := http.NewRequest("POST", oauthUrl, strings.NewReader(data.Encode()))
	req.SetBasicAuth(clientId, clientSecret)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	client := GetHttpClient()

	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	respData, _ := io.ReadAll(res.Body)

	var auth OAuthResponse
	json.Unmarshal(respData, &auth)

	return auth.AccessToken
}
