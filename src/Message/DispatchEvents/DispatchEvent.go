package DispatchEvents

const (
	INTERACTION_CREATE string = "INTERACTION_CREATE"
)

type DispatchEvent struct {
}

func GetDispatchData(t string, data interface{}) interface{} {
	switch t {
	case INTERACTION_CREATE:
		var result = new(InteractionCreateEvent)
		result.Data = new(InteractionDataApplicationCommand)
		return result
	}
	return nil
}
