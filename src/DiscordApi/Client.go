package DiscordApi

import (
	"crypto/tls"
	"net/http"
	"net/url"
	"time"
)

const MAX_TOTAL_LEN int = 5500

var SessionId string
var ResumeGatewayUrl string
var LastSeqId int
var HeartbeatInterval int

const (
	INTERACTION_TYPE_PING                             int = 1
	INTERACTION_TYPE_APPLICATION_COMMAND                  = 2
	INTERACTION_TYPE_MESSAGE_COMPONENT                    = 3
	INTERACTION_TYPE_APPLICATION_COMMAND_AUTOCOMPLETE     = 4
	INTERACTION_TYPE_MODAL_SUBMIT                         = 5
)

const (
	API_BASE_URL string = "https://discord.com/api/v10"
)

var ApiProxy *string

type DiscordRequest struct {
	Type int         `json:"type"`
	Data interface{} `json:"data"`
}

func GetHttpClient() *http.Client {
	var tr *http.Transport
	if ApiProxy != nil && len(*ApiProxy) > 0 {
		proxyUrl, _ := url.Parse(*ApiProxy)
		tr = &http.Transport{
			MaxIdleConns:       10,
			IdleConnTimeout:    30 * time.Second,
			DisableCompression: true,
			Proxy:              http.ProxyURL(proxyUrl),
			TLSClientConfig:    &tls.Config{InsecureSkipVerify: true},
		}
	} else {
		tr = &http.Transport{
			MaxIdleConns:       10,
			IdleConnTimeout:    30 * time.Second,
			DisableCompression: true,
		}
	}

	return &http.Client{Transport: tr}
}
