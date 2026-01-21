package IManager

import "openspy.net/discord-gateway-bot/Message"

const (
	STATE_ESTABLISH int = 1
	STATE_HEARTBEAT int = 3
)

type IManager interface {
	SetTickerInterval(interval int)
	GetTickerInterval() int
	WriteMessage(event Message.Event)
	SetState(state int)
	GetDispatchChannel() chan Message.Event
}
