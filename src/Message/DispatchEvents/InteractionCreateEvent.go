package DispatchEvents

const (
	INTERACTION_TYPE_PING                             int = 1
	INTERACTION_TYPE_APPLICATION_COMMAND                  = 2
	INTERACTION_TYPE_MESSAGE_COMPONENT                    = 3
	INTERACTION_TYPE_APPLICATION_COMMAND_AUTOCOMPLETE     = 4
	INTERACTION_TYPE_MODAL_SUBMIT                         = 5
)

type InteractionDataOption struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}
type InteractionDataApplicationCommand struct {
	Name    string                  `json:"name"`
	Options []InteractionDataOption `json:"options"`
}
type InteractionCreateEvent struct {
	Type  int         `json:"type"`
	Token string      `json:"token"`
	Id    string      `json:"id"`
	Data  interface{} `json:"data"`
}
