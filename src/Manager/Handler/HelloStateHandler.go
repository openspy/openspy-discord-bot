package Handler

import (
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
	var data *Message.HelloEvent = msg.Data.(*Message.HelloEvent)
	h.manager.SetTickerInterval(data.HeartbeatInterval)
	h.manager.SetState(IManager.STATE_HEARTBEAT)
}
