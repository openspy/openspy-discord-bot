package Handler

import (
	"log"
	"os"

	"openspy.net/discord-gateway-bot/DiscordApi"
	"openspy.net/discord-gateway-bot/IManager"
	"openspy.net/discord-gateway-bot/Message"
	"openspy.net/discord-gateway-bot/Message/DispatchEvents"
)

type ReconnectStateHandler struct {
	manager       IManager.IManager
	reconnectSent bool
}

func (h *ReconnectStateHandler) SetManager(manager IManager.IManager) {
	h.manager = manager
}
func (h *ReconnectStateHandler) HandleMessage(msg Message.Event) {
	if msg.Op == Message.DISPATCH_EVENT {
		if msg.Type == DispatchEvents.RESUMED {
			h.manager.SetState(IManager.STATE_HEARTBEAT)
			return
		}
		h.manager.GetDispatchChannel() <- msg
		return
	}
	log.Panicf("got unexpected event: %v", msg)
}

func (h *ReconnectStateHandler) Tick() {
	if h.reconnectSent {
		return
	}
	h.sendReconnect()
}

func (h *ReconnectStateHandler) sendReconnect() {
	h.reconnectSent = true
	var msg Message.Event
	msg.Op = Message.RESUME_EVENT
	var data Message.ResumeEvent
	data.Token = os.Getenv("BOT_TOKEN")
	data.SessionId = DiscordApi.SessionId
	data.Sequence = DiscordApi.LastSeqId
	msg.Data = data

	h.manager.WriteMessage(msg)
}
