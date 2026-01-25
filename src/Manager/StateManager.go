package Manager

import (
	"log"
	"time"

	"github.com/gorilla/websocket"
	"openspy.net/discord-gateway-bot/DiscordApi"
	"openspy.net/discord-gateway-bot/IManager"
	"openspy.net/discord-gateway-bot/Manager/Handler"
	"openspy.net/discord-gateway-bot/Message"
)

type StateManager struct {
	dialer       websocket.Dialer
	address      string
	stateHandler Handler.StateHandler
	connection   *websocket.Conn

	dispatchChannel chan Message.Event

	ticker *time.Ticker
}

func (s *StateManager) SetDispatchChannel(dispatchChannel chan Message.Event) {
	s.dispatchChannel = dispatchChannel
}
func (s *StateManager) SetAddress(address string) {
	s.address = address
}
func (s *StateManager) SetDialer(dialer websocket.Dialer) {
	s.dialer = dialer
}
func (s *StateManager) Start() {
	s.SetState(IManager.STATE_ESTABLISH)

	go s.Run()
}
func (s *StateManager) Tick() {
	for _ = range s.ticker.C {
		if s.stateHandler != nil {
			s.stateHandler.Tick()
		}
	}
}
func (s *StateManager) SetTickerInterval(duration time.Duration) {
	if s.ticker != nil {
		s.ticker.Stop()
	}
	s.ticker = time.NewTicker(duration)
	go s.Tick()
}
func (s *StateManager) SetState(state int) {
	switch state {
	case IManager.STATE_ESTABLISH:
		s.stateHandler = &Handler.HelloStateHandler{}
	case IManager.STATE_HEARTBEAT:
		s.stateHandler = &Handler.HeartbeatStateHandler{}
	case IManager.STATE_RECONNECT:
		s.stateHandler = &Handler.ReconnectStateHandler{}
	}
	s.stateHandler.SetManager(s)
}
func (s *StateManager) Stop() {

}
func (s *StateManager) Run() {
	c, _, err := s.dialer.Dial(s.address, nil)
	if err != nil {
		log.Panic(err)
	}
	s.connection = c
	for {
		_, message, err := s.connection.ReadMessage()
		if err != nil {
			s.handleReadError(err)
			return
		}

		var event Message.Event
		event.UnmarshalJSON(message)

		if event.Op == Message.RECONNECT_EVENT { //need to reconnect - resume
			s.resume()
			return
		} else if event.Op == Message.INVALID_SESSION_EVENT { //need to reconnect - maybe resume
			s.handleInvalidSession(event)
			return
		}
		s.stateHandler.HandleMessage(event)
	}
}

func (s *StateManager) WriteMessage(event Message.Event) {
	s.connection.WriteJSON(event)
}
func (s *StateManager) GetDispatchChannel() chan Message.Event {
	return s.dispatchChannel
}
func (s *StateManager) resume() {
	s.SetAddress(DiscordApi.ResumeGatewayUrl)

	s.SetState(IManager.STATE_RECONNECT)

	go s.Run()
}
func (s *StateManager) reestablish() {
	s.SetAddress(DiscordApi.ResumeGatewayUrl)
	s.SetState(IManager.STATE_ESTABLISH)
	go s.Run()
}
func (s *StateManager) handleInvalidSession(event Message.Event) {
	var reconnect bool = event.Data.(bool)
	if reconnect {
		s.resume()
	} else {
		s.reestablish()
	}
}
func (s *StateManager) handleReadError(err error) {
	if closeErr, ok := err.(*websocket.CloseError); ok {
		if s.isReconnectError(closeErr.Code) {
			s.reestablish()
		} else {
			log.Panic("Got discord error which cannot reconnect")
		}
	}
}
func (s *StateManager) isReconnectError(code int) bool {
	if code < 4000 {
		return true
	}
	switch code {
	case 4000:
		fallthrough
	case 4001:
		fallthrough
	case 4002:
		fallthrough
	case 4003:
		fallthrough
	case 4005:
		fallthrough
	case 4007:
		fallthrough
	case 4008:
		fallthrough
	case 4009:
		return true
	}
	return false
}
