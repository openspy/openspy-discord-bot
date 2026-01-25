package Handler

import (
	"log"

	"openspy.net/discord-gateway-bot/DiscordApi"
	"openspy.net/discord-gateway-bot/IManager"
	"openspy.net/discord-gateway-bot/Message"
)

type HelloStateHandler struct {
	manager IManager.IManager
}

func (h *HelloStateHandler) SetManager(manager IManager.IManager) {
	h.manager = manager
}
func (h *HelloStateHandler) HandleMessage(msg Message.Event) {
	if msg.Op != Message.HELLO_EVENT {
		log.Panicf("Expected hello event!") //probably we should do something better than panic
	}
	var data *Message.HelloEvent = msg.Data.(*Message.HelloEvent)
	DiscordApi.HeartbeatInterval = data.HeartbeatInterval

	h.manager.SetState(IManager.STATE_HEARTBEAT)
}
func (h *HelloStateHandler) Tick() {
}
