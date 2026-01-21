package Message

import (
	"encoding/json"
	"log"

	"openspy.net/discord-gateway-bot/Message/DispatchEvents"
)

const (
	DISPATCH_EVENT              int = 0
	HEARTBEAT_EVENT                 = 1
	IDENTIFY_EVENT                  = 2
	PRESENCE_UPDATE_EVENT           = 3
	VOICE_STATUS_EVENT              = 4
	RESUME_EVENT                    = 6
	RECONNECT_EVENT                 = 7
	REQUEST_GUILD_MEMBERS_EVENT     = 8
	INVALID_SESSION_EVENT           = 9
	HELLO_EVENT                     = 10
	HEARTBEAT_ACK                   = 11
	REQUEST_SOUNDBOARD_SOUNDS       = 31
)

const (
	GATEWAY_UNKNOWN_ERROR             int = 4000
	GATEWAY_UNKNOWN_OPCODE_ERROR          = 4001
	GATEWAY_DECODE_ERROR                  = 4002
	GATEWAY_NOT_AUTHENTICATED_ERROR       = 4003
	GATEWAY_AUTHENTICATION_FAILED         = 4004
	GATEWAY_ALREADY_AUTHENTICATED         = 4005
	GATEWAY_INVALID_SEQ_ERROR             = 4007
	GATEWAY_RATE_LIMITED_ERROR            = 4008
	GATEWAY_SESSION_TIMEOUT_ERROR         = 4009
	GATEWAY_INVALID_SHARD_ERROR           = 4010
	GATEWAY_SHARDING_REQUIRED_ERROR       = 4011
	GATEWAY_INVALID_API_VERSION_ERROR     = 4012
	GATEWAY_INVALID_INTENTS_ERROR         = 4013
	GATEWAY_DISALLOWED_INTENTS_ERROR      = 4014
)

type Event struct {
	Op   int         `json:"op"`
	Type string      `json:"t"`
	Data interface{} `json:"d"`
}

func (b *Event) UnmarshalJSON(data []byte) (Event, error) {
	var typ struct {
		Op   int         `json:"op"`
		Type string      `json:"t"`
		Data interface{} `json:"d"`
	}
	if err := json.Unmarshal(data, &typ); err != nil {
		log.Printf("got err: %s\n", err)
		return Event{}, err
	}
	b.Op = typ.Op
	switch typ.Op {
	case HELLO_EVENT:
		b.Data = new(HelloEvent)
	case DISPATCH_EVENT:
		b.Data = DispatchEvents.GetDispatchData(typ.Type, typ.Data)
	}
	type tmp Event // avoids infinite recursion
	unmarshalErr := json.Unmarshal(data, (*tmp)(b))

	return *b, unmarshalErr

}
