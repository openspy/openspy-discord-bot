package Manager

import (
	"log"

	"github.com/gorilla/websocket"
	"openspy.net/discord-gateway-bot/IManager"
	"openspy.net/discord-gateway-bot/Manager/Handler"
	"openspy.net/discord-gateway-bot/Message"
)

type StateManager struct {
	dialer       websocket.Dialer
	address      string
	stateHandler Handler.StateHandler
	connection   *websocket.Conn

	tickerInterval int

	dispatchChannel chan Message.Event
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
func (s *StateManager) SetTickerInterval(interval int) {
	s.tickerInterval = interval
}
func (s *StateManager) GetTickerInterval() int {
	return s.tickerInterval
}
func (s *StateManager) SetState(state int) {
	switch state {
	case IManager.STATE_ESTABLISH:
		s.stateHandler = &Handler.HelloStateHandler{}
	case IManager.STATE_HEARTBEAT:
		s.stateHandler = &Handler.HeartbeatStateHandler{}
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
			log.Println("read:", err)
			return
		}

		var event Message.Event
		event.UnmarshalJSON(message)
		s.stateHandler.HandleMessage(event)
	}
}

func (s *StateManager) WriteMessage(event Message.Event) {
	s.connection.WriteJSON(event)
}
func (s *StateManager) GetDispatchChannel() chan Message.Event {
	return s.dispatchChannel
}
