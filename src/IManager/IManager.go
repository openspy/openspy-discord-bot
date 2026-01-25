package IManager

import "openspy.net/discord-gateway-bot/Message"

const (
	STATE_ESTABLISH = iota
	STATE_HEARTBEAT
	STATE_RECONNECT
)

type IManager interface {
	WriteMessage(event Message.Event)
	SetState(state int)
	GetDispatchChannel() chan Message.Event
}
