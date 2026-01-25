package main

import (
	"context"
	"crypto/tls"
	"flag"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
	"openspy.net/discord-gateway-bot/DiscordApi"
	"openspy.net/discord-gateway-bot/Manager"
	"openspy.net/discord-gateway-bot/Manager/Handler/Dispatch"
	"openspy.net/discord-gateway-bot/Message"
	"openspy.net/discord-gateway-bot/Message/DispatchEvents"
)

var proxy = flag.String("proxy", "", "http proxy address")
var installMode = flag.Bool("install", false, "Setup / Install bot")

func main() {
	signalChan := make(chan os.Signal, 1)
	flag.Parse()
	signal.Notify(signalChan, os.Interrupt)

	if *installMode == true {
		log.Println("got install mode")
		var oauthToken = DiscordApi.GetOAuthToken()
		installCommands(oauthToken)
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer func() {
		signal.Stop(signalChan)
		cancel()
	}()

	flag.Parse()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	var dialer websocket.Dialer = *websocket.DefaultDialer
	if len(*proxy) > 0 {
		proxyURL, err := url.Parse(*proxy)
		if err != nil {
			log.Panic(err)
		}

		dialer = websocket.Dialer{
			Proxy:           http.ProxyURL(proxyURL),
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			// Other fields like HandshakeTimeout, TLSClientConfig can be set here
		}
	}

	var redisClient *redis.Client = redis.NewClient(GetRedisOptions())

	dispatchChan := make(chan Message.Event, 256)

	gatewayBotInfo := DiscordApi.GetGatewayBot()

	DiscordApi.ResumeGatewayUrl = gatewayBotInfo.Url

	var sm Manager.StateManager
	sm.SetAddress(gatewayBotInfo.Url)
	sm.SetDialer(dialer)
	sm.SetDispatchChannel(dispatchChan)
	sm.SetTickerInterval(500 * time.Millisecond)
	sm.Start()

	var isRunning bool = true
	for {
		select {
		case event := <-dispatchChan:
			handleDispatch(ctx, redisClient, event)
		case <-signalChan:
			cancel()
		case <-ctx.Done():
			log.Println("Shutdown event", ctx.Err())
			isRunning = false
		}
		if !isRunning {
			break
		}
	}
	sm.Stop()
}

func GetRedisOptions() *redis.Options {
	redisOptions := &redis.Options{
		Addr: os.Getenv("REDIS_SERVER"),
	}
	redisUsername := os.Getenv("REDIS_USERNAME")
	redisPassword := os.Getenv("REDIS_PASSWORD")

	if len(redisUsername) > 0 {
		redisOptions.Username = redisUsername
	}
	if len(redisPassword) > 0 {
		redisOptions.Password = redisPassword
	}

	redisUseTLS := os.Getenv("REDIS_USE_TLS")

	if len(redisUseTLS) > 0 {
		useTLSInt, _ := strconv.Atoi(redisUseTLS)

		if useTLSInt == 1 {
			tlsConfig := &tls.Config{
				MinVersion: tls.VersionTLS12,
				//InsecureSkipVerify: true,
				//Certificates: []tls.Certificate{cert}
			}

			redisSkipSSLVerify := os.Getenv("REDIS_INSECURE_TLS")
			if len(redisSkipSSLVerify) > 0 {
				useInsecureTLS, _ := strconv.Atoi(redisSkipSSLVerify)

				if useInsecureTLS == 1 {
					tlsConfig.InsecureSkipVerify = true
				}
			}
			redisOptions.TLSConfig = tlsConfig
		}
	}

	redisOptions.DB = 0
	return redisOptions
}

func installCommands(apiKey string) {
	var appId = os.Getenv("CLIENT_ID")
	var cmd = &DiscordApi.CreateCommandMessage{}
	cmd.Name = "listservers"
	cmd.Description = "Lists game servers for a given game"
	cmd.Type = 1

	var opt = &DiscordApi.CreateCommandOption{}
	opt.Name = "gamename"
	opt.Description = "GameSpy gamename"
	opt.Type = 3
	opt.Required = true
	opt.MinLength = 2
	opt.MaxLength = 32
	cmd.Options = append(cmd.Options, *opt)

	cmd.Send(appId, apiKey)

	cmd = &DiscordApi.CreateCommandMessage{}
	cmd.Name = "serverinfo"
	cmd.Description = "Lists game servers for a given game"
	cmd.Type = 1

	opt = &DiscordApi.CreateCommandOption{}
	opt.Name = "id"
	opt.Description = "Server ID"
	opt.Type = 3
	opt.Required = true
	opt.MinLength = 2
	opt.MaxLength = 48
	cmd.Options = append(cmd.Options, *opt)

	cmd.Send(appId, apiKey)
}

func handleDispatch(ctx context.Context, client *redis.Client, msg Message.Event) {
	var dispatch Dispatch.IDispatch = nil
	switch msg.Type {
	case DispatchEvents.READY:
		dispatch = &Dispatch.ReadyEventDispatchHandler{}
	case DispatchEvents.INTERACTION_CREATE:
		dispatch = &Dispatch.InteractionCreateDispatchHandler{}
	}

	if dispatch != nil {
		dispatch.Handle(ctx, client, msg)
	}
}
