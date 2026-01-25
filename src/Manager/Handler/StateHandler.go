package Handler

import (
	"openspy.net/discord-gateway-bot/IManager"
	"openspy.net/discord-gateway-bot/Message"
)

type StateHandler interface {
	HandleMessage(msg Message.Event)
	SetManager(manager IManager.IManager)
	Tick()
}
